// 讲师相关接口
import request from './request'

/**
 * 获取讲师今日授课列表
 * @returns {Promise} 今日授课数据
 */
export function getTodayCourses() {
  return request({
    url: '/teacher/today-courses',
    method: 'get'
  })
}

/**
 * 获取讲师授课表
 * @param {Object} params - 查询参数
 * @param {string} params.startDate - 开始日期 (YYYY-MM-DD)
 * @param {string} params.endDate - 结束日期 (YYYY-MM-DD)
 * @returns {Promise} 授课表数据
 */
export function getSchedule(params) {
  return request({
    url: '/teacher/schedule',
    method: 'get',
    params
  })
}

/**
 * 获取待评分学员列表
 * @param {Object} params - 查询参数
 * @param {number} params.courseId - 课程ID（可选）
 * @param {string} params.status - 筛选状态：pending/all（可选）
 * @returns {Promise} 待评分学员数据
 */
export function getPendingEvaluations(params) {
  return request({
    url: '/teacher/pending-evaluations',
    method: 'get',
    params
  })
}

/**
 * 提交学员评分
 * @param {Object} data - 评分数据
 * @param {number} data.itemId - 课程安排ID
 * @param {number} data.personId - 学员ID
 * @param {number} data.teacherScore - 讲师评分（0-100）
 * @param {string} data.teacherComment - 讲师评语
 * @param {number} data.scoreRatio - 讲师评分占比（0-1）
 * @returns {Promise} 提交结果
 */
export function submitGrading(data) {
  return request({
    url: '/teacher/submit-grading',
    method: 'post',
    data
  })
}

/**
 * 获取课程成绩统计
 * @param {Object} params - 查询参数
 * @param {number} params.courseId - 课程ID
 * @returns {Promise} 成绩统计数据
 */
export function getCourseStatistics(params) {
  return request({
    url: '/teacher/course-statistics',
    method: 'get',
    params
  })
}

/**
 * 获取讲师授课统计
 * @param {Object} params - 查询参数
 * @param {string} params.startDate - 开始日期（可选）
 * @param {string} params.endDate - 结束日期（可选）
 * @returns {Promise} 授课统计数据
 */
export function getTeachingStatistics(params) {
  return request({
    url: '/teacher/teaching-statistics',
    method: 'get',
    params
  })
}
