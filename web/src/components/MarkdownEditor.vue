<template>
  <div class="markdown-editor-wrapper">
    <MdEditor
      v-model="content"
      :preview="true"
      :toolbars="toolbars"
      :toolbarsExclude="toolbarsExclude"
      :showCodeRowNumber="true"
      :previewTheme="previewTheme"
      :codeTheme="codeTheme"
      :theme="theme"
      :language="language"
      :placeholder="placeholder"
      @onUploadImg="handleUploadImg"
      @onChange="handleChange"
      :style="{ height: editorHeight }"
      class="md-editor"
    />
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { uploadFile } from '@/api/file'
import { ElMessage } from 'element-plus'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  height: {
    type: String,
    default: '600px'
  },
  placeholder: {
    type: String,
    default: '请输入 Markdown 内容...'
  },
  theme: {
    type: String,
    default: 'light' // light / dark
  },
  previewTheme: {
    type: String,
    default: 'default' // default / github / vuepress
  },
  codeTheme: {
    type: String,
    default: 'atom' // atom / a11y / github / gradient / kimbie / paraiso / qtcreator / stackoverflow
  },
  language: {
    type: String,
    default: 'zh-CN' // zh-CN / en-US
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

const content = ref(props.modelValue)
const editorHeight = ref(props.height)

// 工具栏配置（保留所有默认工具）
const toolbars = [
  'bold',
  'underline',
  'italic',
  'strikeThrough',
  '-',
  'title',
  'sub',
  'sup',
  'quote',
  'unorderedList',
  'orderedList',
  'task',
  '-',
  'codeRow',
  'code',
  'link',
  'image',
  'table',
  'mermaid',
  'katex',
  '-',
  'revoke',
  'next',
  'save',
  '=',
  'pageFullscreen',
  'fullscreen',
  'preview',
  'previewOnly',
  'htmlPreview',
  'catalog',
  'github'
]

// 排除的工具栏（如果需要可以排除一些）
const toolbarsExclude = []

// 监听外部值变化
watch(() => props.modelValue, (newValue) => {
  if (content.value !== newValue) {
    content.value = newValue
  }
})

// 监听内部值变化
watch(content, (newValue) => {
  emit('update:modelValue', newValue)
})

// 处理内容变化
const handleChange = (value) => {
  emit('change', value)
}

// 处理图片上传
const handleUploadImg = async (files, callback) => {
  try {
    // 可以同时上传多个图片
    const uploadPromises = Array.from(files).map(async (file) => {
      const res = await uploadFile(file, 'articles')
      return res.data.url
    })

    const urls = await Promise.all(uploadPromises)

    // 回调返回图片 URL 数组
    callback(urls)

    ElMessage.success('图片上传成功')
  } catch (error) {
    ElMessage.error('图片上传失败：' + (error.message || '未知错误'))
    callback([])
  }
}
</script>

<style scoped>
.markdown-editor-wrapper {
  width: 100%;
}

.markdown-editor-wrapper :deep(.md-editor) {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
}

.markdown-editor-wrapper :deep(.md-editor-dark) {
  border-color: #4c4d4f;
}

/* 自定义样式 */
.markdown-editor-wrapper :deep(.md-editor-toolbar) {
  border-bottom: 1px solid #dcdfe6;
  background: #fafafa;
}

.markdown-editor-wrapper :deep(.md-editor-dark .md-editor-toolbar) {
  background: #1e1e1e;
  border-bottom-color: #4c4d4f;
}

.markdown-editor-wrapper :deep(.md-editor-preview-wrapper) {
  padding: 16px;
}

.markdown-editor-wrapper :deep(.md-editor-catalog-wrapper) {
  padding: 8px;
}

/* 滚动条样式 */
.markdown-editor-wrapper :deep(.cm-scroller::-webkit-scrollbar),
.markdown-editor-wrapper :deep(.md-editor-preview::-webkit-scrollbar),
.markdown-editor-wrapper :deep(.md-editor-catalog::-webkit-scrollbar) {
  width: 6px;
  height: 6px;
}

.markdown-editor-wrapper :deep(.cm-scroller::-webkit-scrollbar-thumb),
.markdown-editor-wrapper :deep(.md-editor-preview::-webkit-scrollbar-thumb),
.markdown-editor-wrapper :deep(.md-editor-catalog::-webkit-scrollbar-thumb) {
  background: #dcdfe6;
  border-radius: 3px;
}

.markdown-editor-wrapper :deep(.cm-scroller::-webkit-scrollbar-thumb:hover),
.markdown-editor-wrapper :deep(.md-editor-preview::-webkit-scrollbar-thumb:hover),
.markdown-editor-wrapper :deep(.md-editor-catalog::-webkit-scrollbar-thumb:hover) {
  background: #c0c4cc;
}
</style>
