package planner

import (
	"net/http"
	"strconv"
	"backend/database"

	"github.com/gin-gonic/gin"
)

// AddEmployeesToPlan 为培训计划添加员工（接口5.6）
func AddEmployeesToPlan(c *gin.Context) {
	// 获取路径参数 planId
	planIdStr := c.Param("planId")
	planId, err := strconv.ParseInt(planIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的培训计划ID",
			"data":    nil,
		})
		return
	}

	// 解析请求体
	var req struct {
		EmployeeIds []int64 `json:"employeeIds" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误：" + err.Error(),
			"data":    nil,
		})
		return
	}

	// 验证员工ID列表不为空
	if len(req.EmployeeIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "员工ID列表不能为空",
			"data":    nil,
		})
		return
	}

	// 验证培训计划是否存在
	var plan database.TrainingPlan
	if err := database.DB.Where("plan_id = ?", planId).First(&plan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "培训计划不存在",
			"data":    nil,
		})
		return
	}

	// 验证所有员工ID的合法性
	var persons []database.Person
	if err := database.DB.Where("person_id IN ? AND role = ?", req.EmployeeIds, "员工").Find(&persons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询员工信息失败",
			"data":    nil,
		})
		return
	}

	// 检查是否所有员工都存在且角色正确
	if len(persons) != len(req.EmployeeIds) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "部分员工不存在或角色不是员工",
			"data":    nil,
		})
		return
	}

	// 查询已关联的员工
	var existingRelations []database.PlanEmployee
	database.DB.Where("plan_id = ? AND person_id IN ?", planId, req.EmployeeIds).Find(&existingRelations)
	
	// 构建已存在的员工ID映射
	existingMap := make(map[int64]bool)
	for _, relation := range existingRelations {
		existingMap[relation.PersonID] = true
	}

	// 批量插入新的关联关系
	addedCount := 0
	skippedCount := 0
	for _, employeeId := range req.EmployeeIds {
		if existingMap[employeeId] {
			skippedCount++
			continue
		}

		// 插入新关联
		planEmployee := database.PlanEmployee{
			PlanID:   planId,
			PersonID: employeeId,
		}
		if err := database.DB.Create(&planEmployee).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "添加员工失败",
				"data":    nil,
			})
			return
		}
		addedCount++
	}

	// 返回结果
	message := "添加成功"
	if skippedCount > 0 {
		message = "添加完成，部分员工已存在"
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": message,
		"data": gin.H{
			"addedCount":   addedCount,
			"skippedCount": skippedCount,
		},
	})
}
