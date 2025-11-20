import request from '@/utils/request'

export function getComments(params) {
  return request.get('/comments', { params })
}

export function deleteComment(id) {
  return request.delete(`/comments/${id}`)
}

export function updateCommentStatus(id, status) {
  return request.patch(`/comments/${id}/status`, { status })
}
