package employee

import (
	"backend/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetSchedule 获取员工课程表
func GetSchedule(c *gin.Context) {
	// 从中间件获取用户ID
	personID, exists := c.Get("personId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录", "data": nil})
		return
	}

	// 获取查询参数
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// 验证日期参数
	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "日期参数错误", "data": nil})
		return
	}

	// 验证日期格式并转换
	startTime, err1 := time.Parse("2006-01-02", startDate)
	endTime, err2 := time.Parse("2006-01-02", endDate)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "日期格式错误，应为YYYY-MM-DD", "data": nil})
		return
	}

	// 设置时间范围
	startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
	endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 999999999, endTime.Location())

	// 1. 查询员工参与的培训计划ID
	var planEmployees []database.PlanEmployee
	if err := database.DB.Where("person_id = ?", personID).Find(&planEmployees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "data": nil})
		return
	}

	planIDs := make([]int64, 0, len(planEmployees))
	for _, pe := range planEmployees {
		planIDs = append(planIDs, pe.PlanID)
	}

	if len(planIDs) == 0 {
		// 没有参与任何培训计划，返回空日程
		schedule := buildEmptySchedule(startTime, endTime)
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取成功",
			"data": gin.H{
				"startDate":    startDate,
				"endDate":      endDate,
				"totalCourses": 0,
				"courses":      schedule,
			},
		})
		return
	}

	// 2. 查询日期范围内的课程安排（预加载关联）
	var courseItems []database.PlanCourseItem
	if err := database.DB.
		Preload("Course.Teacher").
		Preload("Plan").
		Where("plan_id IN ?", planIDs).
		Where("class_date >= ? AND class_date <= ?", startTime, endTime).
		Order("class_date ASC, class_begin_time ASC").
		Find(&courseItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "data": nil})
		return
	}

	// 3. 查询该员工的评价记录
	itemIDs := make([]int64, 0, len(courseItems))
	for _, item := range courseItems {
		itemIDs = append(itemIDs, item.ItemID)
	}

	evaluatedMap := make(map[int64]bool)
	if len(itemIDs) > 0 {
		var evaluations []database.AttendanceEvaluation
		database.DB.Where("person_id = ? AND item_id IN ?", personID, itemIDs).Find(&evaluations)
		for _, eval := range evaluations {
			evaluatedMap[eval.ItemID] = true
		}
	}

	// 4. 按日期组织课程
	scheduleMap := make(map[string][]map[string]interface{})
	for _, item := range courseItems {
		dateStr := item.ClassDate.Format("2006-01-02")
		
		course := map[string]interface{}{
			"itemId":         item.ItemID,
			"courseId":       item.CourseID,
			"courseName":     item.Course.CourseName,
			"courseDesc":     item.Course.CourseDesc,
			"courseClass":    item.Course.CourseClass,
			"classBeginTime": item.ClassBeginTime.Format("15:04:05"),
			"classEndTime":   item.ClassEndTime.Format("15:04:05"),
			"location":       item.Location,
			"planId":         item.PlanID,
			"planName":       item.Plan.PlanName,
			"teacherId":      item.Course.TeacherID,
			"teacherName":    item.Course.Teacher.Name,
			"hasEvaluated":   evaluatedMap[item.ItemID],
		}
		
		scheduleMap[dateStr] = append(scheduleMap[dateStr], course)
	}

	// 5. 构建完整日程（包含没有课程的日期）
	schedule := []map[string]interface{}{}
	for d := startTime; !d.After(endTime); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		courses := scheduleMap[dateStr]
		if courses == nil {
			courses = []map[string]interface{}{}
		}

		schedule = append(schedule, map[string]interface{}{
			"date":        dateStr,
			"courseCount": len(courses),
			"courses":     courses,
		})
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"startDate":    startDate,
			"endDate":      endDate,
			"totalCourses": len(courseItems),
			"courses":      schedule,
		},
	})
}

// buildEmptySchedule 构建空日程
func buildEmptySchedule(startTime, endTime time.Time) []map[string]interface{} {
	schedule := []map[string]interface{}{}
	for d := startTime; !d.After(endTime); d = d.AddDate(0, 0, 1) {
		schedule = append(schedule, map[string]interface{}{
			"date":        d.Format("2006-01-02"),
			"courseCount": 0,
			"courses":     []interface{}{},
		})
	}
	return schedule
}
