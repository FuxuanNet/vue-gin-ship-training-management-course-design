package planner

import (
	"net/http"
	"time"
	"backend/database"

	"github.com/gin-gonic/gin"
)

// CreatePlanRequest 创建培训计划请求
type CreatePlanRequest struct {
	PlanName          string `json:"planName" binding:"required"`
	PlanStatus        string `json:"planStatus" binding:"required"`
	PlanStartDatetime string `json:"planStartDatetime" binding:"required"`
	PlanEndDatetime   string `json:"planEndDatetime" binding:"required"`
}

// CreatePlan 创建培训计划（5.2接口）
func CreatePlan(c *gin.Context) {
	// 获取当前用户ID
	personID, exists := c.Get("personId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
			"data":    nil,
		})
		return
	}

	// 绑定请求体
	var req CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误：" + err.Error(),
			"data":    nil,
		})
		return
	}

	// 验证计划状态
	validStatuses := map[string]bool{"规划中": true, "进行中": true, "已完成": true}
	if !validStatuses[req.PlanStatus] {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "计划状态必须为：规划中/进行中/已完成",
			"data":    nil,
		})
		return
	}

	// 解析时间
	startTime, err := time.Parse("2006-01-02 15:04:05", req.PlanStartDatetime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "开始时间格式错误，正确格式：YYYY-MM-DD HH:mm:ss",
			"data":    nil,
		})
		return
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", req.PlanEndDatetime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "结束时间格式错误，正确格式：YYYY-MM-DD HH:mm:ss",
			"data":    nil,
		})
		return
	}

	// 验证时间范围
	if !startTime.Before(endTime) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "开始时间必须早于结束时间",
			"data":    nil,
		})
		return
	}

	// 创建培训计划
	plan := database.TrainingPlan{
		PlanName:          req.PlanName,
		PlanStatus:        req.PlanStatus,
		PlanStartDatetime: startTime,
		PlanEndDatetime:   endTime,
		CreatorID:         personID.(int64),
	}

	if err := database.DB.Create(&plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建失败",
			"data":    nil,
		})
		return
	}

	// 查询制定人姓名
	var creator database.Person
	database.DB.Where("person_id = ?", plan.CreatorID).First(&creator)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data": gin.H{
			"planId":            plan.PlanID,
			"planName":          plan.PlanName,
			"planStatus":        plan.PlanStatus,
			"planStartDatetime": plan.PlanStartDatetime.Format("2006-01-02 15:04:05"),
			"planEndDatetime":   plan.PlanEndDatetime.Format("2006-01-02 15:04:05"),
			"creatorId":         plan.CreatorID,
			"creatorName":       creator.Name,
		},
	})
}
