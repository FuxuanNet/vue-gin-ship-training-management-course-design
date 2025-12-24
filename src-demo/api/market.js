import request from './request'

// 获取市场资源列表
export function getMarketResources(params) {
  return request({
    url: '/api/v1/market',
    method: 'get',
    params
  })
}

// 获取市场资源详情
export function getMarketResourceDetail(id) {
  return request({
    url: `/api/v1/market/${id}`,
    method: 'get'
  })
}

// 获取资源分类
export function getResourceCategories(type) {
  return request({
    url: '/api/v1/market/categories',
    method: 'get',
    params: { type }
  })
}

// 获取热门标签
export function getPopularTags(type, limit) {
  return request({
    url: '/api/v1/market/tags',
    method: 'get',
    params: { type, limit }
  })
}

// 收藏资源
export function favoriteResource(id) {
  return request({
    url: `/api/v1/market/${id}/favorite`,
    method: 'post'
  })
}

// 取消收藏资源
export function unfavoriteResource(id) {
  return request({
    url: `/api/v1/market/${id}/favorite`,
    method: 'delete'
  })
}

// 获取收藏列表
export function getFavorites() {
  return request({
    url: '/api/v1/market/favorites',
    method: 'get'
  })
}

// 获取资源示例数据
export function getResourceSample(id) {
  return request({
    url: `/api/v1/market/sample/${id}`,
    method: 'get',
    responseType: 'blob', // 使用blob方式接收响应
    timeout: 30000 // 增加超时时间
  })
}

// 购买资源
export function purchaseResource(id) {
  return request({
    url: `/api/v1/market/${id}/purchase`,
    method: 'post'
  })
}

// 获取已购买的资源列表
export function getPurchasedResources() {
  return request({
    url: '/api/v1/market/purchased',
    method: 'get'
  })
}

// 获取用户的资源列表
export function getUserResources() {
  return request({
    url: '/api/v1/user/resources',
    method: 'get'
  })
} 