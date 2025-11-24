<template>
  <header class="header">
    <div class="container">
      <div class="header-content">
        <div class="logo" @click="router.push('/')">
          <el-icon :size="28"><Reading /></el-icon>
          <span>个人博客</span>
        </div>

        <nav class="nav">
          <router-link to="/" class="nav-link" exact>首页</router-link>
          <router-link to="/articles" class="nav-link">文章</router-link>

          <!-- 笔记下拉菜单 -->
          <el-dropdown @command="handleNoteCommand" class="notes-dropdown">
            <span class="nav-link notes-link">
              笔记
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item
                  v-for="tag in tags"
                  :key="tag.id"
                  :command="tag.name"
                >
                  {{ tag.name }}
                </el-dropdown-item>
                <el-dropdown-item v-if="tags.length === 0" disabled>
                  暂无分类
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>

          <router-link to="/archive" class="nav-link">归档</router-link>
          <router-link to="/guestbook" class="nav-link">留言板</router-link>
          <router-link to="/stats" class="nav-link">统计</router-link>
          <router-link to="/about" class="nav-link">关于</router-link>
        </nav>

        <div class="header-actions">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索文章..."
            class="search-input"
            clearable
            @input="handleSearchInput"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>

          <div v-if="userStore.isLoggedIn" class="user-info">
            <el-dropdown @command="handleCommand">
              <div class="user-dropdown">
                <el-avatar :size="36" :src="userStore.avatar">
                  {{ userStore.username.charAt(0).toUpperCase() }}
                </el-avatar>
                <span class="username">{{ userStore.username }}</span>
                <el-icon><ArrowDown /></el-icon>
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">
                    <el-icon><User /></el-icon>
                    个人中心
                  </el-dropdown-item>
                  <el-dropdown-item divided command="logout">
                    <el-icon><SwitchButton /></el-icon>
                    退出登录
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>

          <el-button v-else type="primary" @click="router.push('/login')">
            登录
          </el-button>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getTags } from '@/api/tag'
import { ElMessage } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()
const searchKeyword = ref('')
const tags = ref([])
let searchTimer = null

onMounted(() => {
  fetchTags()
})

const fetchTags = async () => {
  try {
    const { data } = await getTags()
    const tagList = Array.isArray(data) ? data : (data.list || [])
    tags.value = tagList
  } catch (error) {
    console.error('Failed to fetch tags:', error)
  }
}

const handleSearchInput = () => {
  // 清除之前的定时器
  if (searchTimer) {
    clearTimeout(searchTimer)
  }

  // 设置新的定时器,500ms后执行搜索
  searchTimer = setTimeout(() => {
    handleSearch()
  }, 500)
}

const handleSearch = () => {
  const keyword = searchKeyword.value.trim()

  // 跳转到文章列表页,如果有关键词则传递,否则显示全部文章
  router.push({
    name: 'Articles',
    query: keyword ? { keyword } : {}
  })
}

const handleNoteCommand = (tagName) => {
  // 跳转到笔记页面,传递标签名称
  router.push({
    name: 'Notes',
    params: { tag: tagName }
  })
}

const handleCommand = (command) => {
  if (command === 'profile') {
    router.push('/profile')
  } else if (command === 'logout') {
    userStore.logout()
    ElMessage.success('已退出登录')
    router.push('/')
  }
}
</script>

<style scoped>
.header {
  background-color: #fff;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 1000;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 0;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 20px;
  font-weight: 600;
  color: #409eff;
  cursor: pointer;
  transition: opacity 0.3s;
}

.logo:hover {
  opacity: 0.8;
}

.nav {
  display: flex;
  gap: 30px;
  flex: 1;
  justify-content: center;
}

.nav-link {
  color: #606266;
  text-decoration: none;
  font-size: 16px;
  transition: color 0.3s;
  position: relative;
}

.nav-link:hover {
  color: #409eff;
}

.nav-link.router-link-active {
  color: #409eff;
  font-weight: 500;
}

.nav-link.router-link-active::after {
  content: '';
  position: absolute;
  bottom: -16px;
  left: 0;
  right: 0;
  height: 2px;
  background-color: #409eff;
}

.notes-dropdown {
  display: inline-flex;
  align-items: center;
}

.notes-link {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  color: #606266;
  font-size: 16px;
  transition: color 0.3s;
  text-decoration: none;
  position: relative;
  line-height: normal;
}

.notes-link:hover {
  color: #409eff;
}

.notes-link .el-icon {
  font-size: 12px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 20px;
}

.search-input {
  width: 200px;
}

.user-info {
  cursor: pointer;
}

.user-dropdown {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 8px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.user-dropdown:hover {
  background-color: #f5f7fa;
}

.username {
  font-size: 14px;
  color: #606266;
}

@media (max-width: 768px) {
  .nav {
    display: none;
  }

  .search-input {
    width: 150px;
  }
}
</style>
