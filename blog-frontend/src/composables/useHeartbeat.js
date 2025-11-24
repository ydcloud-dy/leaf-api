import { onMounted, onUnmounted } from 'vue'
import { sendHeartbeat as sendHeartbeatAPI } from '@/api/stats'

/**
 * 心跳服务 - 保持在线状态
 * 每 30 秒发送一次心跳，Redis 会记录在线状态
 */
export function useHeartbeat() {
  let heartbeatTimer = null

  // 发送心跳请求
  const sendHeartbeat = async () => {
    try {
      await sendHeartbeatAPI()
    } catch (error) {
      console.warn('Heartbeat failed:', error)
    }
  }

  // 启动心跳
  const startHeartbeat = () => {
    // 立即发送一次
    sendHeartbeat()

    // 每 30 秒发送一次
    heartbeatTimer = setInterval(() => {
      sendHeartbeat()
    }, 30000) // 30 秒
  }

  // 停止心跳
  const stopHeartbeat = () => {
    if (heartbeatTimer) {
      clearInterval(heartbeatTimer)
      heartbeatTimer = null
    }
  }

  // 自动管理生命周期
  onMounted(() => {
    startHeartbeat()
  })

  onUnmounted(() => {
    stopHeartbeat()
  })

  return {
    startHeartbeat,
    stopHeartbeat
  }
}
