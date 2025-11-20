<template>
  <el-card class="article-card" shadow="hover" @click="handleClick">
    <div v-if="article.cover" class="cover">
      <img :src="article.cover" :alt="article.title" />
    </div>

    <div class="content">
      <h3 class="title">{{ article.title }}</h3>

      <p class="summary">{{ article.summary || (article.content_markdown || article.content || '')?.substring(0, 150) + '...' }}</p>

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
            {{ article.views || 0 }}
          </span>
          <span class="stat-item">
            <el-icon><ChatDotRound /></el-icon>
            {{ article.comments || 0 }}
          </span>
          <span class="stat-item">
            <el-icon><Star /></el-icon>
            {{ article.likes || 0 }}
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

.cover {
  width: 100%;
  height: 200px;
  overflow: hidden;
  border-radius: 4px;
  margin-bottom: 16px;
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
  padding: 4px 0;
}

.title {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.summary {
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
  margin-bottom: 16px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 3;
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
</style>
