import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'

// 创建 axios 实例
const request = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
request.interceptors.request.use(
  config => {
    const userStore = useUserStore()

    // 添加 token
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }

    return config
  },
  error => {
    console.error('Request error:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  response => {
    const res = response.data

    // 后端统一返回格式：{ code: 0, message: "success", data: {...} }
    // code 为 0 表示成功，其他值表示失败
    if (res.code === 0) {
      return res
    }

    // 处理业务错误
    const errorMsg = res.message || '请求失败'

    // 根据 code 进行不同处理
    if (res.code === 401) {
      ElMessage.error('未授权，请登录')
      const userStore = useUserStore()
      userStore.logout()
      window.location.href = '/login'
    } else if (res.code === 403) {
      ElMessage.error('拒绝访问')
    } else if (res.code === 404) {
      ElMessage.error('请求的资源不存在')
    } else if (res.code === 500) {
      ElMessage.error('服务器错误')
    } else {
      ElMessage.error(errorMsg)
    }

    return Promise.reject(new Error(errorMsg))
  },
  error => {
    console.error('Response error:', error)

    if (error.response) {
      const { status } = error.response

      if (status === 401) {
        ElMessage.error('未授权，请登录')
        const userStore = useUserStore()
        userStore.logout()
        window.location.href = '/login'
      } else {
        ElMessage.error('请求失败')
      }
    } else if (error.request) {
      ElMessage.error('网络错误，请检查网络连接')
    } else {
      ElMessage.error('请求配置错误')
    }

    return Promise.reject(error)
  }
)

export default request
