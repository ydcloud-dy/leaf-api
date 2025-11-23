<template>
  <div class="articles">
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">文章列表</h1>
        <p class="page-subtitle">探索知识的海洋</p>
      </div>

      <!-- 筛选和排序 -->
      <div class="filters">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索文章..."
          clearable
          style="width: 300px"
          @input="handleSearchInput"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>

        <div class="filter-actions">
          <el-select
            v-model="selectedTag"
            placeholder="选择标签"
            clearable
            style="width: 150px"
            @change="handleTagChange"
          >
            <el-option
              v-for="tag in tags"
              :key="tag.id"
              :label="tag.name"
              :value="tag.name"
            />
          </el-select>

          <el-select
            v-model="selectedCategory"
            placeholder="选择分类"
            clearable
            style="width: 150px"
            @change="handleFilter"
          >
            <el-option
              v-for="category in categories"
              :key="category"
              :label="category"
              :value="category"
            />
          </el-select>

          <el-select
            v-model="sortBy"
            placeholder="排序方式"
            style="width: 150px"
            @change="handleFilter"
          >
            <el-option label="最新发布" value="created_at" />
            <el-option label="最多浏览" value="views" />
            <el-option label="最多点赞" value="likes" />
          </el-select>
        </div>
      </div>

      <!-- 文章列表 -->
      <div v-if="!showChapterView" v-loading="loading" class="articles-list">
        <ArticleCard
          v-for="article in articles"
          :key="article.id"
          :article="article"
        />

        <el-empty v-if="!loading && !articles.length" description="暂无文章" />
      </div>

      <!-- 章节目录视图 -->
      <div v-else v-loading="loading" class="chapter-view">
        <el-card v-for="article in articlesWithChapters" :key="article.id" class="article-chapters">
          <template #header>
            <div class="chapter-header" @click="$router.push(`/articles/${article.id}`)">
              <h3>{{ article.title }}</h3>
              <span class="article-meta">
                <el-icon><View /></el-icon>
                {{ article.view_count || 0 }}
              </span>
            </div>
          </template>
          <div class="chapters-list">
            <div
              v-for="(chapter, index) in article.chapters"
              :key="index"
              :class="['chapter-item', `chapter-level-${chapter.level}`]"
              @click="scrollToChapter(article.id, chapter.id)"
            >
              {{ chapter.text }}
            </div>
            <el-empty v-if="!article.chapters || article.chapters.length === 0" description="暂无章节" />
          </div>
        </el-card>

        <el-empty v-if="!loading && !articlesWithChapters.length" description="该标签下暂无文章" />
      </div>

      <!-- 分页 -->
      <div v-if="total > pageSize" class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next, jumper"
          @current-change="handlePageChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getArticles, searchArticles } from '@/api/article'
import { getTags } from '@/api/tag'
import ArticleCard from '@/components/ArticleCard.vue'
import { View } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()

const articles = ref([])
const articlesWithChapters = ref([])
const tags = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const searchKeyword = ref('')
const selectedTag = ref('')  // 改为存储标签名称而不是ID
const selectedCategory = ref('')
const sortBy = ref('created_at')
const showChapterView = ref(false)
const categories = ref(['前端开发', '后端开发', '数据库', '运维部署', '算法', '其他'])
let searchTimer = null

onMounted(() => {
  // 从URL获取查询参数
  searchKeyword.value = route.query.keyword || ''
  selectedCategory.value = route.query.category || ''

  fetchTags()
  fetchArticles()
})

watch(() => route.query, () => {
  searchKeyword.value = route.query.keyword || ''
  selectedCategory.value = route.query.category || ''
  currentPage.value = 1
  fetchArticles()
})

const fetchTags = async () => {
  try {
    const { data } = await getTags()
    const tagList = Array.isArray(data) ? data : (data.list || [])
    tags.value = tagList
  } catch (error) {
    console.error('Failed to fetch tags:', error)
  }
}

const fetchArticles = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      sort: sortBy.value,
      status: 1 // 只获取已发布的文章
    }

    if (searchKeyword.value && searchKeyword.value.trim()) {
      params.keyword = searchKeyword.value.trim()
    }

    if (selectedCategory.value) {
      params.category = selectedCategory.value
    }

    if (selectedTag.value) {
      params.tag = selectedTag.value
    }

    // 如果有搜索关键词才使用searchArticles,否则使用getArticles
    const apiCall = (searchKeyword.value && searchKeyword.value.trim()) ? searchArticles : getArticles
    const { data } = await apiCall(params)

    articles.value = data.list || []
    total.value = data.total || 0

    // 如果选择了标签,提取章节信息
    if (selectedTag.value) {
      extractChapters(articles.value)
    }
  } catch (error) {
    console.error('Failed to fetch articles:', error)
  } finally {
    loading.value = false
  }
}

// 从文章内容中提取章节标题
const extractChapters = (articlesList) => {
  articlesWithChapters.value = articlesList.map(article => {
    const chapters = []

    // 使用正则表达式从markdown内容中提取标题
    const content = article.content_markdown || article.content_html || ''

    // 匹配 markdown 标题 (# 开头)
    const mdHeadingRegex = /^(#{1,6})\s+(.+)$/gm
    let match
    let index = 0

    while ((match = mdHeadingRegex.exec(content)) !== null) {
      const level = match[1].length // # 的数量就是层级
      const text = match[2].trim()

      chapters.push({
        id: `heading-${index}`,
        text: text,
        level: level
      })
      index++
    }

    return {
      ...article,
      chapters: chapters
    }
  })
}

const handleTagChange = () => {
  if (selectedTag.value) {
    showChapterView.value = true
  } else {
    showChapterView.value = false
  }
  currentPage.value = 1
  fetchArticles()
}

const handleSearchInput = () => {
  // 清除之前的定时器
  if (searchTimer) {
    clearTimeout(searchTimer)
  }

  // 设置新的定时器,300ms后执行搜索
  searchTimer = setTimeout(() => {
    currentPage.value = 1
    fetchArticles()
  }, 300)
}

const handleSearch = () => {
  currentPage.value = 1
  fetchArticles()
}

const handleFilter = () => {
  currentPage.value = 1
  fetchArticles()
}

const handlePageChange = () => {
  fetchArticles()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const scrollToChapter = (articleId, chapterId) => {
  // 跳转到文章详情页并滚动到对应章节
  router.push({
    path: `/articles/${articleId}`,
    hash: `#${chapterId}`
  })
}
</script>

<style scoped>
.articles {
  padding: 20px 0;
}

.filters {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding: 20px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.filter-actions {
  display: flex;
  gap: 12px;
}

.articles-list {
  min-height: 500px;
}

/* 章节视图样式 */
.chapter-view {
  min-height: 500px;
}

.article-chapters {
  margin-bottom: 20px;
}

.chapter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
  transition: color 0.3s;
}

.chapter-header:hover {
  color: #409eff;
}

.chapter-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.article-meta {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #909399;
  font-size: 14px;
}

.chapters-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.chapter-item {
  padding: 10px 12px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
  color: #606266;
  font-size: 14px;
  border-left: 3px solid transparent;
}

.chapter-item:hover {
  background-color: #f5f7fa;
  color: #409eff;
  border-left-color: #409eff;
}

/* 章节层级样式 */
.chapter-item.chapter-level-1 {
  padding-left: 12px;
  font-weight: 600;
  font-size: 16px;
}

.chapter-item.chapter-level-2 {
  padding-left: 24px;
  font-size: 15px;
}

.chapter-item.chapter-level-3 {
  padding-left: 36px;
  font-size: 14px;
}

.chapter-item.chapter-level-4 {
  padding-left: 48px;
  font-size: 13px;
  color: #909399;
}

.chapter-item.chapter-level-5 {
  padding-left: 60px;
  font-size: 12px;
  color: #909399;
}

.chapter-item.chapter-level-6 {
  padding-left: 72px;
  font-size: 12px;
  color: #909399;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 40px;
}

@media (max-width: 768px) {
  .filters {
    flex-direction: column;
    gap: 16px;
  }

  .filter-actions {
    width: 100%;
    flex-direction: column;
  }

  .filter-actions .el-select {
    width: 100% !important;
  }

  .chapter-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}
</style>
