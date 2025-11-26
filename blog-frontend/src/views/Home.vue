<template>
  <div class="home">
    <div class="banner">
      <div class="container">
        <h1 class="banner-title">欢迎来到我的博客</h1>
        <p class="banner-subtitle">分享技术，记录生活</p>
      </div>
    </div>

    <div class="container">
      <div class="content-wrapper">
        <!-- 文章列表 -->
        <div class="articles-section">
          <div v-loading="loading" class="articles-list">
            <ArticleCard
              v-for="article in articles"
              :key="article.id"
              :article="article"
            />

            <el-empty v-if="!loading && !articles.length" description="暂无文章" />
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

        <!-- 侧边栏 -->
        <aside class="sidebar">
          <!-- 热门文章 -->
          <el-card class="sidebar-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <el-icon><TrendCharts /></el-icon>
                <span>热门文章</span>
              </div>
            </template>
            <div class="hot-articles">
              <div
                v-for="(article, index) in hotArticles"
                :key="article.id"
                class="hot-article-item"
                @click="router.push(`/articles/${article.id}`)"
              >
                <span class="rank" :class="{ top: index < 3 }">{{ index + 1 }}</span>
                <span class="title">{{ article.title }}</span>
              </div>
            </div>
          </el-card>

          <!-- 标签云 -->
          <el-card class="sidebar-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <el-icon><CollectionTag /></el-icon>
                <span>标签云</span>
              </div>
            </template>
            <div class="tags-cloud">
              <el-tag
                v-for="tag in tags"
                :key="tag"
                class="tag-item"
                @click="handleTagClick(tag)"
              >
                {{ tag }}
              </el-tag>
            </div>
          </el-card>

          <!-- 站点统计 -->
          <el-card class="sidebar-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <el-icon><DataAnalysis /></el-icon>
                <span>站点统计</span>
              </div>
            </template>
            <div class="stats">
              <div class="stat-item">
                <el-icon><Document /></el-icon>
                <span>文章数：{{ stats.articles || 0 }}</span>
              </div>
              <div class="stat-item">
                <el-icon><View /></el-icon>
                <span>访问量：{{ stats.views || 0 }}</span>
              </div>
              <div class="stat-item">
                <el-icon><ChatDotRound /></el-icon>
                <span>评论数：{{ stats.comments || 0 }}</span>
              </div>
            </div>
          </el-card>
        </aside>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getArticles } from '@/api/article'
import { getStats, getHotArticles } from '@/api/stats'
import { getTags } from '@/api/tag'
import ArticleCard from '@/components/ArticleCard.vue'

const router = useRouter()

const articles = ref([])
const hotArticles = ref([])
const tags = ref([])
const stats = ref({})
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const loading = ref(false)

onMounted(() => {
  fetchArticles()
  fetchHotArticles()
  fetchTags()
  fetchStats()
})

const fetchArticles = async () => {
  loading.value = true
  try {
    const { data } = await getArticles({
      page: currentPage.value,
      page_size: pageSize.value,
      status: 1 // 只获取已发布的文章
    })
    articles.value = data.list || []
    total.value = data.total || 0
  } catch (error) {
    console.error('Failed to fetch articles:', error)
  } finally {
    loading.value = false
  }
}

const fetchHotArticles = async () => {
  try {
    const { data } = await getHotArticles()
    hotArticles.value = data || []
  } catch (error) {
    console.error('Failed to fetch hot articles:', error)
  }
}

const fetchTags = async () => {
  try {
    const { data } = await getTags()
    // 提取标签名称 - TagService 返回数组而非 {list: [...]}
    const tagList = Array.isArray(data) ? data : (data.list || [])
    tags.value = tagList.map(tag => tag.name)
  } catch (error) {
    console.error('Failed to fetch tags:', error)
    tags.value = []
  }
}

const fetchStats = async () => {
  try {
    const { data } = await getStats()
    stats.value = {
      articles: data.article_count || 0,
      views: data.total_views || 0,
      comments: data.comment_count || 0
    }
  } catch (error) {
    console.error('Failed to fetch stats:', error)
    // 如果API失败，使用默认值
    stats.value = {
      articles: 0,
      views: 0,
      comments: 0
    }
  }
}

const handlePageChange = () => {
  fetchArticles()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const handleTagClick = (tag) => {
  router.push({
    name: 'Articles',
    query: { tag }
  })
}
</script>

<style scoped>
.banner {
  background:
    linear-gradient(rgba(0, 0, 0, 0.4), rgba(0, 0, 0, 0.4)),
    url('../../img/wukong.png');
  background-size: cover;
  height: 500px;
  background-position: center;
  background-attachment: fixed;
  color: #fff;
  padding: 80px 0;
  margin-bottom: 40px;
}

.banner-title {
  font-size: 48px;
  font-weight: 700;
  margin-top: 100px;
  margin-bottom: 16px;
  text-align: center;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
}

.banner-subtitle {
  font-size: 20px;
  text-align: center;
  opacity: 0.95;
  text-shadow: 1px 1px 3px rgba(0, 0, 0, 0.5);
}

.content-wrapper {
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: 30px;
}

.articles-section {
  min-height: 500px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 40px;
}

.sidebar {
  position: sticky;
  top: 80px;
  height: fit-content;
}

.sidebar-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #303133;
}

.hot-articles {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.hot-article-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.hot-article-item:hover {
  background-color: #f5f7fa;
}

.rank {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #e4e7ed;
  color: #606266;
  border-radius: 50%;
  font-size: 12px;
  font-weight: 600;
  flex-shrink: 0;
}

.rank.top {
  background-color: #f56c6c;
  color: #fff;
}

.hot-article-item .title {
  flex: 1;
  font-size: 14px;
  color: #606266;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tags-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag-item {
  cursor: pointer;
  transition: transform 0.2s;
}

.tag-item:hover {
  transform: scale(1.05);
}

.stats {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #606266;
  font-size: 14px;
}

@media (max-width: 968px) {
  .content-wrapper {
    grid-template-columns: 1fr;
  }

  .sidebar {
    position: static;
  }

  .banner-title {
    font-size: 32px;
  }

  .banner-subtitle {
    font-size: 16px;
  }
}
</style>
