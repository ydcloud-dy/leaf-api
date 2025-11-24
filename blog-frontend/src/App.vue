<template>
  <router-view />
</template>

<script setup>
import { onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { useHeartbeat } from '@/composables/useHeartbeat'
import { useVisitTracking } from '@/composables/useVisitTracking'

const userStore = useUserStore()

// 启动心跳服务（保持在线状态）
useHeartbeat()

// 启动访问追踪服务（内部已自动监听路由变化）
useVisitTracking()

onMounted(() => {
  // 尝试从本地存储恢复用户信息
  userStore.initUser()
})
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  background-color: #f5f7fa;
}

#app {
  min-height: 100vh;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.page-header {
  padding: 40px 0 20px;
  text-align: center;
}

.page-title {
  font-size: 32px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 10px;
}

.page-subtitle {
  font-size: 16px;
  color: #909399;
}
</style>
