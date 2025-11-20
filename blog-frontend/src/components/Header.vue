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
          <router-link to="/archive" class="nav-link">归档</router-link>
          <router-link to="/guestbook" class="nav-link">留言板</router-link>
          <router-link to="/about" class="nav-link">关于</router-link>
        </nav>

        <div class="header-actions">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索文章..."
            class="search-input"
            clearable
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
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const searchKeyword = ref('')

const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    router.push({
      name: 'Articles',
      query: { keyword: searchKeyword.value }
    })
  }
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
