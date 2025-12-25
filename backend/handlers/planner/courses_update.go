package planner

import (
	"net/http"
	"strconv"
	"backend/database"

	"github.com/gin-gonic/gin"
)

// UpdateCourse 修改课程（接口5.10）
func UpdateCourse(c *gin.Context) {
	// 获取路径参数 courseId
	courseIdStr := c.Param("courseId")
	courseId, err := strconv.ParseInt(courseIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的课程ID",
			"data":    nil,
		})
		return
	}

	// 验证课程是否存在
	var course database.Course
	if err := database.DB.Where("course_id = ?", courseId).First(&course).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "课程不存在",
			"data":    nil,
		})
		return
	}

	// 解析请求体
	var req struct {
		CourseName    *string `json:"courseName"`
		CourseDesc    *string `json:"courseDesc"`
		CourseRequire *string `json:"courseRequire"`
		CourseClass   *string `json:"courseClass"`
		TeacherID     *int64  `json:"teacherId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误：" + err.Error(),
			"data":    nil,
		})
		return
	}

	// 构建更新映射
	updates := make(map[string]interface{})

	// 更新课程名称
	if req.CourseName != nil {
		if len(*req.CourseName) == 0 || len(*req.CourseName) > 50 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "课程名称长度必须在1-50字符之间",
				"data":    nil,
			})
			return
		}
		updates["course_name"] = *req.CourseName
	}

	// 更新课程描述
	if req.CourseDesc != nil {
		if len(*req.CourseDesc) > 100 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "课程描述长度不能超过100字符",
				"data":    nil,
			})
			return
		}
		updates["course_desc"] = *req.CourseDesc
	}

	// 更新课程要求
	if req.CourseRequire != nil {
		if len(*req.CourseRequire) > 500 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "课程要求长度不能超过500字符",
				"data":    nil,
			})
			return
		}
		updates["course_require"] = *req.CourseRequire
	}

	// 更新课程类型
	if req.CourseClass != nil {
		if len(*req.CourseClass) == 0 || len(*req.CourseClass) > 20 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "课程类型长度必须在1-20字符之间",
				"data":    nil,
			})
			return
		}
		updates["course_class"] = *req.CourseClass
	}

	// 更新讲师ID
	if req.TeacherID != nil {
		// 验证讲师是否存在且角色正确
		var teacher database.Person
		if err := database.DB.Where("person_id = ? AND role = ?", *req.TeacherID, "讲师").First(&teacher).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "讲师不存在或角色错误",
				"data":    nil,
			})
			return
		}
		updates["teacher_id"] = *req.TeacherID
	}

	// 如果没有更新内容，直接返回
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "没有提供更新内容",
			"data":    nil,
		})
		return
	}

	// 执行更新
	if err := database.DB.Model(&course).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新课程失败",
			"data":    nil,
		})
		return
	}

	// 重新查询更新后的课程（带讲师信息）
	database.DB.Preload("Teacher").Where("course_id = ?", courseId).First(&course)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "修改成功",
		"data": gin.H{
			"courseId":      course.CourseID,
			"courseName":    course.CourseName,
			"courseDesc":    course.CourseDesc,
			"courseRequire": course.CourseRequire,
			"courseClass":   course.CourseClass,
			"teacherId":     course.TeacherID,
			"teacherName":   course.Teacher.Name,
		},
	})
}
