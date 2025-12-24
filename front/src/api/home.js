// 主页相关接口
import request from './request'

/**
 * 获取平台统计数据
 * @returns {Promise} 统计数据
 */
export function getStatistics() {
  return request({
    url: '/home/statistics',
    method: 'get'
  })
}
