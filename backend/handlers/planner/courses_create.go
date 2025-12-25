package planner

import (
	"net/http"
	"strings"
	"backend/database"

	"github.com/gin-gonic/gin"
)

// CreateCourse 创建课程（接口5.9）
func CreateCourse(c *gin.Context) {
	// 解析请求体
	var req struct {
		CourseName    string `json:"courseName" binding:"required"`
		CourseDesc    string `json:"courseDesc"`
		CourseRequire string `json:"courseRequire"`
		CourseClass   string `json:"courseClass" binding:"required"`
		TeacherID     int64  `json:"teacherId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误：" + err.Error(),
			"data":    nil,
		})
		return
	}

	// 验证字段长度
	if len(strings.TrimSpace(req.CourseName)) == 0 || len(req.CourseName) > 50 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "课程名称长度必须在1-50字符之间",
			"data":    nil,
		})
		return
	}

	if len(req.CourseDesc) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "课程描述长度不能超过100字符",
			"data":    nil,
		})
		return
	}

	if len(req.CourseRequire) > 500 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "课程要求长度不能超过500字符",
			"data":    nil,
		})
		return
	}

	if len(strings.TrimSpace(req.CourseClass)) == 0 || len(req.CourseClass) > 20 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "课程类型长度必须在1-20字符之间",
			"data":    nil,
		})
		return
	}

	// 验证讲师是否存在且角色正确
	var teacher database.Person
	if err := database.DB.Where("person_id = ? AND role = ?", req.TeacherID, "讲师").First(&teacher).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "讲师不存在或角色错误",
			"data":    nil,
		})
		return
	}

	// 创建课程
	course := database.Course{
		CourseName:    req.CourseName,
		CourseDesc:    req.CourseDesc,
		CourseRequire: req.CourseRequire,
		CourseClass:   req.CourseClass,
		TeacherID:     req.TeacherID,
	}

	if err := database.DB.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建课程失败",
			"data":    nil,
		})
		return
	}

	// 返回创建的课程信息
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data": gin.H{
			"courseId":      course.CourseID,
			"courseName":    course.CourseName,
			"courseDesc":    course.CourseDesc,
			"courseRequire": course.CourseRequire,
			"courseClass":   course.CourseClass,
			"teacherId":     course.TeacherID,
			"teacherName":   teacher.Name,
		},
	})
}
