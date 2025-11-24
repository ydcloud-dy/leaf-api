<template>
  <div class="about">
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">关于我</h1>
        <p class="page-subtitle">Hello World</p>
      </div>

      <el-card class="about-card" v-loading="loading">
        <div class="about-content">
          <div class="avatar-section">
            <el-avatar :size="150" :src="bloggerInfo.avatar">
              {{ bloggerInfo.nickname.charAt(0) }}
            </el-avatar>
            <h2 class="name">{{ bloggerInfo.nickname }}</h2>
            <p class="bio">{{ bloggerInfo.bio }}</p>
          </div>

          <div class="info-section">
            <h3 class="section-title">个人简介</h3>
            <p class="text">
              {{ bloggerInfo.bio || '暂无简介' }}
            </p>

            <h3 class="section-title" v-if="skills.length > 0">技术栈</h3>
            <div class="skills" v-if="skills.length > 0">
              <el-tag v-for="skill in skills" :key="skill" size="large" effect="plain">
                {{ skill }}
              </el-tag>
            </div>

            <h3 class="section-title" v-if="Object.keys(contacts).length > 0">联系方式</h3>
            <div class="contacts" v-if="Object.keys(contacts).length > 0">
              <div class="contact-item" v-if="contacts.email">
                <el-icon><Message /></el-icon>
                <span>Email: {{ contacts.email }}</span>
              </div>
              <div class="contact-item" v-if="contacts.github">
                <el-icon><Link /></el-icon>
                <span>GitHub: {{ contacts.github }}</span>
              </div>
              <div class="contact-item" v-if="contacts.wechat">
                <el-icon><ChatDotRound /></el-icon>
                <span>WeChat: {{ contacts.wechat }}</span>
              </div>
            </div>

            <h3 class="section-title">关于本站</h3>
            <p class="text">
              本站采用 Vue 3 + Element Plus 构建前端，后端使用 Go + Gin 框架开发。
              所有文章内容均为原创或注明出处，欢迎转载，但请保留原文链接。
            </p>

            <div class="stats-grid">
              <div class="stat-card">
                <el-icon :size="32" color="#409eff"><Document /></el-icon>
                <div class="stat-number">{{ formatNumber(bloggerInfo.article_count) }}</div>
                <div class="stat-label">文章数</div>
              </div>
              <div class="stat-card">
                <el-icon :size="32" color="#67c23a"><View /></el-icon>
                <div class="stat-number">{{ formatNumber(bloggerInfo.total_views) }}</div>
                <div class="stat-label">访问量</div>
              </div>
              <div class="stat-card">
                <el-icon :size="32" color="#e6a23c"><ChatDotRound /></el-icon>
                <div class="stat-number">{{ formatNumber(bloggerInfo.comment_count) }}</div>
                <div class="stat-label">评论数</div>
              </div>
              <div class="stat-card">
                <el-icon :size="32" color="#f56c6c"><Star /></el-icon>
                <div class="stat-number">{{ formatNumber(bloggerInfo.like_count) }}</div>
                <div class="stat-label">获赞数</div>
              </div>
            </div>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getBloggerInfo } from '@/api/stats'

const loading = ref(false)
const bloggerInfo = ref({
  nickname: '博主',
  avatar: '',
  bio: '全栈开发工程师 / 技术爱好者',
  skills: '',
  contacts: '',
  article_count: 0,
  total_views: 0,
  comment_count: 0,
  like_count: 0
})

// 解析技术栈（逗号分隔的字符串）
const skills = computed(() => {
  if (!bloggerInfo.value.skills) return []
  return bloggerInfo.value.skills
    .split(',')
    .map(s => s.trim())
    .filter(s => s)
})

// 解析联系方式（JSON格式）
const contacts = computed(() => {
  if (!bloggerInfo.value.contacts) return {}
  try {
    return JSON.parse(bloggerInfo.value.contacts)
  } catch (e) {
    return {}
  }
})

// 格式化数字（添加k/w单位）
const formatNumber = (num) => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'w'
  }
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num.toString()
}

// 获取博主信息
const fetchBloggerInfo = async () => {
  loading.value = true
  try {
    const { data } = await getBloggerInfo()
    bloggerInfo.value = {
      ...bloggerInfo.value,
      ...data
    }
  } catch (error) {
    console.error('Failed to fetch blogger info:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchBloggerInfo()
})
</script>

<style scoped>
.about {
  padding: 20px 0;
}

.about-card {
  max-width: 900px;
  margin: 0 auto;
}

.about-content {
  display: flex;
  flex-direction: column;
  gap: 40px;
}

.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  color: #fff;
}

.name {
  font-size: 28px;
  font-weight: 600;
  margin: 20px 0 8px;
}

.bio {
  font-size: 16px;
  opacity: 0.9;
}

.info-section {
  padding: 0 20px;
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
  margin: 30px 0 16px;
  padding-bottom: 12px;
  border-bottom: 2px solid #409eff;
}

.text {
  line-height: 1.8;
  color: #606266;
  font-size: 15px;
}

.skills {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.contacts {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.contact-item {
  display: flex;
  align-items: center;
  gap: 12px;
  color: #606266;
  font-size: 15px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-top: 30px;
}

.stat-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 30px 20px;
  background-color: #f5f7fa;
  border-radius: 8px;
  transition: all 0.3s;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stat-number {
  font-size: 28px;
  font-weight: 700;
  color: #303133;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .name {
    font-size: 24px;
  }
}
</style>
