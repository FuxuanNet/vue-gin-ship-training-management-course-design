package planner

import (
	"net/http"
	"strconv"
	"backend/database"

	"github.com/gin-gonic/gin"
)

// DeleteCourse 删除课程（接口5.11）
func DeleteCourse(c *gin.Context) {
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

	// 检查是否有课程安排（plan_course_item）
	var scheduledCount int64
	database.DB.Model(&database.PlanCourseItem{}).
		Where("course_id = ?", courseId).
		Count(&scheduledCount)

	if scheduledCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无法删除，该课程已被安排到培训计划中",
			"data": gin.H{
				"scheduledCount": scheduledCount,
			},
		})
		return
	}

	// 删除课程
	if err := database.DB.Delete(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除课程失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
		"data":    nil,
	})
}
