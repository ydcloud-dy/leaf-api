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

        <el-form-item label="摘要" prop="summary">
          <el-input v-model="form.summary" type="textarea" :rows="3" placeholder="请输入文章摘要" />
        </el-form-item>

        <el-form-item label="封面" prop="cover">
          <el-input v-model="form.cover" placeholder="封面图片URL" />
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
import { ElMessage } from 'element-plus'
import MarkdownEditor from '@/components/MarkdownEditor.vue'

const route = useRoute()
const router = useRouter()

const formRef = ref()
const submitting = ref(false)
const tags = ref([])
const categories = ref([])

const isEdit = computed(() => !!route.params.id)

const form = reactive({
  title: '',
  content_markdown: '',
  summary: '',
  cover: '',
  category_id: null,
  tag_ids: [],
  status: 0
})

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content_markdown: [{ required: true, message: '请输入内容', trigger: 'blur' }],
  category_id: [{ required: true, message: '请选择分类', trigger: 'change' }]
}

const fetchData = async () => {
  const [tagsRes, categoriesRes] = await Promise.all([
    getTags(),
    getCategories()
  ])
  tags.value = tagsRes.data
  categories.value = categoriesRes.data

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
      status: article.status
    })
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
</style>
