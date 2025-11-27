<template>
  <div class="chapter-container">
    <el-card>
      <template #header>
        <div class="header">
          <span>章节管理</span>
          <el-button type="primary" @click="showDialog()">新建章节</el-button>
        </div>
      </template>

      <!-- 筛选 -->
      <el-form :inline="true" class="filter-form">
        <el-form-item label="标签">
          <el-select v-model="filters.tag_id" placeholder="选择标签" clearable @change="fetchChapters">
            <el-option
              v-for="tag in tags"
              :key="tag.id"
              :label="tag.name"
              :value="tag.id"
            />
          </el-select>
        </el-form-item>
      </el-form>

      <!-- 表格 -->
      <el-table :data="chapters" v-loading="loading" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="章节名称" />
        <el-table-column label="父章节" width="150">
          <template #default="{ row }">
            <el-tag v-if="row.parent_id" type="info">{{ getParentChapterName(row.parent_id) }}</el-tag>
            <span v-else style="color: #909399;">一级章节</span>
          </template>
        </el-table-column>
        <el-table-column prop="tag" label="所属标签">
          <template #default="{ row }">
            <el-tag v-if="row.tag">{{ row.tag.name }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="100" />
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="showDialog(row)">编辑</el-button>
            <el-button type="danger" size="small" @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新建/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="editingId ? '编辑章节' : '新建章节'"
      width="500px"
    >
      <el-form :model="form" label-width="80px">
        <el-form-item label="标签" required>
          <el-select v-model="form.tag_id" placeholder="请选择标签" @change="handleTagChange">
            <el-option
              v-for="tag in tags"
              :key="tag.id"
              :label="tag.name"
              :value="tag.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="父章节">
          <el-select v-model="form.parent_id" placeholder="不选则为一级章节" clearable>
            <el-option
              v-for="chapter in availableParentChapters"
              :key="chapter.id"
              :label="chapter.name"
              :value="chapter.id"
            />
          </el-select>
          <div style="color: #909399; font-size: 12px; margin-top: 5px;">
            不选择父章节时,将创建为一级章节
          </div>
        </el-form-item>
        <el-form-item label="章节名称" required>
          <el-input v-model="form.name" placeholder="请输入章节名称" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" placeholder="数字越小越靠前" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const editingId = ref(null)
const chapters = ref([])
const tags = ref([])
const filters = ref({
  tag_id: ''
})
const form = ref({
  tag_id: '',
  parent_id: null,
  name: '',
  sort: 0
})

onMounted(() => {
  fetchTags()
  fetchChapters()
})

const fetchTags = async () => {
  try {
    const { data } = await request.get('/tags')
    tags.value = data || []
  } catch (error) {
    ElMessage.error('获取标签列表失败')
  }
}

const fetchChapters = async () => {
  loading.value = true
  try {
    const params = {}
    if (filters.value.tag_id) {
      params.tag_id = filters.value.tag_id
    }
    const { data } = await request.get('/chapters', { params })
    chapters.value = data || []
  } catch (error) {
    ElMessage.error('获取章节列表失败')
  } finally {
    loading.value = false
  }
}

// 可选的父章节列表（只显示同标签下的一级章节，且排除自己）
const availableParentChapters = computed(() => {
  if (!form.value.tag_id) return []

  return chapters.value.filter(chapter => {
    // 只显示同标签下的章节
    if (chapter.tag_id !== form.value.tag_id) return false

    // 排除自己
    if (editingId.value && chapter.id === editingId.value) return false

    // 只显示一级章节（没有父章节的）
    return !chapter.parent_id
  })
})

// 获取父章节名称
const getParentChapterName = (parentId) => {
  const parent = chapters.value.find(c => c.id === parentId)
  return parent ? parent.name : ''
}

// 标签改变时，清空父章节选择
const handleTagChange = () => {
  form.value.parent_id = null
}

const showDialog = (row = null) => {
  if (row) {
    editingId.value = row.id
    form.value = {
      tag_id: row.tag_id,
      parent_id: row.parent_id || null,
      name: row.name,
      sort: row.sort || 0
    }
  } else {
    editingId.value = null
    form.value = {
      tag_id: '',
      parent_id: null,
      name: '',
      sort: 0
    }
  }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!form.value.tag_id || !form.value.name) {
    ElMessage.warning('请填写必填项')
    return
  }

  submitting.value = true
  try {
    const submitData = {
      tag_id: form.value.tag_id,
      name: form.value.name,
      sort: form.value.sort
    }

    // 只有选择了父章节才添加 parent_id
    if (form.value.parent_id) {
      submitData.parent_id = form.value.parent_id
    }

    if (editingId.value) {
      await request.put(`/chapters/${editingId.value}`, submitData)
      ElMessage.success('更新成功')
    } else {
      await request.post('/chapters', submitData)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchChapters()
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除这个章节吗?', '提示', {
      type: 'warning'
    })
    await request.delete(`/chapters/${id}`)
    ElMessage.success('删除成功')
    fetchChapters()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}
</script>

<style scoped>
.chapter-container {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filter-form {
  margin-bottom: 20px;
}
</style>
