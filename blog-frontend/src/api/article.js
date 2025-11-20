import request from './request'

// 获取文章列表
export function getArticles(params) {
  return request({
    url: '/articles',
    method: 'get',
    params
  })
}

// 获取文章详情
export function getArticleDetail(id) {
  return request({
    url: `/articles/${id}`,
    method: 'get'
  })
}

// 搜索文章
export function searchArticles(params) {
  return request({
    url: '/articles/search',
    method: 'get',
    params
  })
}

// 获取归档文章
export function getArchiveArticles(params) {
  return request({
    url: '/articles/archive',
    method: 'get',
    params
  })
}

// 点赞文章
export function likeArticle(id) {
  return request({
    url: `/articles/${id}/like`,
    method: 'post'
  })
}

// 取消点赞
export function unlikeArticle(id) {
  return request({
    url: `/articles/${id}/like`,
    method: 'delete'
  })
}

// 收藏文章
export function favoriteArticle(id) {
  return request({
    url: `/articles/${id}/favorite`,
    method: 'post'
  })
}

// 取消收藏
export function unfavoriteArticle(id) {
  return request({
    url: `/articles/${id}/favorite`,
    method: 'delete'
  })
}

// 获取我的点赞列表
export function getMyLikes(params) {
  return request({
    url: '/user/likes',
    method: 'get',
    params
  })
}

// 获取我的收藏列表
export function getMyFavorites(params) {
  return request({
    url: '/user/favorites',
    method: 'get',
    params
  })
}
