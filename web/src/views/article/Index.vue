<template>
  <div class="article-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>文章列表</span>
          <div class="header-actions">
            <el-button type="success" @click="importDialogVisible = true">批量导入</el-button>
            <el-button type="primary" @click="$router.push('/articles/create')">新增文章</el-button>
          </div>
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

      <div v-if="selectedArticles.length > 0" style="margin-bottom: 15px">
        <el-alert type="info" :closable="false">
          已选择 {{ selectedArticles.length }} 篇文章
          <el-button size="small" type="success" @click="handleBatchPublish" style="margin-left: 10px">
            批量上架
          </el-button>
          <el-button size="small" type="warning" @click="handleBatchOffline" style="margin-left: 10px">
            批量下架
          </el-button>
          <el-button size="small" type="primary" @click="batchUpdateDialogVisible = true" style="margin-left: 10px">
            批量更新
          </el-button>
          <el-button size="small" type="danger" @click="handleBatchDelete" style="margin-left: 10px">
            批量删除
          </el-button>
        </el-alert>
      </div>

      <el-table :data="articles" v-loading="loading" border @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
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

    <!-- 批量导入对话框 -->
    <el-dialog v-model="importDialogVisible" title="批量导入 Markdown 文件" width="600px">
      <el-alert
        title="导入说明"
        type="info"
        :closable="false"
        style="margin-bottom: 20px"
      >
        <p>1. 支持批量上传 .md 或 .markdown 文件</p>
        <p>2. 文件名将作为文章标题</p>
        <p>3. 文件内容将保持 Markdown 格式存储</p>
        <p>4. 导入的文章默认为草稿状态</p>
      </el-alert>

      <el-upload
        ref="uploadRef"
        :auto-upload="false"
        :on-change="handleFileChange"
        :file-list="fileList"
        accept=".md,.markdown"
        multiple
        drag
      >
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          将文件拖到此处，或<em>点击上传</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            只能上传 .md 或 .markdown 文件
          </div>
        </template>
      </el-upload>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="importDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleImport" :loading="importing">
            确定导入
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 批量更新对话框 -->
    <el-dialog v-model="batchUpdateDialogVisible" title="批量更新" width="600px">
      <div style="margin-bottom: 20px">
        <p>将为选中的 {{ selectedArticles.length }} 篇文章批量更新以下字段（留空则不更新）</p>
      </div>

      <el-form label-width="80px">
        <el-form-item label="封面图片">
          <div style="display: flex; gap: 10px">
            <el-input v-model="batchFormData.cover" placeholder="请输入图片URL或选择已有图片" clearable />
            <el-button @click="coverSelectorVisible = true">选择图片</el-button>
          </div>
          <div v-if="batchFormData.cover" style="margin-top: 10px">
            <el-image :src="batchFormData.cover" style="width: 200px; height: 120px; object-fit: cover" fit="cover" />
          </div>
        </el-form-item>

        <el-form-item label="分类">
          <el-select v-model="batchFormData.categoryId" placeholder="选择分类" clearable style="width: 100%">
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="章节">
          <el-select v-model="batchFormData.chapterId" placeholder="选择章节" clearable style="width: 100%">
            <el-option label="无章节" :value="null" />
            <el-option
              v-for="chapter in chapters"
              :key="chapter.id"
              :label="chapter.name"
              :value="chapter.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="标签">
          <el-select v-model="batchFormData.tagIds" placeholder="选择标签" multiple clearable style="width: 100%">
            <el-option
              v-for="tag in tags"
              :key="tag.id"
              :label="tag.name"
              :value="tag.id"
            />
          </el-select>
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="batchUpdateDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleBatchUpdate" :loading="batchUpdating">
            确定更新
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 图片选择器对话框 -->
    <el-dialog v-model="coverSelectorVisible" title="选择图片" width="800px">
      <div style="max-height: 400px; overflow-y: auto">
        <el-row :gutter="10">
          <el-col v-for="file in fileList2" :key="file.id" :span="6">
            <div
              class="file-item"
              :class="{ selected: batchFormData.cover === file.url }"
              @click="selectCoverImage(file.url)"
            >
              <el-image :src="file.url" fit="cover" style="width: 100%; height: 120px" />
              <div class="file-name">{{ file.name }}</div>
            </div>
          </el-col>
        </el-row>
      </div>
      <template #footer>
        <el-button @click="coverSelectorVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getArticles, updateArticleStatus, deleteArticle } from '@/api/article'
import { getFiles } from '@/api/file'
import { getCategories, getTags } from '@/api/taxonomy'
import { getChapters } from '@/api/chapter'
import { ElMessage, ElMessageBox } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(false)
const articles = ref([])
const importDialogVisible = ref(false)
const importing = ref(false)
const fileList = ref([])
const uploadRef = ref(null)
const selectedArticles = ref([])
const batchUpdateDialogVisible = ref(false)
const batchFormData = reactive({
  cover: '',
  categoryId: null,
  chapterId: null,
  tagIds: []
})
const batchUpdating = ref(false)
const coverSelectorVisible = ref(false)
const fileList2 = ref([])
const categories = ref([])
const tags = ref([])
const chapters = ref([])

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

const handleFileChange = (file, files) => {
  fileList.value = files
}

const handleImport = async () => {
  if (fileList.value.length === 0) {
    ElMessage.warning('请选择要导入的 Markdown 文件')
    return
  }

  importing.value = true
  try {
    const formData = new FormData()
    fileList.value.forEach((file) => {
      formData.append('files', file.raw)
    })

    const res = await request.post('/articles/import', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })

    console.log('导入响应:', res)

    // request 拦截器已经返回 response.data，所以直接访问 res.data
    const result = res.data
    if (result.success > 0) {
      ElMessage.success(`成功导入 ${result.success} 篇文章`)
    }

    // 如果有失败的文件，显示警告信息
    if (result.failed > 0) {
      const failedList = result.failed_files?.join('\n') || ''
      ElMessage.warning({
        message: `${result.failed} 篇文章导入失败${failedList ? ':\n' + failedList : ''}`,
        duration: 5000,
        showClose: true
      })
    }

    importDialogVisible.value = false
    fileList.value = []
    fetchArticles()
  } catch (error) {
    console.error('导入错误:', error)
    console.error('错误响应:', error.response)
    const errorMsg = error.response?.data?.message || error.message || '导入失败'
    ElMessage.error({
      message: errorMsg,
      duration: 5000,
      showClose: true
    })
  } finally {
    importing.value = false
  }
}

const handleSelectionChange = (selection) => {
  selectedArticles.value = selection
}

const handleBatchUpdate = async () => {
  if (selectedArticles.value.length === 0) {
    ElMessage.warning('请先选择要更新的文章')
    return
  }

  // 检查是否至少有一个字段要更新
  const hasUpdate = batchFormData.cover ||
                    batchFormData.categoryId ||
                    batchFormData.chapterId !== null ||
                    batchFormData.tagIds.length > 0

  if (!hasUpdate) {
    ElMessage.warning('请至少填写一个要更新的字段')
    return
  }

  batchUpdating.value = true
  try {
    const payload = {
      article_ids: selectedArticles.value.map(a => a.id)
    }

    if (batchFormData.cover) {
      payload.cover = batchFormData.cover
    }
    if (batchFormData.categoryId) {
      payload.category_id = batchFormData.categoryId
    }
    if (batchFormData.chapterId !== null) {
      payload.chapter_id = batchFormData.chapterId
    }
    if (batchFormData.tagIds.length > 0) {
      payload.tag_ids = batchFormData.tagIds
    }

    await request.post('/articles/batch-update-fields', payload)

    ElMessage.success(`成功更新 ${selectedArticles.value.length} 篇文章`)
    batchUpdateDialogVisible.value = false

    // 重置表单
    batchFormData.cover = ''
    batchFormData.categoryId = null
    batchFormData.chapterId = null
    batchFormData.tagIds = []

    selectedArticles.value = []
    fetchArticles()
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '批量更新失败')
  } finally {
    batchUpdating.value = false
  }
}

const selectCoverImage = (url) => {
  batchFormData.cover = url
  coverSelectorVisible.value = false
}

const fetchFiles = async () => {
  try {
    const res = await getFiles({ page: 1, page_size: 100 })
    fileList2.value = res.data.list || []
  } catch (error) {
    console.error('获取文件列表失败:', error)
  }
}

const fetchCategories = async () => {
  try {
    const res = await getCategories()
    categories.value = res.data || []
  } catch (error) {
    console.error('获取分类列表失败:', error)
  }
}

const fetchTags = async () => {
  try {
    const res = await getTags()
    tags.value = res.data || []
  } catch (error) {
    console.error('获取标签列表失败:', error)
  }
}

const fetchChapters = async () => {
  try {
    const res = await getChapters()
    // 后端返回的是 res.data 直接就是数组,不需要 res.data.data
    chapters.value = res.data || []
    console.log('获取到的章节列表:', chapters.value)
  } catch (error) {
    console.error('获取章节列表失败:', error)
  }
}

const handleBatchDelete = async () => {
  if (selectedArticles.value.length === 0) {
    ElMessage.warning('请先选择要删除的文章')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedArticles.value.length} 篇文章吗？此操作不可恢复！`,
      '批量删除确认',
      {
        type: 'warning',
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        confirmButtonClass: 'el-button--danger'
      }
    )

    await request.post('/articles/batch-delete', {
      article_ids: selectedArticles.value.map(a => a.id)
    })

    ElMessage.success(`成功删除 ${selectedArticles.value.length} 篇文章`)
    selectedArticles.value = []
    fetchArticles()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '批量删除失败')
    }
  }
}

const handleBatchPublish = async () => {
  if (selectedArticles.value.length === 0) {
    ElMessage.warning('请先选择要上架的文章')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要上架选中的 ${selectedArticles.value.length} 篇文章吗？`,
      '批量上架确认',
      {
        type: 'success',
        confirmButtonText: '确定上架',
        cancelButtonText: '取消'
      }
    )

    // 使用批量更新字段接口,只更新状态
    const payload = {
      article_ids: selectedArticles.value.map(a => a.id)
    }

    // 批量更新状态为已发布
    for (const article of selectedArticles.value) {
      await request.patch(`/articles/${article.id}/status`, { status: 1 })
    }

    ElMessage.success(`成功上架 ${selectedArticles.value.length} 篇文章`)
    selectedArticles.value = []
    fetchArticles()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '批量上架失败')
    }
  }
}

const handleBatchOffline = async () => {
  if (selectedArticles.value.length === 0) {
    ElMessage.warning('请先选择要下架的文章')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要下架选中的 ${selectedArticles.value.length} 篇文章吗？`,
      '批量下架确认',
      {
        type: 'warning',
        confirmButtonText: '确定下架',
        cancelButtonText: '取消'
      }
    )

    // 批量更新状态为已下架
    for (const article of selectedArticles.value) {
      await request.patch(`/articles/${article.id}/status`, { status: 2 })
    }

    ElMessage.success(`成功下架 ${selectedArticles.value.length} 篇文章`)
    selectedArticles.value = []
    fetchArticles()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '批量下架失败')
    }
  }
}

onMounted(() => {
  fetchArticles()
  fetchFiles()
  fetchCategories()
  fetchTags()
  fetchChapters()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.search-form {
  margin-bottom: 20px;
}

.file-item {
  cursor: pointer;
  border: 2px solid transparent;
  border-radius: 4px;
  padding: 5px;
  margin-bottom: 10px;
  transition: all 0.3s;
}

.file-item:hover {
  border-color: #409eff;
}

.file-item.selected {
  border-color: #67c23a;
  background-color: #f0f9ff;
}

.file-name {
  font-size: 12px;
  text-align: center;
  margin-top: 5px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
