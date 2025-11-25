import request from '@/utils/request'

export function login(data) {
  return request.post('/auth/login', data)
}

export function logout() {
  return request.post('/auth/logout')
}

export function getProfile() {
  return request.get('/auth/profile')
}

export function updateProfile(data) {
  return request.put('/auth/profile', data)
}
