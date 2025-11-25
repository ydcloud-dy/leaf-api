<template>
  <div class="profile">
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">个人中心</h1>
      </div>

      <div class="profile-content">
        <!-- 用户信息卡片 -->
        <el-card class="user-card">
          <div class="user-info">
            <el-avatar :size="100" :src="userStore.avatar">
              {{ userStore.username.charAt(0).toUpperCase() }}
            </el-avatar>
            <div class="user-details">
              <h2 class="username">{{ userStore.username }}</h2>
              <p class="user-role">
                <el-tag :type="userStore.isAdmin ? 'danger' : 'primary'">
                  {{ userStore.isAdmin ? '管理员' : '普通用户' }}
                </el-tag>
              </p>
            </div>
          </div>

          <el-divider />

          <div class="user-stats">
            <div class="stat-item">
              <div class="stat-number">{{ stats.articles }}</div>
              <div class="stat-label">发布文章</div>
            </div>
            <div class="stat-item">
              <div class="stat-number">{{ stats.likes }}</div>
              <div class="stat-label">点赞数</div>
            </div>
            <div class="stat-item">
              <div class="stat-number">{{ stats.favorites }}</div>
              <div class="stat-label">收藏数</div>
            </div>
            <div class="stat-item">
              <div class="stat-number">{{ stats.comments }}</div>
              <div class="stat-label">评论数</div>
            </div>
          </div>
        </el-card>

        <!-- 内容区域 -->
        <el-card class="content-card">
          <el-tabs v-model="activeTab" class="profile-tabs">
            <!-- 我的点赞 -->
            <el-tab-pane label="我的点赞" name="likes">
              <div v-loading="loading" class="tab-content">
                <ArticleCard
                  v-for="article in likedArticles"
                  :key="article.id"
                  :article="article"
                />
                <el-empty v-if="!loading && !likedArticles.length" description="暂无点赞记录" />

                <div v-if="likesTotal > likesPageSize" class="pagination">
                  <el-pagination
                    v-model:current-page="likesPage"
                    :page-size="likesPageSize"
                    :total="likesTotal"
                    layout="prev, pager, next"
                    @current-change="fetchLikedArticles"
                  />
                </div>
              </div>
            </el-tab-pane>

            <!-- 我的收藏 -->
            <el-tab-pane label="我的收藏" name="favorites">
              <div v-loading="loading" class="tab-content">
                <ArticleCard
                  v-for="article in favoritedArticles"
                  :key="article.id"
                  :article="article"
                />
                <el-empty v-if="!loading && !favoritedArticles.length" description="暂无收藏记录" />

                <div v-if="favoritesTotal > favoritesPageSize" class="pagination">
                  <el-pagination
                    v-model:current-page="favoritesPage"
                    :page-size="favoritesPageSize"
                    :total="favoritesTotal"
                    layout="prev, pager, next"
                    @current-change="fetchFavoritedArticles"
                  />
                </div>
              </div>
            </el-tab-pane>

            <!-- 账号设置 -->
            <el-tab-pane label="账号设置" name="settings">
              <div class="tab-content">
                <el-form
                  ref="settingsFormRef"
                  :model="settingsForm"
                  label-width="100px"
                  style="max-width: 600px"
                >
                  <el-form-item label="头像">
                    <div class="avatar-uploader">
                      <el-upload
                        class="avatar-upload"
                        :action="uploadUrl"
                        :headers="uploadHeaders"
                        :show-file-list="false"
                        :on-success="handleAvatarSuccess"
                        :before-upload="beforeAvatarUpload"
                        accept="image/*"
                      >
                        <el-avatar :size="80" :src="settingsForm.avatar" class="avatar-preview">
                          {{ settingsForm.username?.charAt(0)?.toUpperCase() }}
                        </el-avatar>
                        <div class="avatar-upload-text">点击上传头像</div>
                      </el-upload>
                    </div>
                  </el-form-item>

                  <el-form-item label="用户名">
                    <el-input v-model="settingsForm.username" disabled />
                  </el-form-item>

                  <el-form-item label="邮箱">
                    <el-input v-model="settingsForm.email" />
                  </el-form-item>

                  <el-form-item label="昵称">
                    <el-input v-model="settingsForm.nickname" />
                  </el-form-item>

                  <el-form-item label="个人简介">
                    <el-input
                      v-model="settingsForm.bio"
                      type="textarea"
                      :rows="4"
                      maxlength="200"
                      show-word-limit
                    />
                  </el-form-item>

                  <el-form-item label="技术栈">
                    <el-input
                      v-model="settingsForm.skills"
                      type="textarea"
                      :rows="3"
                      placeholder="用逗号分隔，如：Vue.js, React, Go, MySQL"
                    />
                    <div class="form-tip">多个技术用逗号分隔</div>
                  </el-form-item>

                  <el-form-item label="联系方式">
                    <el-input
                      v-model="settingsForm.contacts"
                      type="textarea"
                      :rows="4"
                      placeholder='JSON格式，如：{"email":"xxx@email.com","github":"username","wechat":"微信号"}'
                    />
                    <div class="form-tip">JSON格式，支持 email、github、wechat 等字段</div>
                  </el-form-item>

                  <el-form-item>
                    <el-button type="primary" :loading="saving" @click="handleSaveSettings">
                      保存设置
                    </el-button>
                    <el-button @click="handleResetSettings">重置</el-button>
                  </el-form-item>
                </el-form>

                <el-divider />

                <h3 style="margin-bottom: 20px">修改密码</h3>
                <el-form
                  ref="passwordFormRef"
                  :model="passwordForm"
                  :rules="passwordRules"
                  label-width="100px"
                  style="max-width: 600px"
                >
                  <el-form-item label="原密码" prop="oldPassword">
                    <el-input v-model="passwordForm.oldPassword" type="password" show-password />
                  </el-form-item>

                  <el-form-item label="新密码" prop="newPassword">
                    <el-input v-model="passwordForm.newPassword" type="password" show-password />
                  </el-form-item>

                  <el-form-item label="确认密码" prop="confirmPassword">
                    <el-input v-model="passwordForm.confirmPassword" type="password" show-password />
                  </el-form-item>

                  <el-form-item>
                    <el-button type="primary" :loading="changingPassword" @click="handleChangePassword">
                      修改密码
                    </el-button>
                  </el-form-item>
                </el-form>
              </div>
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useUserStore } from '@/stores/user'
import { getMyLikes, getMyFavorites } from '@/api/article'
import { updateUserInfo, changePassword, getUserStats } from '@/api/auth'
import ArticleCard from '@/components/ArticleCard.vue'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()

const activeTab = ref('likes')
const loading = ref(false)
const saving = ref(false)
const changingPassword = ref(false)

const stats = ref({
  articles: 0,
  likes: 0,
  favorites: 0,
  comments: 0
})

const likedArticles = ref([])
const likesPage = ref(1)
const likesPageSize = ref(10)
const likesTotal = ref(0)

const favoritedArticles = ref([])
const favoritesPage = ref(1)
const favoritesPageSize = ref(10)
const favoritesTotal = ref(0)

const settingsFormRef = ref(null)
const passwordFormRef = ref(null)

const settingsForm = reactive({
  username: '',
  email: '',
  nickname: '',
  bio: '',
  avatar: '',
  skills: '',
  contacts: ''
})

// 文件上传配置
const uploadUrl = '/files/upload'
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${userStore.token}`
}))

const handleAvatarSuccess = (response) => {
  if (response.code === 0 && response.data) {
    settingsForm.avatar = response.data.url
    // 立即更新 userStore，让左侧头像显示出来
    userStore.updateUser({ avatar: response.data.url })
    ElMessage.success('头像上传成功')
  } else {
    ElMessage.error(response.message || '上传失败')
  }
}

const beforeAvatarUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isImage) {
    ElMessage.error('只能上传图片文件！')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过 2MB！')
    return false
  }
  return true
}

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const passwordRules = {
  oldPassword: [
    { required: true, message: '请输入原密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

onMounted(() => {
  initUserSettings()
  fetchLikedArticles()
  fetchStats()
})

// 监听切换tab时加载对应数据
watch(activeTab, (newTab) => {
  if (newTab === 'likes' && likedArticles.value.length === 0) {
    fetchLikedArticles()
  } else if (newTab === 'favorites' && favoritedArticles.value.length === 0) {
    fetchFavoritedArticles()
  }
})

const initUserSettings = () => {
  settingsForm.username = userStore.user?.username || ''
  settingsForm.email = userStore.user?.email || ''
  settingsForm.nickname = userStore.user?.nickname || ''
  settingsForm.bio = userStore.user?.bio || ''
  settingsForm.avatar = userStore.user?.avatar || ''
  settingsForm.skills = userStore.user?.skills || ''
  settingsForm.contacts = userStore.user?.contacts || ''
}

const fetchStats = async () => {
  try {
    const { data } = await getUserStats()
    stats.value = {
      articles: data.articles_count || 0,
      likes: data.likes_count || 0,
      favorites: data.favorites_count || 0,
      comments: data.comments_count || 0
    }
  } catch (error) {
    console.error('Failed to fetch stats:', error)
    // 如果API调用失败，保持默认值
    stats.value = {
      articles: 0,
      likes: 0,
      favorites: 0,
      comments: 0
    }
  }
}

const fetchLikedArticles = async () => {
  loading.value = true
  try {
    const { data } = await getMyLikes({
      page: likesPage.value,
      page_size: likesPageSize.value
    })
    likedArticles.value = data.list || []
    likesTotal.value = data.total || 0
  } catch (error) {
    console.error('Failed to fetch liked articles:', error)
  } finally {
    loading.value = false
  }
}

const fetchFavoritedArticles = async () => {
  loading.value = true
  try {
    const { data } = await getMyFavorites({
      page: favoritesPage.value,
      page_size: favoritesPageSize.value
    })
    favoritedArticles.value = data.list || []
    favoritesTotal.value = data.total || 0
  } catch (error) {
    console.error('Failed to fetch favorited articles:', error)
  } finally {
    loading.value = false
  }
}

const handleSaveSettings = async () => {
  saving.value = true
  try {
    const { data } = await updateUserInfo({
      email: settingsForm.email,
      nickname: settingsForm.nickname,
      bio: settingsForm.bio,
      avatar: settingsForm.avatar,
      skills: settingsForm.skills,
      contacts: settingsForm.contacts
    })
    // 使用后端返回的完整数据更新store
    userStore.updateUser(data)
    ElMessage.success('设置保存成功')
  } catch (error) {
    console.error('Failed to save settings:', error)
  } finally {
    saving.value = false
  }
}

const handleResetSettings = () => {
  initUserSettings()
  ElMessage.info('已重置为原始设置')
}

const handleChangePassword = async () => {
  if (!passwordFormRef.value) return

  await passwordFormRef.value.validate(async (valid) => {
    if (!valid) return

    changingPassword.value = true
    try {
      await changePassword({
        old_password: passwordForm.oldPassword,
        new_password: passwordForm.newPassword
      })
      ElMessage.success('密码修改成功，请重新登录')

      // 清空表单
      passwordForm.oldPassword = ''
      passwordForm.newPassword = ''
      passwordForm.confirmPassword = ''

      // 退出登录
      setTimeout(() => {
        userStore.logout()
        window.location.href = '/login'
      }, 1500)
    } catch (error) {
      console.error('Failed to change password:', error)
    } finally {
      changingPassword.value = false
    }
  })
}
</script>

<style scoped>
.profile {
  padding: 20px 0;
}

.profile-content {
  display: grid;
  grid-template-columns: 300px 1fr;
  gap: 20px;
}

.user-card {
  height: fit-content;
  position: sticky;
  top: 80px;
}

.user-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 20px 0;
}

.user-details {
  text-align: center;
}

.username {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 8px;
}

.user-role {
  margin: 0;
}

.user-stats {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
  padding: 20px 0;
}

.stat-item {
  text-align: center;
}

.stat-number {
  font-size: 28px;
  font-weight: 700;
  color: #409eff;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.content-card {
  min-height: 600px;
}

.tab-content {
  min-height: 400px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}

@media (max-width: 968px) {
  .profile-content {
    grid-template-columns: 1fr;
  }

  .user-card {
    position: static;
  }
}

.avatar-uploader {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.avatar-upload {
  cursor: pointer;
}

.avatar-preview {
  border: 2px dashed #dcdfe6;
  transition: border-color 0.3s;
}

.avatar-preview:hover {
  border-color: #409eff;
}

.avatar-upload-text {
  font-size: 12px;
  color: #909399;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
</style>
