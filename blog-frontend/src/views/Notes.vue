<template>
  <div class="notes-page">
    <div class="notes-container">
      <!-- 左侧边栏 - 分类列表 -->
      <aside class="notes-sidebar">
        <div class="sidebar-header">
          <el-icon><Notebook /></el-icon>
          <span>笔记分类</span>
        </div>
        <div class="category-list">
          <div
            v-for="tag in tags"
            :key="tag.id"
            :class="['category-item', { active: currentTag === tag.name }]"
            @click="selectTag(tag.name)"
          >
            {{ tag.name }}
          </div>
        </div>
      </aside>

      <!-- 右侧主内容区 - 章节目录 -->
      <main class="notes-main" v-loading="loading">
        <div v-if="currentTag" class="chapters-content">
          <h2 class="content-title">{{ currentTag }}</h2>

          <!-- 章节列表 -->
          <div v-for="chapter in chapters" :key="chapter.id" class="chapter-section">
            <!-- 章节标题 -->
            <div
              class="chapter-header"
              :class="{ expanded: expandedChapters.includes(chapter.id) }"
              @click="toggleChapter(chapter.id)"
            >
              <el-icon class="expand-icon">
                <ArrowRight v-if="!expandedChapters.includes(chapter.id)" />
                <ArrowDown v-else />
              </el-icon>
              <h3 class="chapter-name">{{ chapter.name }}</h3>
              <span class="article-count">({{ chapter.articles?.length || 0 }})</span>
            </div>

            <!-- 章节下的文章列表 -->
            <div v-show="expandedChapters.includes(chapter.id)" class="articles-list">
              <div v-for="article in chapter.articles" :key="article.id" class="article-item" @click="goToArticle(article.id)">
                <!-- 文章标题行 -->
                <div class="article-title-row">
                  <el-icon class="article-icon"><Document /></el-icon>
                  <span class="article-title">
                    {{ article.title }}
                  </span>
                </div>

                <div class="article-meta">
                  <span><el-icon><View /></el-icon> {{ article.view_count || 0 }}</span>
                  <span><el-icon><Calendar /></el-icon> {{ formatDate(article.created_at) }}</span>
                </div>
              </div>

              <el-empty v-if="!chapter.articles || chapter.articles.length === 0" description="该章节下暂无文章" :image-size="60" />
            </div>
          </div>

          <el-empty v-if="!loading && chapters.length === 0" description="该分类下暂无章节" />
        </div>

        <div v-else class="empty-state">
          <el-empty description="请从左侧选择一个分类" />
        </div>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getChaptersByTag } from '@/api/chapter'
import { getTags } from '@/api/tag'
import { Notebook, ArrowRight, ArrowDown, Document, View, Calendar } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()

const tags = ref([])
const chapters = ref([])
const loading = ref(false)
const currentTag = ref('')
const expandedChapters = ref([])

onMounted(() => {
  fetchTags()
  // 从路由参数获取标签
  if (route.params.tag) {
    currentTag.value = route.params.tag
    fetchChapters()
  }
})

watch(() => route.params.tag, (newTag) => {
  if (newTag) {
    currentTag.value = newTag
    fetchChapters()
  }
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

const fetchChapters = async () => {
  if (!currentTag.value) return

  loading.value = true
  try {
    const { data } = await getChaptersByTag(currentTag.value)
    chapters.value = data || []
    // 默认展开第一个章节
    if (chapters.value.length > 0) {
      expandedChapters.value = [chapters.value[0].id]
    }
  } catch (error) {
    console.error('Failed to fetch chapters:', error)
  } finally {
    loading.value = false
  }
}

const selectTag = (tagName) => {
  currentTag.value = tagName
  router.push({ name: 'Notes', params: { tag: tagName } })
  fetchChapters()
}

const toggleChapter = (chapterId) => {
  const index = expandedChapters.value.indexOf(chapterId)
  if (index > -1) {
    expandedChapters.value.splice(index, 1)
  } else {
    expandedChapters.value.push(chapterId)
  }
}

const goToArticle = (articleId) => {
  router.push(`/articles/${articleId}`)
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}
</script>

<style scoped>
.notes-page {
  min-height: calc(100vh - 60px);
  background-color: #f5f7fa;
}

.notes-container {
  display: flex;
  max-width: 1400px;
  margin: 0 auto;
  min-height: calc(100vh - 60px);
}

/* 左侧边栏 */
.notes-sidebar {
  width: 260px;
  background-color: #fff;
  border-right: 1px solid #e4e7ed;
  flex-shrink: 0;
}

.sidebar-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 20px;
  font-size: 18px;
  font-weight: 600;
  color: #409eff;
  border-bottom: 1px solid #e4e7ed;
}

.category-list {
  padding: 10px 0;
}

.category-item {
  padding: 12px 20px;
  color: #606266;
  cursor: pointer;
  transition: all 0.3s;
  border-left: 3px solid transparent;
}

.category-item:hover {
  background-color: #f5f7fa;
  color: #409eff;
}

.category-item.active {
  background-color: #ecf5ff;
  color: #409eff;
  border-left-color: #409eff;
}

/* 右侧主内容 */
.notes-main {
  flex: 1;
  padding: 30px;
  overflow-y: auto;
  background-color: #f5f7fa;
}

.content-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 30px;
  padding-bottom: 15px;
  border-bottom: 2px solid #409eff;
}

.chapter-section {
  margin-bottom: 15px;
}

.chapter-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 15px 20px;
  background-color: #409eff;
  color: #fff;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.3);
}

.chapter-header:hover {
  background-color: #66b1ff;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.4);
}

.chapter-header.expanded {
  border-radius: 8px 8px 0 0;
}

.expand-icon {
  flex-shrink: 0;
  font-size: 16px;
  color: #fff;
  transition: transform 0.3s;
}

.chapter-name {
  flex: 1;
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #fff;
}

.article-count {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
  font-weight: normal;
}

.articles-list {
  background-color: #fff;
  border: 1px solid #e4e7ed;
  border-top: none;
  border-radius: 0 0 8px 8px;
  padding: 15px;
}

.article-item {
  padding: 15px;
  margin-bottom: 10px;
  border-radius: 6px;
  background-color: #fafafa;
  border-left: 3px solid #409eff;
  transition: all 0.3s;
  cursor: pointer;
}

.article-item:hover {
  background-color: #f0f7ff;
  box-shadow: 0 2px 6px rgba(64, 158, 255, 0.15);
  transform: translateX(4px);
}

.article-item:last-child {
  margin-bottom: 0;
}

.article-title-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.article-icon {
  color: #409eff;
  font-size: 18px;
  flex-shrink: 0;
}

.article-title {
  flex: 1;
  font-size: 16px;
  font-weight: 500;
  color: #303133;
  transition: color 0.3s;
}

.article-item:hover .article-title {
  color: #409eff;
}

.article-meta {
  display: flex;
  gap: 15px;
  padding-left: 28px;
  font-size: 12px;
  color: #909399;
}

.article-meta span {
  display: flex;
  align-items: center;
  gap: 4px;
}

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 400px;
}

/* 响应式 */
@media (max-width: 768px) {
  .notes-container {
    flex-direction: column;
  }

  .notes-sidebar {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid #e4e7ed;
  }

  .category-list {
    display: flex;
    flex-wrap: wrap;
    padding: 10px;
    gap: 10px;
  }

  .category-item {
    padding: 8px 16px;
    border-radius: 20px;
    border-left: none;
    background-color: #f5f7fa;
  }

  .category-item.active {
    background-color: #409eff;
    color: #fff;
  }
}
</style>
