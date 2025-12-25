// 课程大纲制定者相关接口
import request from './request'

/**
 * 获取培训计划列表
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码，默认1
 * @param {number} params.pageSize - 每页数量，默认10
 * @param {string} params.status - 计划状态筛选：规划中/进行中/已完成
 * @param {string} params.startDate - 开始日期筛选（YYYY-MM-DD）
 * @param {string} params.endDate - 结束日期筛选（YYYY-MM-DD）
 * @param {string} params.keyword - 计划名称关键词搜索
 * @param {string} params.sortBy - 排序字段
 * @param {string} params.sortOrder - 排序方向：asc/desc
 * @returns {Promise} 培训计划列表数据
 */
export function getPlansList(params) {
  return request({
    url: '/planner/plans',
    method: 'get',
    params
  })
}

/**
 * 创建培训计划
 * @param {Object} data - 培训计划数据
 * @param {string} data.planName - 计划名称
 * @param {string} data.planStatus - 计划状态：规划中/进行中/已完成
 * @param {string} data.planStartDatetime - 开始时间（YYYY-MM-DD HH:mm:ss）
 * @param {string} data.planEndDatetime - 结束时间（YYYY-MM-DD HH:mm:ss）
 * @returns {Promise} 创建结果
 */
export function createPlan(data) {
  return request({
    url: '/planner/plans',
    method: 'post',
    data
  })
}

/**
 * 获取培训计划详情
 * @param {number} planId - 培训计划ID
 * @returns {Promise} 培训计划详情数据
 */
export function getPlanDetail(planId) {
  return request({
    url: `/planner/plans/${planId}`,
    method: 'get'
  })
}

/**
 * 修改培训计划
 * @param {number} planId - 培训计划ID
 * @param {Object} data - 要修改的数据
 * @param {string} data.planName - 计划名称（可选）
 * @param {string} data.planStatus - 计划状态（可选）
 * @param {string} data.planStartDatetime - 开始时间（可选）
 * @param {string} data.planEndDatetime - 结束时间（可选）
 * @returns {Promise} 修改结果
 */
export function updatePlan(planId, data) {
  return request({
    url: `/planner/plans/${planId}`,
    method: 'put',
    data
  })
}

/**
 * 删除培训计划
 * @param {number} planId - 培训计划ID
 * @returns {Promise} 删除结果
 */
export function deletePlan(planId) {
  return request({
    url: `/planner/plans/${planId}`,
    method: 'delete'
  })
}

/**
 * 为培训计划添加员工
 * @param {number} planId - 培训计划ID
 * @param {Object} data - 员工数据
 * @param {Array<number>} data.employeeIds - 员工ID数组
 * @returns {Promise} 添加结果
 */
export function addEmployeesToPlan(planId, data) {
  return request({
    url: `/planner/plans/${planId}/employees`,
    method: 'post',
    data
  })
}

/**
 * 从培训计划移除员工
 * @param {number} planId - 培训计划ID
 * @param {number} employeeId - 员工ID
 * @param {boolean} force - 是否强制删除
 * @returns {Promise} 移除结果
 */
export function removeEmployeeFromPlan(planId, employeeId, force = false) {
  return request({
    url: `/planner/plans/${planId}/employees/${employeeId}`,
    method: 'delete',
    params: { force }
  })
}

/**
 * 获取课程列表
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码，默认1
 * @param {number} params.pageSize - 每页数量，默认10
 * @param {string} params.courseClass - 课程类型筛选
 * @param {string} params.keyword - 课程名称关键词搜索
 * @param {number} params.teacherId - 按讲师筛选
 * @returns {Promise} 课程列表数据
 */
export function getCoursesList(params) {
  return request({
    url: '/planner/courses',
    method: 'get',
    params
  })
}

/**
 * 创建课程
 * @param {Object} data - 课程数据
 * @param {string} data.courseName - 课程名称
 * @param {string} data.courseDesc - 课程描述（可选）
 * @param {string} data.courseRequire - 课程要求（可选）
 * @param {string} data.courseClass - 课程类型
 * @param {number} data.teacherId - 讲师ID
 * @returns {Promise} 创建结果
 */
export function createCourse(data) {
  return request({
    url: '/planner/courses',
    method: 'post',
    data
  })
}

/**
 * 修改课程
 * @param {number} courseId - 课程ID
 * @param {Object} data - 要修改的数据
 * @param {string} data.courseName - 课程名称（可选）
 * @param {string} data.courseDesc - 课程描述（可选）
 * @param {string} data.courseRequire - 课程要求（可选）
 * @param {string} data.courseClass - 课程类型（可选）
 * @param {number} data.teacherId - 讲师ID（可选）
 * @returns {Promise} 修改结果
 */
export function updateCourse(courseId, data) {
  return request({
    url: `/planner/courses/${courseId}`,
    method: 'put',
    data
  })
}

/**
 * 删除课程
 * @param {number} courseId - 课程ID
 * @returns {Promise} 删除结果
 */
export function deleteCourse(courseId) {
  return request({
    url: `/planner/courses/${courseId}`,
    method: 'delete'
  })
}
