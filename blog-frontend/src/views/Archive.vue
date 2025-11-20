<template>
  <div class="archive">
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">文章归档</h1>
        <p class="page-subtitle">时光荏苒，记录点滴</p>
      </div>

      <div v-loading="loading" class="archive-content">
        <div v-for="year in archives" :key="year.year" class="year-group">
          <div class="year-header">
            <h2 class="year-title">{{ year.year }}</h2>
            <span class="year-count">共 {{ year.count }} 篇</span>
          </div>

          <div v-for="month in year.months" :key="month.month" class="month-group">
            <div class="month-header">
              <h3 class="month-title">{{ month.month }}月</h3>
            </div>

            <div class="articles-list">
              <div
                v-for="article in month.articles"
                :key="article.id"
                class="article-item"
                @click="router.push(`/articles/${article.id}`)"
              >
                <div class="article-date">
                  {{ formatDay(article.created_at) }}
                </div>
                <div class="article-info">
                  <h4 class="article-title">{{ article.title }}</h4>
                  <div class="article-meta">
                    <el-tag
                      v-if="article.category"
                      size="small"
                      type="primary"
                      effect="plain"
                    >
                      {{ article.category }}
                    </el-tag>
                    <span class="meta-item">
                      <el-icon><View /></el-icon>
                      {{ article.views || 0 }}
                    </span>
                    <span class="meta-item">
                      <el-icon><ChatDotRound /></el-icon>
                      {{ article.comments || 0 }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <el-empty v-if="!loading && !archives.length" description="暂无文章" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getArchiveArticles } from '@/api/article'
import dayjs from 'dayjs'

const router = useRouter()

const archives = ref([])
const loading = ref(false)

onMounted(() => {
  fetchArchives()
})

const fetchArchives = async () => {
  loading.value = true
  try {
    const { data } = await getArchiveArticles()
    archives.value = processArchiveData(data.list || [])
  } catch (error) {
    console.error('Failed to fetch archives:', error)
    // 如果接口不存在，使用备用方案
    archives.value = []
  } finally {
    loading.value = false
  }
}

// 处理归档数据
const processArchiveData = (articles) => {
  const yearMap = new Map()

  articles.forEach(article => {
    const date = dayjs(article.created_at)
    const year = date.year()
    const month = date.month() + 1

    if (!yearMap.has(year)) {
      yearMap.set(year, new Map())
    }

    const monthMap = yearMap.get(year)
    if (!monthMap.has(month)) {
      monthMap.set(month, [])
    }

    monthMap.get(month).push(article)
  })

  // 转换为数组格式并排序
  const result = []
  const sortedYears = Array.from(yearMap.keys()).sort((a, b) => b - a)

  sortedYears.forEach(year => {
    const monthMap = yearMap.get(year)
    const sortedMonths = Array.from(monthMap.keys()).sort((a, b) => b - a)

    const months = sortedMonths.map(month => ({
      month,
      articles: monthMap.get(month).sort((a, b) =>
        dayjs(b.created_at).unix() - dayjs(a.created_at).unix()
      )
    }))

    const totalCount = months.reduce((sum, month) => sum + month.articles.length, 0)

    result.push({
      year,
      count: totalCount,
      months
    })
  })

  return result
}

const formatDay = (date) => {
  return dayjs(date).format('MM-DD')
}
</script>

<style scoped>
.archive {
  padding: 20px 0;
}

.archive-content {
  min-height: 500px;
}

.year-group {
  margin-bottom: 60px;
}

.year-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 30px;
  padding-bottom: 16px;
  border-bottom: 2px solid #409eff;
}

.year-title {
  font-size: 32px;
  font-weight: 700;
  color: #409eff;
  margin: 0;
}

.year-count {
  font-size: 14px;
  color: #909399;
  background-color: #f5f7fa;
  padding: 4px 12px;
  border-radius: 12px;
}

.month-group {
  margin-bottom: 40px;
}

.month-header {
  margin-bottom: 20px;
}

.month-title {
  font-size: 20px;
  font-weight: 600;
  color: #606266;
  margin: 0;
  padding-left: 12px;
  border-left: 3px solid #67c23a;
}

.articles-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.article-item {
  display: flex;
  gap: 20px;
  padding: 16px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.3s;
}

.article-item:hover {
  transform: translateX(8px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.article-date {
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  font-size: 18px;
  font-weight: 600;
  border-radius: 8px;
  flex-shrink: 0;
}

.article-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.article-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.article-meta {
  display: flex;
  align-items: center;
  gap: 16px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
  color: #909399;
}

@media (max-width: 768px) {
  .year-title {
    font-size: 24px;
  }

  .article-item {
    flex-direction: column;
  }

  .article-date {
    width: 100%;
    height: 40px;
  }
}
</style>
