import request from '@/utils/request'

export function uploadFile(file, folder = 'uploads') {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('folder', folder)
  return request.post('/files/upload', formData)
}

export function getFiles(params) {
  return request.get('/files', { params })
}

export function deleteFile(id) {
  return request.delete(`/files/${id}`)
}
