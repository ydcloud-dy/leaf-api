<template>
  <div class="profile-container">
    <el-card>
      <template #header>
        <span>个人信息</span>
      </template>

      <el-form ref="profileFormRef" :model="profileForm" :rules="profileRules" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="profileForm.username" disabled />
        </el-form-item>

        <el-form-item label="邮箱" prop="email">
          <el-input v-model="profileForm.email" />
        </el-form-item>

        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="profileForm.nickname" />
        </el-form-item>

        <el-form-item label="个人简介" prop="bio">
          <el-input v-model="profileForm.bio" type="textarea" :rows="3" maxlength="500" show-word-limit />
        </el-form-item>

        <el-form-item label="技术栈" prop="skills">
          <el-input v-model="profileForm.skills" placeholder="用逗号分隔，如：Go,Python,Vue" />
        </el-form-item>

        <el-form-item label="联系方式" prop="contacts">
          <el-input v-model="profileForm.contacts" type="textarea" :rows="3" placeholder='JSON格式，如：{"email":"xxx@example.com"}' />
        </el-form-item>

        <el-form-item label="设为博主">
          <el-switch v-model="profileForm.isBlogger" />
          <div class="form-tip">开启后，您的信息将显示在博客前台的"关于"页面</div>
        </el-form-item>

        <el-form-item label="头像" prop="avatar">
          <div class="avatar-uploader-wrapper">
            <el-upload
              class="avatar-uploader"
              :action="uploadAction"
              :headers="uploadHeaders"
              :show-file-list="false"
              :on-success="handleAvatarSuccess"
              :before-upload="beforeAvatarUpload"
              accept="image/*"
            >
              <img v-if="profileForm.avatar" :src="profileForm.avatar" class="avatar" />
              <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
            </el-upload>
            <div class="upload-tip">建议上传 200x200 像素的图片</div>
          </div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleUpdateProfile" :loading="profileLoading">
            更新信息
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card style="margin-top: 20px;">
      <template #header>
        <span>修改密码</span>
      </template>

      <el-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" label-width="100px">
        <el-form-item label="新密码" prop="password">
          <el-input v-model="passwordForm.password" type="password" placeholder="请输入新密码" show-password />
        </el-form-item>

        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="passwordForm.confirmPassword" type="password" placeholder="请再次输入新密码" show-password />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleChangePassword" :loading="passwordLoading">
            修改密码
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { updateProfile } from '@/api/auth'
import request from '@/utils/request'

const userStore = useUserStore()

const profileFormRef = ref()
const passwordFormRef = ref()
const profileLoading = ref(false)
const passwordLoading = ref(false)

// 上传配置
const uploadAction = computed(() => '/api/files/upload')
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${userStore.token}`
}))

const profileForm = reactive({
  username: '',
  email: '',
  avatar: '',
  nickname: '',
  bio: '',
  skills: '',
  contacts: '',
  isBlogger: false
})

const passwordForm = reactive({
  password: '',
  confirmPassword: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== passwordForm.password) {
    callback(new Error('两次输入密码不一致'))
  } else {
    callback()
  }
}

const profileRules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ]
}

const passwordRules = {
  password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 50, message: '密码长度在 6 到 50 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const loadProfile = () => {
  profileForm.username = userStore.userInfo.username
  profileForm.email = userStore.userInfo.email
  profileForm.avatar = userStore.userInfo.avatar || ''
  profileForm.nickname = userStore.userInfo.nickname || ''
  profileForm.bio = userStore.userInfo.bio || ''
  profileForm.skills = userStore.userInfo.skills || ''
  profileForm.contacts = userStore.userInfo.contacts || ''
  profileForm.isBlogger = userStore.userInfo.is_blogger || false
}

// 头像上传前校验
const beforeAvatarUpload = (file) => {
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

// 头像上传成功
const handleAvatarSuccess = (response) => {
  if (response.code === 0) {
    profileForm.avatar = response.data.url
    ElMessage.success('头像上传成功')
  } else {
    ElMessage.error(response.message || '上传失败')
  }
}

const handleUpdateProfile = async () => {
  const valid = await profileFormRef.value.validate().catch(() => false)
  if (!valid) return

  profileLoading.value = true
  try {
    const { data } = await updateProfile({
      email: profileForm.email,
      avatar: profileForm.avatar,
      nickname: profileForm.nickname,
      bio: profileForm.bio,
      skills: profileForm.skills,
      contacts: profileForm.contacts,
      is_blogger: profileForm.isBlogger
    })

    // 使用后端返回的完整用户数据更新store
    userStore.userInfo = data
    localStorage.setItem('userInfo', JSON.stringify(data))

    ElMessage.success('个人信息更新成功')
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '更新失败')
  } finally {
    profileLoading.value = false
  }
}

const handleChangePassword = async () => {
  const valid = await passwordFormRef.value.validate().catch(() => false)
  if (!valid) return

  passwordLoading.value = true
  try {
    const userId = userStore.userInfo.id
    await request.put(`/users/${userId}`, {
      password: passwordForm.password
    })

    ElMessage.success('密码修改成功，请重新登录')

    // 清空密码表单
    passwordForm.password = ''
    passwordForm.confirmPassword = ''
    passwordFormRef.value.resetFields()

    // 3秒后自动退出登录
    setTimeout(() => {
      userStore.logout()
      window.location.href = '/login'
    }, 3000)
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '修改密码失败')
  } finally {
    passwordLoading.value = false
  }
}

onMounted(() => {
  loadProfile()
})
</script>

<style scoped>
.profile-container {
  max-width: 600px;
}

.profile-container :deep(.el-card__header) {
  font-weight: bold;
}

.avatar-uploader-wrapper {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.avatar-uploader :deep(.el-upload) {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: border-color 0.3s;
  width: 178px;
  height: 178px;
}

.avatar-uploader :deep(.el-upload:hover) {
  border-color: #409eff;
}

.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar {
  width: 178px;
  height: 178px;
  display: block;
  object-fit: cover;
}

.upload-tip {
  font-size: 12px;
  color: #909399;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-left: 10px;
}
</style>
