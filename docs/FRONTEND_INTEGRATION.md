# 前端集成指南 - 在线统计功能

本文档提供前端集成在线人数统计和访问时长记录功能的示例代码。

## 📦 安装依赖

项目已有 axios，无需额外安装。

## 🔌 API 接口说明

### 1. 心跳接口（保持在线状态）
- **接口**: `POST /blog/heartbeat`
- **说明**: 前端每 30 秒调用一次，保持用户在线状态
- **认证**: 不需要（支持未登录用户）

### 2. 记录访问时长
- **接口**: `POST /blog/visit`
- **说明**: 页面关闭或切换时上报停留时长
- **认证**: 不需要（支持未登录用户）
- **参数**:
  ```json
  {
    "path": "/blog/articles/123",
    "duration": 120  // 秒
  }
  ```

### 3. 获取统计数据
- **接口**: `GET /blog/stats`
- **说明**: 获取网站统计数据（包含在线人数和平均访问时长）
- **返回数据**:
  ```json
  {
    "code": 200,
    "data": {
      "article_count": 92,
      "chapter_count": 1124,
      "category_count": 11,
      "tag_count": 9,
      "user_count": 5,
      "comment_count": 12,
      "total_views": 123456,
      "today_views": 2913,
      "online_count": 9,           // 当前在线人数
      "avg_visit_duration": 77,    // 平均访问时长（秒）
      "site_runtime": 1849         // 网站运行天数
    }
  }
  ```

## 📝 Vue 3 集成示例

### 1. 创建统计 API 服务

在 `src/api/stats.js` 中创建：

```javascript
import request from '@/utils/request'

// 发送心跳
export function sendHeartbeat() {
  return request({
    url: '/blog/heartbeat',
    method: 'post'
  })
}

// 记录访问时长
export function recordVisitDuration(data) {
  return request({
    url: '/blog/visit',
    method: 'post',
    data
  })
}

// 获取统计数据
export function getStats() {
  return request({
    url: '/blog/stats',
    method: 'get'
  })
}
```

### 2. 创建心跳服务 (Composable)

在 `src/composables/useHeartbeat.js` 中创建：

```javascript
import { onMounted, onUnmounted } from 'vue'
import { sendHeartbeat } from '@/api/stats'

export function useHeartbeat() {
  let heartbeatTimer = null

  // 启动心跳
  const startHeartbeat = () => {
    // 立即发送一次心跳
    sendHeartbeat().catch(() => {
      console.warn('Heart beat failed')
    })

    // 每 30 秒发送一次心跳
    heartbeatTimer = setInterval(() => {
      sendHeartbeat().catch(() => {
        console.warn('Heart beat failed')
      })
    }, 30000) // 30 秒
  }

  // 停止心跳
  const stopHeartbeat = () => {
    if (heartbeatTimer) {
      clearInterval(heartbeatTimer)
      heartbeatTimer = null
    }
  }

  // 组件挂载时启动
  onMounted(() => {
    startHeartbeat()
  })

  // 组件卸载时停止
  onUnmounted(() => {
    stopHeartbeat()
  })

  return {
    startHeartbeat,
    stopHeartbeat
  }
}
```

### 3. 创建访问时长追踪服务 (Composable)

在 `src/composables/useVisitTracking.js` 中创建：

```javascript
import { onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { recordVisitDuration } from '@/api/stats'

export function useVisitTracking() {
  const route = useRoute()
  let startTime = 0
  let currentPath = ''

  // 记录页面访问时长
  const recordVisit = () => {
    if (!startTime || !currentPath) return

    const duration = Math.floor((Date.now() - startTime) / 1000) // 转换为秒

    // 只记录停留时间超过 3 秒的访问
    if (duration < 3) return

    // 使用 sendBeacon 确保数据发送（即使页面关闭）
    const data = JSON.stringify({
      path: currentPath,
      duration: duration
    })

    if (navigator.sendBeacon) {
      const blob = new Blob([data], { type: 'application/json' })
      navigator.sendBeacon('/api/blog/visit', blob)
    } else {
      // 降级方案：使用普通 AJAX
      recordVisitDuration({
        path: currentPath,
        duration: duration
      }).catch(() => {
        console.warn('Failed to record visit duration')
      })
    }
  }

  // 开始追踪
  const startTracking = (path) => {
    // 先记录上一个页面
    if (startTime && currentPath) {
      recordVisit()
    }

    // 开始新页面追踪
    currentPath = path || route.path
    startTime = Date.now()
  }

  // 停止追踪
  const stopTracking = () => {
    recordVisit()
    startTime = 0
    currentPath = ''
  }

  // 监听页面可见性变化
  const handleVisibilityChange = () => {
    if (document.hidden) {
      // 页面隐藏时记录
      recordVisit()
    } else {
      // 页面重新可见时重置计时
      startTime = Date.now()
    }
  }

  // 页面卸载时记录
  const handleBeforeUnload = () => {
    recordVisit()
  }

  onMounted(() => {
    startTracking()

    // 监听页面可见性变化
    document.addEventListener('visibilitychange', handleVisibilityChange)

    // 监听页面卸载
    window.addEventListener('beforeunload', handleBeforeUnload)
  })

  onUnmounted(() => {
    stopTracking()
    document.removeEventListener('visibilitychange', handleVisibilityChange)
    window.removeEventListener('beforeunload', handleBeforeUnload)
  })

  return {
    startTracking,
    stopTracking,
    recordVisit
  }
}
```

### 4. 在根组件中使用

在 `App.vue` 或布局组件中：

```vue
<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script setup>
import { watch } from 'vue'
import { useRoute } from 'vue-router'
import { useHeartbeat } from '@/composables/useHeartbeat'
import { useVisitTracking } from '@/composables/useVisitTracking'

const route = useRoute()

// 启动心跳
useHeartbeat()

// 启动访问追踪
const { startTracking, recordVisit } = useVisitTracking()

// 监听路由变化，记录上一页访问时长并开始新页面追踪
watch(() => route.path, (newPath) => {
  recordVisit() // 记录上一页
  startTracking(newPath) // 开始追踪新页面
})
</script>
```

### 5. 在统计页面展示数据

```vue
<template>
  <div class="stats-container">
    <h1>📊 网站统计</h1>

    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">🔄</div>
        <div class="stat-label">网站运行时长</div>
        <div class="stat-value">{{ stats.site_runtime }}天</div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">👥</div>
        <div class="stat-label">24小时访问量</div>
        <div class="stat-value">{{ stats.today_views }}次</div>
      </div>

      <div class="stat-card highlight">
        <div class="stat-icon">🌐</div>
        <div class="stat-label">当前在线人数</div>
        <div class="stat-value">{{ stats.online_count }}人</div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">⏱️</div>
        <div class="stat-label">平均访问时长</div>
        <div class="stat-value">{{ formatDuration(stats.avg_visit_duration) }}</div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">📝</div>
        <div class="stat-label">文章篇数</div>
        <div class="stat-value">{{ stats.article_count }}篇</div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">📔</div>
        <div class="stat-label">笔记篇数</div>
        <div class="stat-value">{{ stats.chapter_count }}篇</div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">📚</div>
        <div class="stat-label">文章分类数</div>
        <div class="stat-value">{{ stats.category_count }}个</div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">🏷️</div>
        <div class="stat-label">文章标签数</div>
        <div class="stat-value">{{ stats.tag_count }}个</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getStats } from '@/api/stats'

const stats = ref({
  article_count: 0,
  chapter_count: 0,
  category_count: 0,
  tag_count: 0,
  user_count: 0,
  comment_count: 0,
  total_views: 0,
  today_views: 0,
  online_count: 0,
  avg_visit_duration: 0,
  site_runtime: 0
})

// 格式化时长（秒转为分钟）
const formatDuration = (seconds) => {
  if (seconds < 60) {
    return `${Math.round(seconds)}秒/页`
  }
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = Math.round(seconds % 60)
  return `${minutes}分${remainingSeconds}秒/页`
}

// 加载统计数据
const loadStats = async () => {
  try {
    const res = await getStats()
    stats.value = res.data
  } catch (error) {
    console.error('Failed to load stats:', error)
  }
}

onMounted(() => {
  loadStats()

  // 每 30 秒刷新一次统计数据（在线人数会实时变化）
  setInterval(loadStats, 30000)
})
</script>

<style scoped>
.stats-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-top: 30px;
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  text-align: center;
  transition: transform 0.3s, box-shadow 0.3s;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.stat-card.highlight {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.stat-icon {
  font-size: 36px;
  margin-bottom: 12px;
}

.stat-label {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.stat-card.highlight .stat-label {
  color: rgba(255, 255, 255, 0.9);
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #333;
}

.stat-card.highlight .stat-value {
  color: white;
}
</style>
```

## 🔧 可选优化

### 1. 使用 localStorage 缓存心跳状态

避免多个标签页重复发送心跳：

```javascript
// 在 useHeartbeat.js 中添加
const TAB_ID = `tab_${Date.now()}_${Math.random()}`
const HEARTBEAT_KEY = 'app_heartbeat'

const startHeartbeat = () => {
  const sendHeartbeatIfNeeded = () => {
    const lastHeartbeat = localStorage.getItem(HEARTBEAT_KEY)
    const now = Date.now()

    // 如果最后一次心跳距现在超过 25 秒，则发送
    if (!lastHeartbeat || now - parseInt(lastHeartbeat) > 25000) {
      sendHeartbeat().then(() => {
        localStorage.setItem(HEARTBEAT_KEY, now.toString())
      }).catch(() => {
        console.warn('Heart beat failed')
      })
    }
  }

  sendHeartbeatIfNeeded()
  heartbeatTimer = setInterval(sendHeartbeatIfNeeded, 30000)
}
```

### 2. 添加统计数据缓存

减少 API 请求频率：

```javascript
// 使用 Pinia store 缓存统计数据
import { defineStore } from 'pinia'
import { getStats } from '@/api/stats'

export const useStatsStore = defineStore('stats', {
  state: () => ({
    stats: null,
    lastUpdate: 0,
    cacheTime: 30000 // 30秒缓存
  }),

  actions: {
    async fetchStats(force = false) {
      const now = Date.now()

      // 如果缓存有效且不是强制刷新，直接返回
      if (!force && this.stats && (now - this.lastUpdate) < this.cacheTime) {
        return this.stats
      }

      try {
        const res = await getStats()
        this.stats = res.data
        this.lastUpdate = now
        return this.stats
      } catch (error) {
        console.error('Failed to fetch stats:', error)
        return this.stats
      }
    }
  }
})
```

## 📊 后端数据初始化

在数据库的 `settings` 表中添加网站启动时间：

```sql
INSERT INTO settings (`key`, `value`, updated_at)
VALUES ('site_start_date', '2020-01-01', NOW())
ON DUPLICATE KEY UPDATE value = value;
```

## ✅ 功能测试清单

- [ ] 打开网站后，在线人数 +1
- [ ] 关闭网站后 60 秒内，在线人数 -1
- [ ] 浏览文章后关闭，记录访问时长
- [ ] 统计接口返回正确的在线人数
- [ ] 统计接口返回平均访问时长（秒）
- [ ] 多个标签页不会重复发送心跳（如果启用 localStorage 优化）

## 🐛 常见问题

### 1. Redis 连接失败

**错误**: `Failed to initialize Redis`

**解决方案**:
- 确保 Redis 服务已启动
- 检查 `config.yaml` 中的 Redis 配置
- 验证 Redis 端口和密码是否正确

### 2. 在线人数始终为 0

**可能原因**:
- Redis 未启动或连接失败
- 前端心跳请求未发送
- 心跳间隔设置不合理

**解决方案**:
- 查看后端日志，确认 Redis 连接成功
- 在浏览器 Network 面板检查心跳请求
- 确认路由中间件配置正确

### 3. 访问时长记录失败

**可能原因**:
- sendBeacon API 不支持（旧浏览器）
- 跨域问题
- 请求参数格式错误

**解决方案**:
- 使用降级方案（普通 AJAX）
- 配置 CORS 允许 sendBeacon
- 检查请求 payload 格式

## 📚 参考资料

- [Navigator.sendBeacon API](https://developer.mozilla.org/en-US/docs/Web/API/Navigator/sendBeacon)
- [Redis 官方文档](https://redis.io/documentation)
- [Vue 3 Composables](https://vuejs.org/guide/reusability/composables.html)
