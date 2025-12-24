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
	personID, exists := c.Get("person_id")
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

	// 验证日期格式
	_, err1 := time.Parse("2006-01-02", startDate)
	_, err2 := time.Parse("2006-01-02", endDate)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "日期格式错误，应为YYYY-MM-DD", "data": nil})
		return
	}

	// 查询课程表
	rows, err := database.DB.Raw(`
		SELECT 
			pci.item_id,
			pci.course_id,
			c.course_name,
			c.course_desc,
			c.course_class,
			pci.class_date,
			pci.class_begin_time,
			pci.class_end_time,
			pci.location,
			pci.plan_id,
			tp.plan_name,
			c.teacher_id,
			p.name AS teacher_name,
			CASE 
				WHEN ae.person_id IS NOT NULL THEN true 
				ELSE false 
			END AS has_evaluated
		FROM plan_employee pe
		JOIN plan_course_item pci ON pe.plan_id = pci.plan_id
		JOIN course c ON pci.course_id = c.course_id
		JOIN training_plan tp ON pci.plan_id = tp.plan_id
		JOIN person p ON c.teacher_id = p.person_id
		LEFT JOIN attendance_evaluation ae ON ae.item_id = pci.item_id AND ae.person_id = ?
		WHERE pe.person_id = ? 
			AND pci.class_date >= ? 
			AND pci.class_date <= ?
		ORDER BY pci.class_date ASC, pci.class_begin_time ASC
	`, personID, personID, startDate, endDate).Rows()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "data": nil})
		return
	}
	defer rows.Close()

	// 按日期组织课程
	scheduleMap := make(map[string][]map[string]interface{})
	totalCourses := 0

	for rows.Next() {
		var (
			itemID         int64
			courseID       int64
			courseName     string
			courseDesc     string
			courseClass    string
			classDate      string
			classBeginTime string
			classEndTime   string
			location       string
			planID         int64
			planName       string
			teacherID      int64
			teacherName    string
			hasEvaluated   bool
		)

		rows.Scan(
			&itemID, &courseID, &courseName, &courseDesc, &courseClass,
			&classDate, &classBeginTime, &classEndTime, &location,
			&planID, &planName, &teacherID, &teacherName, &hasEvaluated,
		)

		course := map[string]interface{}{
			"itemId":         itemID,
			"courseId":       courseID,
			"courseName":     courseName,
			"courseDesc":     courseDesc,
			"courseClass":    courseClass,
			"classBeginTime": classBeginTime,
			"classEndTime":   classEndTime,
			"location":       location,
			"planId":         planID,
			"planName":       planName,
			"teacherId":      teacherID,
			"teacherName":    teacherName,
			"hasEvaluated":   hasEvaluated,
		}

		scheduleMap[classDate] = append(scheduleMap[classDate], course)
		totalCourses++
	}

	// 构建返回数据（包含所有日期，即使没有课程）
	schedule := []map[string]interface{}{}
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
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
			"totalCourses": totalCourses,
			"schedule":     schedule,
		},
	})
}
