<template>
  <div class="settings">
    <el-card>
      <template #header>
        <span>系统设置</span>
      </template>

      <el-form
        ref="formRef"
        :model="form"
        label-width="120px"
        v-loading="loading"
        style="max-width: 600px"
      >
        <el-form-item label="站点名称">
          <el-input v-model="form.site_name" placeholder="请输入站点名称" />
        </el-form-item>

        <el-form-item label="站点描述">
          <el-input
            v-model="form.site_description"
            type="textarea"
            :rows="3"
            placeholder="请输入站点描述"
          />
        </el-form-item>

        <el-form-item label="站点Logo">
          <el-input v-model="form.site_logo" placeholder="Logo图片URL" />
        </el-form-item>

        <el-form-item label="备案号">
          <el-input v-model="form.icp_number" placeholder="请输入ICP备案号" />
        </el-form-item>

        <el-form-item label="版权信息">
          <el-input v-model="form.copyright" placeholder="请输入版权信息" />
        </el-form-item>

        <el-form-item label="联系邮箱">
          <el-input v-model="form.contact_email" placeholder="请输入联系邮箱" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            保存设置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getSettings, updateSettings } from '@/api/system'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const submitting = ref(false)

const form = reactive({
  site_name: '',
  site_description: '',
  site_logo: '',
  icp_number: '',
  copyright: '',
  contact_email: ''
})

const fetchSettings = async () => {
  loading.value = true
  try {
    const res = await getSettings()
    Object.assign(form, res.data)
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    await updateSettings(form)
    ElMessage.success('保存成功')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchSettings()
})
</script>

<style scoped>
</style>
