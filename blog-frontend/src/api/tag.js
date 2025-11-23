import request from './request'

// 获取标签列表
export function getTags(params) {
  return request({
    url: '/tags',
    method: 'get',
    params
  })
}
