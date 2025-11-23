<template>
  <div class="article-edit">
    <el-card>
      <template #header>
        <span>{{ isEdit ? '编辑文章' : '创建文章' }}</span>
      </template>

      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入文章标题" />
        </el-form-item>

        <el-form-item label="分类" prop="category_id">
          <el-select v-model="form.category_id" placeholder="请选择分类">
            <el-option
              v-for="item in categories"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="标签" prop="tag_ids">
          <el-select v-model="form.tag_ids" multiple placeholder="请选择标签">
            <el-option
              v-for="item in tags"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="章节">
          <el-select v-model="form.chapter_id" placeholder="请选择章节(可选)" clearable>
            <el-option
              v-for="item in chapters"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            >
              <span>{{ item.name }}</span>
              <span v-if="item.tag" style="float: right; color: #8492a6; font-size: 13px">{{ item.tag.name }}</span>
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="摘要" prop="summary">
          <el-input v-model="form.summary" type="textarea" :rows="3" placeholder="请输入文章摘要" />
        </el-form-item>

        <el-form-item label="封面" prop="cover">
          <div class="cover-uploader-wrapper">
            <el-upload
              class="cover-uploader"
              :action="uploadAction"
              :headers="uploadHeaders"
              :show-file-list="false"
              :on-success="handleCoverSuccess"
              :before-upload="beforeCoverUpload"
              accept="image/*"
            >
              <img v-if="form.cover" :src="form.cover" class="cover-image" />
              <div v-else class="cover-placeholder">
                <el-icon class="cover-icon"><Plus /></el-icon>
                <div class="cover-text">点击上传封面</div>
              </div>
            </el-upload>
            <div class="upload-tip">建议上传 16:9 比例的图片，大小不超过 2MB</div>
          </div>
        </el-form-item>

        <el-form-item label="内容" prop="content_markdown">
          <MarkdownEditor v-model="form.content_markdown" />
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :value="0">草稿</el-radio>
            <el-radio :value="1">发布</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            {{ isEdit ? '更新' : '创建' }}
          </el-button>
          <el-button @click="$router.back()">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getArticle, createArticle, updateArticle } from '@/api/article'
import { getTags, getCategories } from '@/api/taxonomy'
import { getChapters } from '@/api/chapter'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import MarkdownEditor from '@/components/MarkdownEditor.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const formRef = ref()
const submitting = ref(false)
const tags = ref([])
const categories = ref([])
const chapters = ref([])

const isEdit = computed(() => !!route.params.id)

// 上传配置
const uploadAction = computed(() => '/api/files/upload')
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${userStore.token}`
}))

const form = reactive({
  title: '',
  content_markdown: '',
  summary: '',
  cover: '',
  category_id: null,
  tag_ids: [],
  chapter_id: null,
  status: 0
})

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content_markdown: [{ required: true, message: '请输入内容', trigger: 'blur' }],
  category_id: [{ required: true, message: '请选择分类', trigger: 'change' }]
}

const fetchData = async () => {
  const [tagsRes, categoriesRes, chaptersRes] = await Promise.all([
    getTags(),
    getCategories(),
    getChapters()
  ])
  tags.value = tagsRes.data
  categories.value = categoriesRes.data
  chapters.value = chaptersRes.data || []

  if (isEdit.value) {
    const res = await getArticle(route.params.id)
    const article = res.data
    Object.assign(form, {
      title: article.title,
      content_markdown: article.content_markdown,
      summary: article.summary,
      cover: article.cover,
      category_id: article.category_id,
      tag_ids: article.tags?.map(t => t.id) || [],
      chapter_id: article.chapter_id,
      status: article.status
    })
  }
}

// 封面上传前校验
const beforeCoverUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过 2MB!')
    return false
  }
  return true
}

// 封面上传成功
const handleCoverSuccess = (response) => {
  if (response.code === 0) {
    form.cover = response.data.url
    ElMessage.success('封面上传成功')
  } else {
    ElMessage.error(response.message || '上传失败')
  }
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  // 验证分类必填
  if (!form.category_id) {
    ElMessage.warning('请选择分类')
    return
  }

  submitting.value = true
  try {
    // 准备提交数据
    const submitData = {
      title: form.title,
      content_markdown: form.content_markdown,
      summary: form.summary,
      cover: form.cover,
      category_id: form.category_id,
      tag_ids: form.tag_ids || [],
      chapter_id: form.chapter_id || null,
      status: form.status
    }

    if (isEdit.value) {
      await updateArticle(route.params.id, submitData)
      ElMessage.success('更新成功')
    } else {
      await createArticle(submitData)
      ElMessage.success('创建成功')
    }
    router.push('/articles')
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.article-edit {
  max-width: 1400px;
}

.article-edit :deep(.el-form-item__content) {
  line-height: normal;
}

.cover-uploader-wrapper {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.cover-uploader :deep(.el-upload) {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: border-color 0.3s;
  width: 360px;
  height: 202px;
}

.cover-uploader :deep(.el-upload:hover) {
  border-color: #409eff;
}

.cover-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: #8c939d;
}

.cover-icon {
  font-size: 28px;
}

.cover-text {
  font-size: 14px;
}

.cover-image {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
}

.upload-tip {
  font-size: 12px;
  color: #909399;
}
</style>
