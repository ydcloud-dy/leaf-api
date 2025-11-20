import request from '@/utils/request'

export function getTags() {
  return request.get('/tags')
}

export function createTag(data) {
  return request.post('/tags', data)
}

export function deleteTag(id) {
  return request.delete(`/tags/${id}`)
}

export function getCategories() {
  return request.get('/categories')
}

export function createCategory(data) {
  return request.post('/categories', data)
}

export function deleteCategory(id) {
  return request.delete(`/categories/${id}`)
}
