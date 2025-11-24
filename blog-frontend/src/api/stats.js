import request from './request'

// 获取站点统计数据
export const getStats = () => {
  return request.get('/stats')
}

// 获取热门文章（按浏览量排序，最多10篇）
export const getHotArticles = () => {
  return request.get('/stats/hot-articles')
}

// 发送心跳（保持在线状态）
export const sendHeartbeat = () => {
  return request.post('/heartbeat')
}

// 记录访问时长
export const recordVisitDuration = (data) => {
  return request.post('/visit', data)
}

// 获取博主信息
export const getBloggerInfo = () => {
  return request.get('/blogger')
}
