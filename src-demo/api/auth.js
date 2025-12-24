import request, { API_BASE_URL } from './request'

// 用户注册
export function register(data) {
  return request({
    url: '/api/v1/auth/register',
    method: 'post',
    data
  })
}

// 用户登录
export function login(data) {
  return request({
    url: '/api/v1/auth/login',
    method: 'post',
    data
  })
}

// 用户退出
export function logout() {
  return request({
    url: '/api/v1/auth/logout',
    method: 'post'
  })
}

// 获取当前用户信息
export function getCurrentUser() {
  return request({
    url: '/api/v1/auth/current',
    method: 'get'
  })
}

// 修改密码
export function updatePassword(data) {
  return request({
    url: '/api/v1/auth/password',
    method: 'put',
    data
  })
}

// 上传头像
export function uploadAvatar(file) {
  const formData = new FormData()
  formData.append('file', file)
  
  return request({
    url: '/api/v1/auth/avatar',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 获取头像URL
export function getAvatarUrl(userId) {
  return `${API_BASE_URL}/api/v1/auth/avatar/${userId}`
}

// 充值余额
export function depositBalance(data) {
  return request({
    url: '/api/v1/auth/deposit',
    method: 'post',
    data
  })
}

// 提现余额
export function withdrawBalance(data) {
  return request({
    url: '/api/v1/auth/withdraw',
    method: 'post',
    data
  })
}

// 获取交易记录
export function getTransactions(params) {
  return request({
    url: '/api/v1/auth/transactions',
    method: 'get',
    params
  })
} 