<template>
  <div class="article-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>文章列表</span>
          <el-button type="primary" @click="$router.push('/articles/create')">新增文章</el-button>
        </div>
      </template>

      <el-form inline class="search-form">
        <el-form-item>
          <el-input v-model="searchParams.keyword" placeholder="搜索标题" clearable />
        </el-form-item>
        <el-form-item>
          <el-select v-model="searchParams.status" placeholder="状态" clearable>
            <el-option label="草稿" value="0" />
            <el-option label="已发布" value="1" />
            <el-option label="已下架" value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchArticles">搜索</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="articles" v-loading="loading" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
        <el-table-column label="分类" width="120">
          <template #default="{ row }">
            {{ row.category?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="标签" width="150">
          <template #default="{ row }">
            <el-tag v-for="tag in row.tags" :key="tag.id" size="small" style="margin-right: 4px">
              {{ tag.name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusMap[row.status].type">
              {{ statusMap[row.status].label }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="view_count" label="浏览" width="80" />
        <el-table-column prop="like_count" label="点赞" width="80" />
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="$router.push(`/articles/${row.id}/edit`)">编辑</el-button>
            <el-button
              v-if="row.status !== 1"
              size="small"
              type="success"
              @click="handleUpdateStatus(row, 1)"
            >上架</el-button>
            <el-button
              v-if="row.status === 1"
              size="small"
              type="warning"
              @click="handleUpdateStatus(row, 2)"
            >下架</el-button>
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
        @size-change="fetchArticles"
        @current-change="fetchArticles"
        style="margin-top: 20px"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getArticles, updateArticleStatus, deleteArticle } from '@/api/article'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const articles = ref([])

const statusMap = {
  0: { label: '草稿', type: 'info' },
  1: { label: '已发布', type: 'success' },
  2: { label: '已下架', type: 'warning' }
}

const searchParams = reactive({
  keyword: '',
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

const fetchArticles = async () => {
  loading.value = true
  try {
    const res = await getArticles({
      page: pagination.page,
      limit: pagination.limit,
      ...searchParams
    })
    articles.value = res.data.list
    pagination.total = res.data.total
  } finally {
    loading.value = false
  }
}

const handleUpdateStatus = async (row, status) => {
  const action = status === 1 ? '上架' : '下架'
  await ElMessageBox.confirm(`确定要${action}文章 "${row.title}" 吗？`, '提示', {
    type: 'warning'
  })
  await updateArticleStatus(row.id, status)
  ElMessage.success(`${action}成功`)
  fetchArticles()
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm(`确定要删除文章 "${row.title}" 吗？`, '提示', {
    type: 'warning'
  })
  await deleteArticle(row.id)
  ElMessage.success('删除成功')
  fetchArticles()
}

onMounted(() => {
  fetchArticles()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}
</style>
