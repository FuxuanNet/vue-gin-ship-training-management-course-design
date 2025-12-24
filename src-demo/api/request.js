import axios from 'axios'
import router from '../router'

// API基础URL配置
export const API_BASE_URL = ''

// 创建 axios 实例
const service = axios.create({
  baseURL: API_BASE_URL, // API 基础URL
  timeout: 15000 // 请求超时时间
})

// 请求拦截器
service.interceptors.request.use(
  config => {
    // 从 localStorage 获取 token
    const token = localStorage.getItem('token')
    
    // 如果 token 存在，则添加到请求头
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    
    return config
  },
  error => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  response => {
    // 如果响应类型是blob，直接返回不做处理
    if (response.config.responseType === 'blob') {
      // 检查是否是JSON格式的错误响应
      const contentType = response.headers['content-type']
      if (contentType && contentType.includes('application/json')) {
        // 可能是JSON格式的错误响应，尝试读取
        return new Promise((resolve, reject) => {
          const reader = new FileReader()
          reader.onload = () => {
            try {
              const result = JSON.parse(reader.result)
              if (result.code !== 200) {
                reject(result)
              } else {
                resolve(response) // 正常的JSON响应
              }
            } catch (e) {
              // 如果无法解析为JSON，则视为正常的blob响应
              resolve(response)
            }
          }
          reader.onerror = () => resolve(response) // 读取错误时仍然返回原始响应
          reader.readAsText(response.data)
        })
      }
      return response // 正常的blob响应
    }

    const res = response.data
    
    // 根据状态码判断请求是否成功
    if (res.code !== 200) {
      // 处理特定错误
      if (res.code === 401) {
        // 401 未授权，清除 token 并重定向到登录页
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        router.push('/login')
      }
      
      // 直接返回错误消息，而不是创建新的Error对象
      return Promise.reject(res)
    } else {
      return res
    }
  },
  error => {
    console.error('响应错误:', error)
    
    // 处理 401 错误
    if (error.response && error.response.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      router.push('/login')
    }
    
    // 尝试提取更有用的错误信息
    let errorMsg = '网络错误，请检查您的网络连接'
    
    if (error.response) {
      if (error.response.data) {
        if (typeof error.response.data === 'object' && error.response.data.message) {
          errorMsg = error.response.data.message
        } else if (typeof error.response.data === 'string') {
          errorMsg = error.response.data
        }
      } else {
        errorMsg = `服务器错误 (${error.response.status})`
      }
    } else if (error.request) {
      errorMsg = '服务器未响应，请检查网络连接'
    } else if (error.message) {
      errorMsg = error.message
    }
    
    return Promise.reject({ message: errorMsg })
  }
)

export default service 