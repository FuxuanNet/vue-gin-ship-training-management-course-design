// 员工相关接口
import request from './request'

/**
 * 获取员工课程表
 * @param {Object} params - 查询参数
 * @param {string} params.startDate - 开始日期 (YYYY-MM-DD)
 * @param {string} params.endDate - 结束日期 (YYYY-MM-DD)
 * @returns {Promise} 课程表数据
 */
export function getSchedule(params) {
  return request({
    url: '/employee/schedule',
    method: 'get',
    params
  })
}

/**
 * 获取待自评课程列表
 * @param {Object} params - 查询参数
 * @param {string} params.status - 筛选状态：pending（待自评）/all（全部）
 * @param {number} params.limit - 返回数量限制
 * @returns {Promise} 待自评课程数据
 */
export function getPendingEvaluations(params) {
  return request({
    url: '/employee/pending-evaluations',
    method: 'get',
    params
  })
}

/**
 * 提交课程自评
 * @param {Object} data - 自评数据
 * @param {number} data.itemId - 课程安排ID
 * @param {string} data.selfComment - 自评内容（50-1000字）
 * @param {number} data.understanding - 理解程度（1-5，可选）
 * @param {number} data.difficulty - 难度感受（1-5，可选）
 * @param {number} data.satisfaction - 满意度（1-5，可选）
 * @returns {Promise} 提交结果
 */
export function submitEvaluation(data) {
  return request({
    url: '/employee/submit-evaluation',
    method: 'post',
    data
  })
}

/**
 * 获取员工成绩列表
 * @param {Object} params - 查询参数
 * @param {number} params.planId - 培训计划ID筛选（可选）
 * @param {string} params.courseClass - 课程类型筛选（可选）
 * @param {string} params.startDate - 开始日期（可选）
 * @param {string} params.endDate - 结束日期（可选）
 * @returns {Promise} 成绩数据
 */
export function getScores(params) {
  return request({
    url: '/employee/scores',
    method: 'get',
    params
  })
}

/**
 * 获取课程类型成绩分析（用于雷达图）
 * @returns {Promise} 课程类型成绩分析数据
 */
export function getCourseTypeScores() {
  return request({
    url: '/employee/course-type-scores',
    method: 'get'
  })
}

/**
 * 获取员工学习进度
 * @returns {Promise} 学习进度数据
 */
export function getLearningProgress() {
  return request({
    url: '/employee/learning-progress',
    method: 'get'
  })
}
