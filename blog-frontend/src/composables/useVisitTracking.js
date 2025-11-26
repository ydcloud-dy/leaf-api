import { onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { recordVisitDuration } from '@/api/stats'

/**
 * 访问追踪服务 - 记录页面停留时长
 * 用于统计平均访问时长和24小时PV
 */
export function useVisitTracking() {
  const router = useRouter()
  let startTime = 0
  let currentPath = ''
  let isReady = false

  // 记录页面访问（进入时立即记录，duration=0）
  const recordPageEnter = async (path) => {
    if (!path) return

    try {
      await recordVisitDuration({
        path: path,
        duration: 0
      })
    } catch (error) {
      console.warn('Failed to record page enter:', error)
    }
  }

  // 记录页面访问时长（离开时记录实际时长）
  const recordVisit = () => {
    if (!startTime || !currentPath) return

    const duration = Math.floor((Date.now() - startTime) / 1000)

    // 只记录停留时间超过 3 秒的访问
    if (duration < 3) return

    const data = {
      path: currentPath,
      duration: duration
    }

    // 使用 sendBeacon（即使页面关闭也能发送）
    if (navigator.sendBeacon) {
      try {
        const blob = new Blob([JSON.stringify(data)], { type: 'application/json' })
        // 使用相对路径，由 Nginx 代理到后端
        navigator.sendBeacon('/blog/visit', blob)
      } catch (error) {
        console.warn('sendBeacon failed:', error)
      }
    }
  }

  // 开始追踪新页面
  const startTracking = (path) => {
    currentPath = path
    startTime = Date.now()
    recordPageEnter(path)
  }

  // 监听页面可见性变化
  const handleVisibilityChange = () => {
    if (document.hidden) {
      recordVisit()
    } else {
      startTime = Date.now()
    }
  }

  // 页面卸载时记录
  const handleBeforeUnload = () => {
    recordVisit()
  }

  onMounted(() => {
    // 等待路由完全准备好后再开始追踪
    router.isReady().then(() => {
      if (isReady) return
      isReady = true

      // 记录初始页面
      startTracking(router.currentRoute.value.path)
    })

    // 监听路由变化（afterEach 确保导航完成后才触发）
    router.afterEach((to, from) => {
      if (!isReady) return

      // 记录上一页停留时长
      if (from.path && from.path !== to.path) {
        recordVisit()
      }

      // 开始追踪新页面
      startTracking(to.path)
    })

    document.addEventListener('visibilitychange', handleVisibilityChange)
    window.addEventListener('beforeunload', handleBeforeUnload)
  })

  onUnmounted(() => {
    recordVisit()
    document.removeEventListener('visibilitychange', handleVisibilityChange)
    window.removeEventListener('beforeunload', handleBeforeUnload)
  })

  return {
    startTracking,
    recordVisit
  }
}
