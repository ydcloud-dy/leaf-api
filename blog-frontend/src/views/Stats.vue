<template>
  <div class="stats">
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">📊 网站统计</h1>
        <p class="page-subtitle">实时统计数据，每30秒自动刷新</p>
      </div>

      <div class="stats-grid" v-loading="loading">
        <!-- 网站运行时长 -->
        <div class="stat-card">
          <div class="stat-icon">🔄</div>
          <div class="stat-label">网站运行时长</div>
          <div class="stat-value">{{ stats.site_runtime }}天</div>
        </div>

        <!-- 24小时访问量 -->
        <div class="stat-card">
          <div class="stat-icon">📈</div>
          <div class="stat-label">24小时访问量（PV）</div>
          <div class="stat-value">{{ stats.today_views }}次</div>
        </div>

        <!-- 当前在线人数 -->
        <div class="stat-card highlight">
          <div class="stat-icon">🌐</div>
          <div class="stat-label">当前在线人数</div>
          <div class="stat-value">{{ stats.online_count }}人</div>
        </div>

        <!-- 平均访问时长 -->
        <div class="stat-card">
          <div class="stat-icon">⏱️</div>
          <div class="stat-label">平均访问时长</div>
          <div class="stat-value">{{ formatDuration(stats.avg_visit_duration) }}</div>
        </div>

        <!-- 文章篇数 -->
        <div class="stat-card">
          <div class="stat-icon">📝</div>
          <div class="stat-label">文章篇数</div>
          <div class="stat-value">{{ stats.article_count }}篇</div>
        </div>

        <!-- 笔记篇数 -->
        <div class="stat-card">
          <div class="stat-icon">📔</div>
          <div class="stat-label">笔记篇数</div>
          <div class="stat-value">{{ stats.chapter_count }}篇</div>
        </div>

        <!-- 文章分类数 -->
        <div class="stat-card">
          <div class="stat-icon">📚</div>
          <div class="stat-label">文章分类数</div>
          <div class="stat-value">{{ stats.category_count }}个</div>
        </div>

        <!-- 文章标签数 -->
        <div class="stat-card">
          <div class="stat-icon">🏷️</div>
          <div class="stat-label">文章标签数</div>
          <div class="stat-value">{{ stats.tag_count }}个</div>
        </div>

        <!-- 总浏览量 -->
        <div class="stat-card">
          <div class="stat-icon">👀</div>
          <div class="stat-label">总浏览量</div>
          <div class="stat-value">{{ formatNumber(stats.total_views) }}</div>
        </div>

        <!-- 评论总数 -->
        <div class="stat-card">
          <div class="stat-icon">💬</div>
          <div class="stat-label">评论总数</div>
          <div class="stat-value">{{ stats.comment_count }}条</div>
        </div>

        <!-- 用户总数 -->
        <div class="stat-card">
          <div class="stat-icon">👥</div>
          <div class="stat-label">注册用户数</div>
          <div class="stat-value">{{ stats.user_count }}人</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { getStats } from '@/api/stats'

const loading = ref(false)
const stats = ref({
  article_count: 0,
  chapter_count: 0,
  category_count: 0,
  tag_count: 0,
  user_count: 0,
  comment_count: 0,
  total_views: 0,
  today_views: 0,
  online_count: 0,
  avg_visit_duration: 0,
  site_runtime: 0
})

let refreshTimer = null

// 格式化时长（秒转为分钟或秒）
const formatDuration = (seconds) => {
  if (seconds === 0) return '0秒/页'
  if (seconds < 60) {
    return `${Math.round(seconds)}秒/页`
  }
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = Math.round(seconds % 60)
  if (remainingSeconds === 0) {
    return `${minutes}分钟/页`
  }
  return `${minutes}分${remainingSeconds}秒/页`
}

// 格式化数字（添加千分位）
const formatNumber = (num) => {
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

// 加载统计数据
const loadStats = async () => {
  try {
    loading.value = true
    const res = await getStats()
    if (res.data) {
      stats.value = res.data
    }
  } catch (error) {
    console.error('Failed to load stats:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadStats()

  // 每 30 秒自动刷新统计数据
  refreshTimer = setInterval(() => {
    loadStats()
  }, 30000)
})

onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
})
</script>

<style scoped>
.stats {
  min-height: 100vh;
  padding: 20px 0 60px;
  background:
    url('https://images.unsplash.com/photo-1557683316-973673baf926?w=1920&q=80');
  background-size: cover;
  background-position: center;
  background-attachment: fixed;
}

.page-header {
  padding: 40px 0 30px;
  text-align: center;
}

.page-title {
  font-size: 36px;
  font-weight: 700;
  color: white;
  margin-bottom: 10px;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
}

.page-subtitle {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.95);
  text-shadow: 1px 1px 3px rgba(0, 0, 0, 0.5);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 24px;
  margin-top: 40px;
}

.stat-card {
  background: white;
  border-radius: 16px;
  padding: 32px 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  text-align: center;
  transition: all 0.3s ease;
  cursor: default;
}

.stat-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.stat-card.highlight {
  background: #409eff;
  color: white;
  position: relative;
  overflow: hidden;
}

.stat-card.highlight::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.15) 0%, rgba(255, 255, 255, 0) 100%);
  pointer-events: none;
}

.stat-icon {
  font-size: 48px;
  margin-bottom: 16px;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

.stat-label {
  font-size: 15px;
  color: #666;
  margin-bottom: 12px;
  font-weight: 500;
}

.stat-card.highlight .stat-label {
  color: rgba(255, 255, 255, 0.95);
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: #333;
  line-height: 1;
}

.stat-card.highlight .stat-value {
  color: white;
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    gap: 16px;
  }

  .stat-card {
    padding: 24px 16px;
  }

  .stat-icon {
    font-size: 36px;
    margin-bottom: 12px;
  }

  .stat-label {
    font-size: 13px;
  }

  .stat-value {
    font-size: 24px;
  }

  .page-title {
    font-size: 28px;
  }
}
</style>
