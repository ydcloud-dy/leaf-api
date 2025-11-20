<template>
  <div class="guestbook">
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">留言板</h1>
        <p class="page-subtitle">欢迎留下您的足迹</p>
      </div>

      <el-card class="guestbook-card">
        <!-- 发表留言 -->
        <div v-if="userStore.isLoggedIn" class="message-editor">
          <h3 class="editor-title">发表留言</h3>
          <el-input
            v-model="messageContent"
            type="textarea"
            :rows="6"
            placeholder="写下你想说的话..."
            maxlength="500"
            show-word-limit
          />
          <div class="editor-actions">
            <el-button type="primary" :loading="submitting" @click="handleSubmit">
              发表留言
            </el-button>
          </div>
        </div>

        <div v-else class="login-tip">
          <el-alert
            title="请先登录后发表留言"
            type="info"
            :closable="false"
            show-icon
          >
            <template #default>
              <el-button type="primary" size="small" @click="router.push('/login')">
                去登录
              </el-button>
            </template>
          </el-alert>
        </div>

        <!-- 留言列表 -->
        <div class="messages-section">
          <h3 class="section-title">
            全部留言 ({{ total }})
          </h3>

          <div v-loading="loading" class="messages-list">
            <div v-for="message in messages" :key="message.id" class="message-item">
              <div class="message-avatar">
                <el-avatar :size="50" :src="message.user?.avatar">
                  {{ message.user?.username?.charAt(0).toUpperCase() }}
                </el-avatar>
              </div>

              <div class="message-content">
                <div class="message-header">
                  <span class="username">{{ message.user?.username || '匿名用户' }}</span>
                  <span class="time">{{ formatDate(message.created_at) }}</span>
                </div>

                <div class="message-text">{{ message.content }}</div>
              </div>
            </div>

            <el-empty v-if="!loading && !messages.length" description="暂无留言，快来抢沙发吧！" />
          </div>

          <!-- 分页 -->
          <div v-if="total > pageSize" class="pagination">
            <el-pagination
              v-model:current-page="currentPage"
              :page-size="pageSize"
              :total="total"
              layout="prev, pager, next"
              @current-change="handlePageChange"
            />
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getGuestbookMessages, createGuestbookMessage } from '@/api/comment'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const router = useRouter()
const userStore = useUserStore()

const messages = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const submitting = ref(false)
const messageContent = ref('')

onMounted(() => {
  fetchMessages()
})

const fetchMessages = async () => {
  loading.value = true
  try {
    const res = await getGuestbookMessages({
      page: currentPage.value,
      limit: pageSize.value
    })
    messages.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (error) {
    console.error('Failed to fetch messages:', error)
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  if (!messageContent.value.trim()) {
    ElMessage.warning('请输入留言内容')
    return
  }

  submitting.value = true
  try {
    await createGuestbookMessage({
      content: messageContent.value
    })
    ElMessage.success('留言发表成功')
    messageContent.value = ''
    currentPage.value = 1
    await fetchMessages()
  } catch (error) {
    console.error('Failed to submit message:', error)
    ElMessage.error('留言发表失败')
  } finally {
    submitting.value = false
  }
}

const handlePageChange = () => {
  fetchMessages()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const formatDate = (date) => {
  return dayjs(date).fromNow()
}
</script>

<style scoped>
.guestbook {
  padding: 20px 0;
}

.guestbook-card {
  max-width: 900px;
  margin: 0 auto;
}

.message-editor {
  padding-bottom: 30px;
  margin-bottom: 30px;
  border-bottom: 1px solid #ebeef5;
}

.editor-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 16px;
}

.editor-actions {
  margin-top: 12px;
  text-align: right;
}

.login-tip {
  margin-bottom: 30px;
}

.messages-section {
  min-height: 400px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 20px;
}

.messages-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.message-item {
  display: flex;
  gap: 16px;
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 8px;
  transition: all 0.3s;
}

.message-item:hover {
  background-color: #ecf5ff;
}

.message-content {
  flex: 1;
}

.message-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.username {
  font-weight: 500;
  color: #303133;
  font-size: 16px;
}

.time {
  font-size: 13px;
  color: #909399;
}

.message-text {
  color: #606266;
  line-height: 1.8;
  font-size: 15px;
  margin-bottom: 12px;
  white-space: pre-wrap;
}

.message-actions {
  display: flex;
  gap: 16px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}
</style>
