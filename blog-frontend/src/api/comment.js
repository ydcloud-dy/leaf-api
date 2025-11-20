import request from './request'

// 获取评论列表
export function getComments(params) {
  return request({
    url: '/comments',
    method: 'get',
    params
  })
}

// 获取文章评论
export function getArticleComments(articleId, params) {
  return request({
    url: `/articles/${articleId}/comments`,
    method: 'get',
    params
  })
}

// 发表评论
export function createComment(data) {
  return request({
    url: '/comments',
    method: 'post',
    data
  })
}

// 回复评论
export function replyComment(id, data) {
  return request({
    url: `/comments/${id}/reply`,
    method: 'post',
    data
  })
}

// 删除评论
export function deleteComment(id) {
  return request({
    url: `/comments/${id}`,
    method: 'delete'
  })
}

// 点赞评论
export function likeComment(id) {
  return request({
    url: `/comments/${id}/like`,
    method: 'post'
  })
}

// 取消点赞评论
export function unlikeComment(id) {
  return request({
    url: `/comments/${id}/like`,
    method: 'delete'
  })
}

// 获取留言板消息列表
export function getGuestbookMessages(params) {
  return request({
    url: '/guestbook',
    method: 'get',
    params
  })
}

// 创建留言板消息
export function createGuestbookMessage(data) {
  return request({
    url: '/guestbook',
    method: 'post',
    data
  })
}
