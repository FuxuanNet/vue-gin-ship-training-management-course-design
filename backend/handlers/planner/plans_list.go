package planner

import (
	"net/http"
	"strconv"
	"backend/database"

	"github.com/gin-gonic/gin"
)

// PlanListItem 培训计划列表项
type PlanListItem struct {
	PlanID            int64  `json:"planId"`
	PlanName          string `json:"planName"`
	PlanStatus        string `json:"planStatus"`
	PlanStartDatetime string `json:"planStartDatetime"`
	PlanEndDatetime   string `json:"planEndDatetime"`
	CreatorID         int64  `json:"creatorId"`
	CreatorName       string `json:"creatorName"`
	EmployeeCount     int    `json:"employeeCount"`
	CourseCount       int    `json:"courseCount"`
}

// GetPlansList 获取培训计划列表（5.1接口）
func GetPlansList(c *gin.Context) {
	// 获取查询参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	status := c.Query("status")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	keyword := c.Query("keyword")
	sortBy := c.DefaultQuery("sortBy", "plan_start_datetime")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	// 转换分页参数
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 构建查询
	query := database.DB.Table("training_plan tp").
		Select(`
			tp.plan_id,
			tp.plan_name,
			tp.plan_status,
			DATE_FORMAT(tp.plan_start_datetime, '%Y-%m-%d %H:%i:%s') as plan_start_datetime,
			DATE_FORMAT(tp.plan_end_datetime, '%Y-%m-%d %H:%i:%s') as plan_end_datetime,
			tp.creator_id,
			p.name as creator_name,
			COALESCE(employee_counts.count, 0) as employee_count,
			COALESCE(course_counts.count, 0) as course_count
		`).
		Joins("JOIN person p ON tp.creator_id = p.person_id").
		Joins(`LEFT JOIN (
			SELECT plan_id, COUNT(*) as count 
			FROM plan_employee 
			GROUP BY plan_id
		) employee_counts ON tp.plan_id = employee_counts.plan_id`).
		Joins(`LEFT JOIN (
			SELECT plan_id, COUNT(*) as count 
			FROM plan_course_item 
			GROUP BY plan_id
		) course_counts ON tp.plan_id = course_counts.plan_id`)

	// 添加筛选条件
	if status != "" {
		query = query.Where("tp.plan_status = ?", status)
	}
	if startDate != "" {
		query = query.Where("DATE(tp.plan_start_datetime) >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("DATE(tp.plan_end_datetime) <= ?", endDate)
	}
	if keyword != "" {
		query = query.Where("tp.plan_name LIKE ?", "%"+keyword+"%")
	}

	// 统计总数
	var total int64
	countQuery := query
	if err := countQuery.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
			"data":    nil,
		})
		return
	}

	// 添加排序
	orderClause := "tp." + sortBy + " " + sortOrder
	query = query.Order(orderClause)

	// 添加分页
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// 执行查询
	var plans []PlanListItem
	if err := query.Scan(&plans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
			"data":    nil,
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
			"list":     plans,
		},
	})
}
