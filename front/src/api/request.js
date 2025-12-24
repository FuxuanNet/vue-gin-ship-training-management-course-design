// 网络请求基础配置和封装
import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建 axios 实例
const request = axios.create({
  baseURL: 'http://localhost:8080/api', // 后端API基础地址
  timeout: 10000, // 请求超时时间（毫秒）
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器：自动添加 Session-ID
request.interceptors.request.use(
  (config) => {
    // 从 localStorage 获取 Session-ID
    const sessionId = localStorage.getItem('sessionId')
    if (sessionId) {
      config.headers['Session-ID'] = sessionId
    }
    return config
  },
  (error) => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器：统一处理响应和错误
request.interceptors.response.use(
  (response) => {
    const res = response.data

    // 根据后端返回的 code 判断请求是否成功
    if (res.code === 200) {
      return res // 返回成功数据
    } else {
      // 非 200 状态码，显示错误消息
      ElMessage.error(res.message || '请求失败')
      return Promise.reject(new Error(res.message || '请求失败'))
    }
  },
  (error) => {
    console.error('响应错误:', error)
    
    // 处理网络错误
    if (error.response) {
      const status = error.response.status
      switch (status) {
        case 401:
          ElMessage.error('未登录或登录已过期，请重新登录')
          // 清除本地存储的会话信息
          localStorage.removeItem('sessionId')
          localStorage.removeItem('userInfo')
          // 跳转到登录页
          window.location.href = '/login'
          break
        case 403:
          ElMessage.error('没有权限访问')
          break
        case 404:
          ElMessage.error('请求的资源不存在')
          break
        case 500:
          ElMessage.error('服务器内部错误')
          break
        default:
          ElMessage.error(error.response.data?.message || '请求失败')
      }
    } else if (error.request) {
      ElMessage.error('网络连接失败，请检查网络')
    } else {
      ElMessage.error('请求配置错误')
    }
    
    return Promise.reject(error)
  }
)

export default request
