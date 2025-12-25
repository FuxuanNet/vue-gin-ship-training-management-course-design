package planner

import (
	"net/http"
	"strconv"
	"backend/database"

	"github.com/gin-gonic/gin"
)

// DeletePlan 删除培训计划（5.4接口）
func DeletePlan(c *gin.Context) {
	// 获取路径参数
	planIDStr := c.Param("planId")
	planID, err := strconv.ParseInt(planIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的计划ID",
			"data":    nil,
		})
		return
	}

	// 查询计划是否存在
	var plan database.TrainingPlan
	if err := database.DB.Where("plan_id = ?", planID).First(&plan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "培训计划不存在",
			"data":    nil,
		})
		return
	}

	// 检查是否存在关联的课程安排
	var courseItemCount int64
	database.DB.Model(&database.PlanCourseItem{}).Where("plan_id = ?", planID).Count(&courseItemCount)
	if courseItemCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无法删除，该计划下存在课程安排或关联员工，请先删除相关数据",
			"data":    nil,
		})
		return
	}

	// 检查是否存在关联的员工
	var employeeCount int64
	database.DB.Model(&database.PlanEmployee{}).Where("plan_id = ?", planID).Count(&employeeCount)
	if employeeCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无法删除，该计划下存在课程安排或关联员工，请先删除相关数据",
			"data":    nil,
		})
		return
	}

	// 删除培训计划
	if err := database.DB.Delete(&plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除失败",
			"data":    nil,
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
		"data":    nil,
	})
}
