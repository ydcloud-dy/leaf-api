<template>
  <div class="comment-section">
    <h3 class="section-title">
      评论 ({{ total }})
    </h3>

    <!-- 发表评论 -->
    <div v-if="userStore.isLoggedIn" class="comment-editor">
      <el-input
        v-model="commentContent"
        type="textarea"
        :rows="4"
        placeholder="写下你的评论..."
        maxlength="500"
        show-word-limit
      />
      <div class="editor-actions">
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          发表评论
        </el-button>
      </div>
    </div>

    <div v-else class="login-tip">
      <el-alert
        title="请先登录后发表评论"
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

    <!-- 评论列表 -->
    <div v-loading="loading" class="comment-list">
      <div v-for="comment in comments" :key="comment.id" class="comment-item">
        <div class="comment-avatar">
          <el-avatar :size="40" :src="comment.user?.avatar">
            {{ comment.user?.username?.charAt(0).toUpperCase() }}
          </el-avatar>
        </div>

        <div class="comment-content">
          <div class="comment-header">
            <span class="username">{{ comment.user?.username || '匿名用户' }}</span>
            <span class="time">{{ formatDate(comment.created_at) }}</span>
          </div>

          <div class="comment-text">{{ comment.content }}</div>

          <div class="comment-actions">
            <el-button
              text
              size="small"
              @click="handleLikeComment(comment)"
              :type="comment.is_liked ? 'primary' : ''"
            >
              <el-icon><Promotion /></el-icon>
              {{ comment.like_count || 0 }}
            </el-button>
            <el-button
              text
              size="small"
              @click="handleReply(comment)"
            >
              <el-icon><ChatDotRound /></el-icon>
              回复
            </el-button>
            <el-button
              v-if="canDelete(comment)"
              text
              size="small"
              type="danger"
              @click="handleDelete(comment)"
            >
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </div>

          <!-- 回复输入框 -->
          <div v-if="replyingTo === comment.id" class="reply-editor">
            <el-input
              v-model="replyContent"
              type="textarea"
              :rows="3"
              :placeholder="`回复 @${comment.user?.username}...`"
              maxlength="300"
            />
            <div class="reply-actions">
              <el-button size="small" @click="cancelReply">取消</el-button>
              <el-button
                size="small"
                type="primary"
                :loading="submitting"
                @click="submitReply(comment)"
              >
                回复
              </el-button>
            </div>
          </div>

          <!-- 回复列表 -->
          <div v-if="comment.replies?.length" class="replies">
            <div v-for="reply in comment.replies" :key="reply.id" class="reply-item">
              <el-avatar :size="32" :src="reply.user?.avatar">
                {{ reply.user?.username?.charAt(0).toUpperCase() }}
              </el-avatar>
              <div class="reply-content">
                <div class="reply-header">
                  <span class="username">{{ reply.user?.username }}</span>
                  <span v-if="reply.reply_to_user" class="reply-to">
                    <el-icon><Right /></el-icon>
                    <span>{{ reply.reply_to_user.username }}</span>
                  </span>
                  <span class="time">{{ formatDate(reply.created_at) }}</span>
                </div>
                <div class="reply-text">{{ reply.content }}</div>
                <div class="comment-actions">
                  <el-button
                    text
                    size="small"
                    @click="handleLikeComment(reply)"
                    :type="reply.is_liked ? 'primary' : ''"
                  >
                    <el-icon><Promotion /></el-icon>
                    {{ reply.like_count || 0 }}
                  </el-button>
                  <el-button
                    text
                    size="small"
                    @click="handleReplyToReply(comment, reply)"
                  >
                    <el-icon><ChatDotRound /></el-icon>
                    回复
                  </el-button>
                  <el-button
                    v-if="canDeleteReply(comment, reply)"
                    text
                    size="small"
                    type="danger"
                    @click="handleDelete(reply, comment)"
                  >
                    <el-icon><Delete /></el-icon>
                    删除
                  </el-button>
                </div>

                <!-- 子回复输入框 -->
                <div v-if="replyingTo === reply.id" class="reply-editor">
                  <el-input
                    v-model="replyContent"
                    type="textarea"
                    :rows="3"
                    :placeholder="`回复 @${reply.user?.username}...`"
                    maxlength="300"
                  />
                  <div class="reply-actions">
                    <el-button size="small" @click="cancelReply">取消</el-button>
                    <el-button
                      size="small"
                      type="primary"
                      :loading="submitting"
                      @click="submitReply(comment, reply)"
                    >
                      回复
                    </el-button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <el-empty v-if="!loading && !comments.length" description="暂无评论" />

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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getArticleComments, createComment, likeComment, unlikeComment, deleteComment } from '@/api/comment'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Promotion, ChatDotRound, Delete, Right } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const props = defineProps({
  articleId: {
    type: [Number, String],
    required: true
  }
})

const router = useRouter()
const userStore = useUserStore()

const comments = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const submitting = ref(false)
const commentContent = ref('')
const replyContent = ref('')
const replyingTo = ref(null)

onMounted(() => {
  fetchComments()
})

const fetchComments = async () => {
  loading.value = true
  try {
    const res = await getArticleComments(props.articleId, {
      page: currentPage.value,
      limit: pageSize.value
    })
    comments.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (error) {
    console.error('Failed to fetch comments:', error)
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  if (!commentContent.value.trim()) {
    ElMessage.warning('请输入评论内容')
    return
  }

  submitting.value = true
  try {
    await createComment({
      article_id: props.articleId,
      content: commentContent.value
    })
    ElMessage.success('评论发表成功')
    commentContent.value = ''
    currentPage.value = 1
    await fetchComments()
  } catch (error) {
    console.error('Failed to submit comment:', error)
  } finally {
    submitting.value = false
  }
}

const handleReply = (comment) => {
  replyingTo.value = comment.id
  replyContent.value = ''
}

const handleReplyToReply = (parentComment, reply) => {
  replyingTo.value = reply.id
  replyContent.value = ''
}

const cancelReply = () => {
  replyingTo.value = null
  replyContent.value = ''
}

const submitReply = async (parentComment, targetReply = null) => {
  if (!replyContent.value.trim()) {
    ElMessage.warning('请输入回复内容')
    return
  }

  submitting.value = true
  try {
    await createComment({
      article_id: props.articleId,
      parent_id: parentComment.id,
      reply_to_user_id: targetReply ? targetReply.user_id : parentComment.user_id,
      content: replyContent.value
    })
    ElMessage.success('回复成功')
    cancelReply()
    await fetchComments()
  } catch (error) {
    console.error('Failed to reply:', error)
    ElMessage.error('回复失败')
  } finally {
    submitting.value = false
  }
}

const handleLikeComment = async (comment) => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }

  try {
    if (comment.is_liked) {
      await unlikeComment(comment.id)
      comment.like_count = (comment.like_count || 0) - 1
      comment.is_liked = false
    } else {
      await likeComment(comment.id)
      comment.like_count = (comment.like_count || 0) + 1
      comment.is_liked = true
    }
  } catch (error) {
    console.error('Failed to like comment:', error)
    ElMessage.error('操作失败')
  }
}

const canDelete = (comment) => {
  if (!userStore.isLoggedIn) return false
  return comment.user_id === userStore.user?.id
}

const canDeleteReply = (parentComment, reply) => {
  if (!userStore.isLoggedIn) return false
  return reply.user_id === userStore.user?.id || parentComment.user_id === userStore.user?.id
}

const handleDelete = async (comment, parentComment = null) => {
  const message = parentComment
    ? '确定要删除这条回复吗？'
    : comment.replies?.length
      ? '删除父评论将同时删除所有子评论，确定要删除吗？'
      : '确定要删除这条评论吗？'

  try {
    await ElMessageBox.confirm(message, '提示', {
      type: 'warning',
      confirmButtonText: '确定',
      cancelButtonText: '取消'
    })

    await deleteComment(comment.id)
    ElMessage.success('删除成功')
    await fetchComments()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete comment:', error)
      ElMessage.error('删除失败')
    }
  }
}

const handlePageChange = () => {
  fetchComments()
}

const formatDate = (date) => {
  return dayjs(date).fromNow()
}
</script>

<style scoped>
.comment-section {
  margin-top: 40px;
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 20px;
  color: #303133;
}

.comment-editor {
  margin-bottom: 30px;
}

.editor-actions {
  margin-top: 10px;
  text-align: right;
}

.login-tip {
  margin-bottom: 30px;
}

.comment-list {
  min-height: 200px;
}

.comment-item {
  display: flex;
  gap: 12px;
  padding: 20px 0;
  border-bottom: 1px solid #ebeef5;
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-content {
  flex: 1;
}

.comment-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.username {
  font-weight: 500;
  color: #303133;
}

.time {
  font-size: 13px;
  color: #909399;
}

.comment-text {
  color: #606266;
  line-height: 1.6;
  margin-bottom: 12px;
  white-space: pre-wrap;
}

.comment-actions {
  display: flex;
  gap: 16px;
}

.reply-editor {
  margin-top: 12px;
  padding: 12px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.reply-actions {
  margin-top: 8px;
  text-align: right;
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}

.replies {
  margin-top: 16px;
  padding-left: 20px;
  border-left: 2px solid #ebeef5;
}

.reply-item {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}

.reply-item:last-child {
  margin-bottom: 0;
}

.reply-content {
  flex: 1;
}

.reply-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
  flex-wrap: wrap;
}

.reply-to {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #909399;
  font-size: 13px;
}

.reply-to .el-icon {
  font-size: 12px;
}

.reply-text {
  color: #606266;
  font-size: 14px;
  line-height: 1.6;
  margin-bottom: 8px;
  white-space: pre-wrap;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}
</style>
