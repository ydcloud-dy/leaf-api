import request from './request'

// 获取站点统计数据
export const getStats = () => {
  return request.get('/stats')
}

// 获取热门文章（按浏览量排序，最多10篇）
export const getHotArticles = () => {
  return request.get('/stats/hot-articles')
}
