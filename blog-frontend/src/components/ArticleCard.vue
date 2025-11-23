<template>
  <el-card class="article-card" shadow="hover" @click="handleClick">
    <div class="card-wrapper">
      <div v-if="article.cover" class="cover">
        <img :src="article.cover" :alt="article.title" />
      </div>

      <div class="content">
        <h3 class="title">{{ article.title }}</h3>

        <p class="summary">{{ getSummary() }}</p>

        <div class="meta">
          <div class="tags">
            <el-tag
              v-if="article.category"
              type="primary"
              size="small"
              effect="plain"
            >
              {{ typeof article.category === 'object' ? article.category.name : article.category }}
            </el-tag>
            <el-tag
              v-for="tag in article.tags?.slice(0, 3)"
              :key="tag.id || tag"
              size="small"
              effect="plain"
            >
              {{ typeof tag === 'object' ? tag.name : tag }}
            </el-tag>
          </div>

          <div class="stats">
            <span class="stat-item">
              <el-icon><View /></el-icon>
              {{ article.view_count || 0 }}
            </span>
            <span class="stat-item">
              <el-icon><ChatDotRound /></el-icon>
              {{ article.comment_count || 0 }}
            </span>
            <span class="stat-item">
              <el-icon><Star /></el-icon>
              {{ article.like_count || 0 }}
            </span>
          </div>
        </div>

        <div class="footer">
          <div class="author">
            <el-avatar :size="24" :src="article.author?.avatar">
              {{ article.author?.username?.charAt(0).toUpperCase() }}
            </el-avatar>
            <span>{{ article.author?.username || '匿名' }}</span>
          </div>
          <div class="date">
            {{ formatDate(article.created_at) }}
          </div>
        </div>
      </div>
    </div>
  </el-card>
</template>

<script setup>
import { useRouter } from 'vue-router'
import dayjs from 'dayjs'

const props = defineProps({
  article: {
    type: Object,
    required: true
  }
})

const router = useRouter()

const handleClick = () => {
  router.push(`/articles/${props.article.id}`)
}

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD')
}

const getSummary = () => {
  // 优先使用 summary 字段
  if (props.article.summary && props.article.summary.trim()) {
    return props.article.summary
  }

  // 然后尝试从 content_markdown 或 content 中提取
  const content = props.article.content_markdown || props.article.content || ''

  if (content.trim()) {
    // 移除 markdown 标记和多余空白
    const plainText = content
      .replace(/[#*`>\[\]()]/g, '') // 移除常见 markdown 符号
      .replace(/\s+/g, ' ')          // 将多个空白字符替换为单个空格
      .trim()

    if (plainText.length > 200) {
      return plainText.substring(0, 200) + '...'
    }
    return plainText
  }

  return '暂无简介'
}
</script>

<style scoped>
.article-card {
  cursor: pointer;
  transition: transform 0.3s;
  margin-bottom: 20px;
}

.article-card:hover {
  transform: translateY(-4px);
}

.card-wrapper {
  display: flex;
  gap: 20px;
}

.cover {
  flex-shrink: 0;
  width: 280px;
  height: 180px;
  overflow: hidden;
  border-radius: 8px;
}

.cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.article-card:hover .cover img {
  transform: scale(1.05);
}

.content {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-width: 0;
}

.title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: 1.4;
}

.summary {
  font-size: 14px;
  color: #606266;
  line-height: 1.8;
  margin-bottom: 16px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 4;
  -webkit-box-orient: vertical;
}

.meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid #ebeef5;
}

.tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.stats {
  display: flex;
  gap: 16px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
  color: #909399;
}

.footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.author {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #606266;
}

.date {
  font-size: 14px;
  color: #909399;
}

/* 响应式布局 */
@media (max-width: 768px) {
  .card-wrapper {
    flex-direction: column;
  }

  .cover {
    width: 100%;
    height: 200px;
  }

  .title {
    font-size: 20px;
  }

  .summary {
    -webkit-line-clamp: 3;
  }
}
</style>
