package planner

import (
	"net/http"
	"strconv"
	"backend/database"

	"github.com/gin-gonic/gin"
)

// GetCoursesList 获取课程列表（接口5.8）
func GetCoursesList(c *gin.Context) {
	// 获取查询参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	courseClass := c.Query("courseClass")
	keyword := c.Query("keyword")
	teacherIdStr := c.Query("teacherId")

	// 解析分页参数
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// 构建查询
	query := database.DB.Model(&database.Course{})

	// 筛选条件：课程类型
	if courseClass != "" {
		query = query.Where("course_class = ?", courseClass)
	}

	// 筛选条件：关键词搜索（课程名称）
	if keyword != "" {
		query = query.Where("course_name LIKE ?", "%"+keyword+"%")
	}

	// 筛选条件：讲师ID
	if teacherIdStr != "" {
		teacherId, err := strconv.ParseInt(teacherIdStr, 10, 64)
		if err == nil {
			query = query.Where("teacher_id = ?", teacherId)
		}
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

	// 查询课程列表（带讲师信息）
	var courses []database.Course
	offset := (page - 1) * pageSize
	if err := query.Preload("Teacher").
		Offset(offset).
		Limit(pageSize).
		Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
			"data":    nil,
		})
		return
	}

	// 构建响应数据
	type CourseResponse struct {
		CourseID       int64  `json:"courseId"`
		CourseName     string `json:"courseName"`
		CourseDesc     string `json:"courseDesc"`
		CourseRequire  string `json:"courseRequire"`
		CourseClass    string `json:"courseClass"`
		TeacherID      int64  `json:"teacherId"`
		TeacherName    string `json:"teacherName"`
		ScheduledCount int64  `json:"scheduledCount"`
	}

	list := make([]CourseResponse, 0, len(courses))
	for _, course := range courses {
		// 查询该课程的安排次数
		var scheduledCount int64
		database.DB.Model(&database.PlanCourseItem{}).
			Where("course_id = ?", course.CourseID).
			Count(&scheduledCount)

		list = append(list, CourseResponse{
			CourseID:       course.CourseID,
			CourseName:     course.CourseName,
			CourseDesc:     course.CourseDesc,
			CourseRequire:  course.CourseRequire,
			CourseClass:    course.CourseClass,
			TeacherID:      course.TeacherID,
			TeacherName:    course.Teacher.Name,
			ScheduledCount: scheduledCount,
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
