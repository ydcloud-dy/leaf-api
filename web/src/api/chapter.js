import request from '@/utils/request'

export function getChapters(params) {
  return request.get('/chapters', { params })
}

export function getChapter(id) {
  return request.get(`/chapters/${id}`)
}

export function createChapter(data) {
  return request.post('/chapters', data)
}

export function updateChapter(id, data) {
  return request.put(`/chapters/${id}`, data)
}

export function deleteChapter(id) {
  return request.delete(`/chapters/${id}`)
}
