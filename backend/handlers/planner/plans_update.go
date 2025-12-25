package planner

import (
	"net/http"
	"strconv"
	"time"
	"backend/database"

	"github.com/gin-gonic/gin"
)

// UpdatePlanRequest 更新培训计划请求
type UpdatePlanRequest struct {
	PlanName          string `json:"planName"`
	PlanStatus        string `json:"planStatus"`
	PlanStartDatetime string `json:"planStartDatetime"`
	PlanEndDatetime   string `json:"planEndDatetime"`
}

// UpdatePlan 修改培训计划（5.3接口）
func UpdatePlan(c *gin.Context) {
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

	// 绑定请求体
	var req UpdatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误：" + err.Error(),
			"data":    nil,
		})
		return
	}

	// 更新字段
	updates := make(map[string]interface{})

	if req.PlanName != "" {
		updates["plan_name"] = req.PlanName
	}

	if req.PlanStatus != "" {
		validStatuses := map[string]bool{"规划中": true, "进行中": true, "已完成": true}
		if !validStatuses[req.PlanStatus] {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "计划状态必须为：规划中/进行中/已完成",
				"data":    nil,
			})
			return
		}
		updates["plan_status"] = req.PlanStatus
	}

	// 处理开始时间
	var startTime time.Time
	if req.PlanStartDatetime != "" {
		startTime, err = time.Parse("2006-01-02 15:04:05", req.PlanStartDatetime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "开始时间格式错误，正确格式：YYYY-MM-DD HH:mm:ss",
				"data":    nil,
			})
			return
		}
		updates["plan_start_datetime"] = startTime
	} else {
		startTime = plan.PlanStartDatetime
	}

	// 处理结束时间
	var endTime time.Time
	if req.PlanEndDatetime != "" {
		endTime, err = time.Parse("2006-01-02 15:04:05", req.PlanEndDatetime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "结束时间格式错误，正确格式：YYYY-MM-DD HH:mm:ss",
				"data":    nil,
			})
			return
		}
		updates["plan_end_datetime"] = endTime
	} else {
		endTime = plan.PlanEndDatetime
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

	// 更新数据库
	if len(updates) > 0 {
		if err := database.DB.Model(&plan).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "修改失败",
				"data":    nil,
			})
			return
		}
	}

	// 重新查询更新后的计划
	database.DB.Where("plan_id = ?", planID).First(&plan)

	// 查询制定人姓名
	var creator database.Person
	database.DB.Where("person_id = ?", plan.CreatorID).First(&creator)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "修改成功",
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
