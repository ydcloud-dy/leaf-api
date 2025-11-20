import request from '@/utils/request'

export function getArticles(params) {
  return request.get('/articles', { params })
}

export function getArticle(id) {
  return request.get(`/articles/${id}`)
}

export function createArticle(data) {
  return request.post('/articles', data)
}

export function updateArticle(id, data) {
  return request.put(`/articles/${id}`, data)
}

export function updateArticleStatus(id, status) {
  return request.patch(`/articles/${id}/status`, { status })
}

export function deleteArticle(id) {
  return request.delete(`/articles/${id}`)
}

export function getArticleLikes(id, params) {
  return request.get(`/articles/${id}/likes`, { params })
}

export function getArticleFavorites(id, params) {
  return request.get(`/articles/${id}/favorites`, { params })
}

export function getArticleViews(id) {
  return request.get(`/articles/${id}/views`)
}
