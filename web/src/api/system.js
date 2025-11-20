import request from '@/utils/request'

export function getStats() {
  return request.get('/stats')
}

export function getHotArticles() {
  return request.get('/stats/hot-articles')
}

export function getSettings() {
  return request.get('/settings')
}

export function updateSettings(data) {
  return request.put('/settings', data)
}
