package planner

import (
	"net/http"
	"strconv"
	"backend/database"

	"github.com/gin-gonic/gin"

)

// GetCourseItemsList 获取课程安排列表（接口5.12）
func GetCourseItemsList(c *gin.Context) {
	// 获取查询参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "20")
	planIdStr := c.Query("planId")
	courseIdStr := c.Query("courseId")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	sortBy := c.DefaultQuery("sortBy", "class_date")
	sortOrder := c.DefaultQuery("sortOrder", "asc")

	// 解析分页参数
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 20
	}

	// 构建查询
	query := database.DB.Model(&database.PlanCourseItem{}).
		Joins("LEFT JOIN training_plan ON plan_course_item.plan_id = training_plan.plan_id").
		Joins("LEFT JOIN course ON plan_course_item.course_id = course.course_id").
		Joins("LEFT JOIN person ON course.teacher_id = person.person_id")

	// 筛选条件：培训计划
	if planIdStr != "" {
		planId, err := strconv.ParseInt(planIdStr, 10, 64)
		if err == nil {
			query = query.Where("plan_course_item.plan_id = ?", planId)
		}
	}

	// 筛选条件：课程
	if courseIdStr != "" {
		courseId, err := strconv.ParseInt(courseIdStr, 10, 64)
		if err == nil {
			query = query.Where("plan_course_item.course_id = ?", courseId)
		}
	}

	// 筛选条件：日期范围
	if startDate != "" {
		query = query.Where("plan_course_item.class_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("plan_course_item.class_date <= ?", endDate)
	}

	// 查询总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
			"data":    nil,
		})
		return
	}

	// 排序
	orderClause := "plan_course_item." + sortBy + " " + sortOrder
	query = query.Order(orderClause)

	// 查询课程安排列表
	var items []database.PlanCourseItem
	offset := (page - 1) * pageSize
	if err := query.Preload("Plan").Preload("Course").Preload("Course.Teacher").
		Offset(offset).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
			"data":    nil,
		})
		return
	}

	// 构建响应数据
	type ItemResponse struct {
		ItemID         int64  `json:"itemId"`
		PlanID         int64  `json:"planId"`
		PlanName       string `json:"planName"`
		CourseID       int64  `json:"courseId"`
		CourseName     string `json:"courseName"`
		CourseClass    string `json:"courseClass"`
		TeacherID      int64  `json:"teacherId"`
		TeacherName    string `json:"teacherName"`
		ClassDate      string `json:"classDate"`
		ClassBeginTime string `json:"classBeginTime"`
		ClassEndTime   string `json:"classEndTime"`
		Location       string `json:"location"`
	}

	list := make([]ItemResponse, 0, len(items))
	for _, item := range items {
		list = append(list, ItemResponse{
			ItemID:         item.ItemID,
			PlanID:         item.PlanID,
			PlanName:       item.Plan.PlanName,
			CourseID:       item.CourseID,
			CourseName:     item.Course.CourseName,
			CourseClass:    item.Course.CourseClass,
			TeacherID:      item.Course.TeacherID,
			TeacherName:    item.Course.Teacher.Name,
			ClassDate:      item.ClassDate.Format("2006-01-02"),
			ClassBeginTime: item.ClassBeginTime.Format("15:04:05"),
			ClassEndTime:   item.ClassEndTime.Format("15:04:05"),
			Location:       item.Location,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
			"list":     list,
		},
	})
}
