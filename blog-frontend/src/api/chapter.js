import request from './request'

// 获取标签下的章节及文章
export function getChaptersByTag(tag) {
  return request({
    url: `/chapters/${tag}`,
    method: 'get'
  })
}

// 获取章节列表
export function getChapters(params) {
  return request({
    url: '/chapters',
    method: 'get',
    params
  })
}

// 创建章节
export function createChapter(data) {
  return request({
    url: '/chapters',
    method: 'post',
    data
  })
}

// 更新章节
export function updateChapter(id, data) {
  return request({
    url: `/chapters/${id}`,
    method: 'put',
    data
  })
}

// 删除章节
export function deleteChapter(id) {
  return request({
    url: `/chapters/${id}`,
    method: 'delete'
  })
}
