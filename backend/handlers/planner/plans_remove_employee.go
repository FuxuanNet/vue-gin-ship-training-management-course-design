package planner

import (
	"net/http"
	"strconv"
	"backend/database"

	"github.com/gin-gonic/gin"
)

// RemoveEmployeeFromPlan 从培训计划移除员工（接口5.7）
func RemoveEmployeeFromPlan(c *gin.Context) {
	// 获取路径参数
	planIdStr := c.Param("planId")
	employeeIdStr := c.Param("employeeId")
	
	planId, err := strconv.ParseInt(planIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的培训计划ID",
			"data":    nil,
		})
		return
	}

	employeeId, err := strconv.ParseInt(employeeIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的员工ID",
			"data":    nil,
		})
		return
	}

	// 获取查询参数 force（是否强制删除）
	forceStr := c.Query("force")
	force := forceStr == "true"

	// 验证培训计划和员工关联是否存在
	var planEmployee database.PlanEmployee
	if err := database.DB.Where("plan_id = ? AND person_id = ?", planId, employeeId).First(&planEmployee).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "该员工不在此培训计划中",
			"data":    nil,
		})
		return
	}

	// 检查是否存在评价记录
	var evaluationCount int64
	database.DB.Model(&database.AttendanceEvaluation{}).
		Joins("JOIN plan_course_item ON attendance_evaluation.item_id = plan_course_item.item_id").
		Where("plan_course_item.plan_id = ? AND attendance_evaluation.person_id = ?", planId, employeeId).
		Count(&evaluationCount)

	// 如果有评价记录且不是强制删除，返回警告
	if evaluationCount > 0 && !force {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "该员工已有课程评价记录，建议不要移除。如需强制删除，请添加force=true参数",
			"data": gin.H{
				"evaluationCount": evaluationCount,
			},
		})
		return
	}

	// 如果是强制删除，先删除评价记录
	if force && evaluationCount > 0 {
		// 获取该员工在该计划下的所有课程安排ID
		var itemIds []int64
		database.DB.Model(&database.PlanCourseItem{}).
			Where("plan_id = ?", planId).
			Pluck("item_id", &itemIds)

		// 删除评价记录
		if len(itemIds) > 0 {
			database.DB.Where("item_id IN ? AND person_id = ?", itemIds, employeeId).
				Delete(&database.AttendanceEvaluation{})
		}
	}

	// 删除关联关系
	if err := database.DB.Where("plan_id = ? AND person_id = ?", planId, employeeId).
		Delete(&database.PlanEmployee{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "移除员工失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "移除成功",
		"data":    nil,
	})
}
