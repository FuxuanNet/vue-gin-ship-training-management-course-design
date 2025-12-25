package planner

import (
	"net/http"
	"time"
	"backend/database"

	"github.com/gin-gonic/gin"
)

// CreateCourseItem 创建课程安排（接口5.13）
func CreateCourseItem(c *gin.Context) {
	// 解析请求体
	var req struct {
		PlanID         int64  `json:"planId" binding:"required"`
		CourseID       int64  `json:"courseId" binding:"required"`
		ClassDate      string `json:"classDate" binding:"required"`
		ClassBeginTime string `json:"classBeginTime" binding:"required"`
		ClassEndTime   string `json:"classEndTime" binding:"required"`
		Location       string `json:"location" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误：" + err.Error(),
			"data":    nil,
		})
		return
	}

	// 验证地点长度
	if len(req.Location) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "上课地点长度不能超过100字符",
			"data":    nil,
		})
		return
	}

	// 验证培训计划是否存在
	var plan database.TrainingPlan
	if err := database.DB.Where("plan_id = ?", req.PlanID).First(&plan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "培训计划不存在",
			"data":    nil,
		})
		return
	}

	// 验证课程是否存在并获取讲师信息
	var course database.Course
	if err := database.DB.Preload("Teacher").Where("course_id = ?", req.CourseID).First(&course).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "课程不存在",
			"data":    nil,
		})
		return
	}

	// 检查讲师时间冲突
	var conflictCount int64
	database.DB.Model(&database.PlanCourseItem{}).
		Joins("LEFT JOIN course ON plan_course_item.course_id = course.course_id").
		Where("course.teacher_id = ?", course.TeacherID).
		Where("plan_course_item.class_date = ?", req.ClassDate).
		Where("(plan_course_item.class_begin_time < ? AND plan_course_item.class_end_time > ?) OR "+
			"(plan_course_item.class_begin_time < ? AND plan_course_item.class_end_time > ?) OR "+
			"(plan_course_item.class_begin_time >= ? AND plan_course_item.class_end_time <= ?)",
			req.ClassEndTime, req.ClassBeginTime,
			req.ClassEndTime, req.ClassBeginTime,
			req.ClassBeginTime, req.ClassEndTime).
		Count(&conflictCount)

	if conflictCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "时间冲突：讲师" + course.Teacher.Name + "在该时间段已有其他课程安排",
			"data":    nil,
		})
		return
	}

	// 解析日期
	classDate, err := time.Parse("2006-01-02", req.ClassDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "日期格式错误，请使用 YYYY-MM-DD 格式",
			"data":    nil,
		})
		return
	}

	// 验证时间格式（HH:mm:ss）
	if _, err := time.Parse("15:04:05", req.ClassBeginTime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "开始时间格式错误，请使用 HH:mm:ss 格式",
			"data":    nil,
		})
		return
	}

	if _, err := time.Parse("15:04:05", req.ClassEndTime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "结束时间格式错误，请使用 HH:mm:ss 格式",
			"data":    nil,
		})
		return
	}

	// 创建课程安排
	item := database.PlanCourseItem{
		PlanID:         req.PlanID,
		CourseID:       req.CourseID,
		ClassDate:      classDate,
		ClassBeginTime: req.ClassBeginTime,
		ClassEndTime:   req.ClassEndTime,
		Location:       req.Location,
	}

	if err := database.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建课程安排失败",
			"data":    nil,
		})
		return
	}

	// 返回创建的课程安排信息
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data": gin.H{
			"itemId":         item.ItemID,
			"planId":         item.PlanID,
			"planName":       plan.PlanName,
			"courseId":       item.CourseID,
			"courseName":     course.CourseName,
			"classDate":      item.ClassDate,
			"classBeginTime": item.ClassBeginTime,
			"classEndTime":   item.ClassEndTime,
			"location":       item.Location,
		},
	})
}
