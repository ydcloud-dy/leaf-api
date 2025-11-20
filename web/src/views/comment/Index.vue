<template>
  <div class="comment-management">
    <el-card>
      <template #header>
        <span>评论列表</span>
      </template>

      <el-form inline class="search-form">
        <el-form-item>
          <el-select v-model="searchParams.status" placeholder="状态" clearable>
            <el-option label="待审核" value="0" />
            <el-option label="已通过" value="1" />
            <el-option label="已拒绝" value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchComments">搜索</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="comments" v-loading="loading" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="content" label="评论内容" min-width="200" show-overflow-tooltip />
        <el-table-column label="用户" width="120">
          <template #default="{ row }">
            {{ row.user?.username || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="文章" width="200">
          <template #default="{ row }">
            {{ row.article?.title || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusMap[row.status].type">
              {{ statusMap[row.status].label }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="评论时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 0"
              size="small"
              type="success"
              @click="handleUpdateStatus(row, 1)"
            >通过</el-button>
            <el-button
              v-if="row.status === 0"
              size="small"
              type="warning"
              @click="handleUpdateStatus(row, 2)"
            >拒绝</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.limit"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @size-change="fetchComments"
        @current-change="fetchComments"
        style="margin-top: 20px"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getComments, updateCommentStatus, deleteComment } from '@/api/comment'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const comments = ref([])

const statusMap = {
  0: { label: '待审核', type: 'warning' },
  1: { label: '已通过', type: 'success' },
  2: { label: '已拒绝', type: 'danger' }
}

const searchParams = reactive({
  status: ''
})

const pagination = reactive({
  page: 1,
  limit: 10,
  total: 0
})

const formatDate = (date) => {
  return new Date(date).toLocaleString('zh-CN')
}

const fetchComments = async () => {
  loading.value = true
  try {
    const res = await getComments({
      page: pagination.page,
      limit: pagination.limit,
      ...searchParams
    })
    comments.value = res.data.list
    pagination.total = res.data.total
  } finally {
    loading.value = false
  }
}

const handleUpdateStatus = async (row, status) => {
  const action = status === 1 ? '通过' : '拒绝'
  await updateCommentStatus(row.id, status)
  ElMessage.success(`已${action}`)
  fetchComments()
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm('确定要删除这条评论吗？', '提示', {
    type: 'warning'
  })
  await deleteComment(row.id)
  ElMessage.success('删除成功')
  fetchComments()
}

onMounted(() => {
  fetchComments()
})
</script>

<style scoped>
.search-form {
  margin-bottom: 20px;
}
</style>
