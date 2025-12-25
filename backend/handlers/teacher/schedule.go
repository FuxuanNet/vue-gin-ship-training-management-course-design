package teacher

import (
	"net/http"
	"time"

	"backend/database"

	"github.com/gin-gonic/gin"
)

// GetSchedule 获取讲师授课表
func GetSchedule(c *gin.Context) {
	// 获取当前用户信息
	personID, exists := c.Get("personId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
			"data":    nil,
		})
		return
	}
	
	teacherID := personID.(int64)

	// 获取查询参数
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少必填参数：startDate 或 endDate",
			"data":    nil,
		})
		return
	}

	// 验证日期格式
	startTime, err1 := time.Parse("2006-01-02", startDate)
	endTime, err2 := time.Parse("2006-01-02", endDate)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "日期格式错误，应为 YYYY-MM-DD",
			"data":    nil,
		})
		return
	}

	// 查询指定日期范围内的授课课程
	var courseItems []database.PlanCourseItem
	err := database.DB.
		Preload("Course").
		Preload("Plan").
		Joins("JOIN course ON plan_course_item.course_id = course.course_id").
		Where("course.teacher_id = ? AND plan_course_item.class_date BETWEEN ? AND ?", 
			teacherID, startDate, endDate).
		Order("plan_course_item.class_date ASC, plan_course_item.class_begin_time ASC").
		Find(&courseItems).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询授课表失败",
			"data":    nil,
		})
		return
	}

	// 统计每个课程的学员数和已评分数
	itemIDs := make([]int64, len(courseItems))
	for i, item := range courseItems {
		itemIDs[i] = item.ItemID
	}

	// 批量查询学员数和已评分数
	type CountResult struct {
		ItemID         int64
		StudentCount   int64
		EvaluatedCount int64
	}
	var countResults []CountResult
	if len(itemIDs) > 0 {
		database.DB.Raw(`
			SELECT 
				item_id,
				COUNT(*) as student_count,
				SUM(CASE WHEN teacher_comment IS NOT NULL AND teacher_comment != '' THEN 1 ELSE 0 END) as evaluated_count
			FROM attendance_evaluation
			WHERE item_id IN ?
			GROUP BY item_id
		`, itemIDs).Scan(&countResults)
	}

	// 构建统计映射
	countMap := make(map[int64]CountResult)
	for _, result := range countResults {
		countMap[result.ItemID] = result
	}

	// 按日期分组课程
	scheduleMap := make(map[string][]map[string]interface{})
	weekDays := []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
	
	for _, item := range courseItems {
		dateStr := item.ClassDate.Format("2006-01-02")
		
		counts := countMap[item.ItemID]
		
		course := map[string]interface{}{
			"itemId":         item.ItemID,
			"courseId":       item.CourseID,
			"courseName":     item.Course.CourseName,
			"courseDesc":     item.Course.CourseDesc,
			"courseRequire":  item.Course.CourseRequire,
			"courseClass":    item.Course.CourseClass,
			"classBeginTime": item.ClassBeginTime.Format("15:04:05"),
			"classEndTime":   item.ClassEndTime.Format("15:04:05"),
			"location":       item.Location,
			"planId":         item.PlanID,
			"planName":       item.Plan.PlanName,
			"studentCount":   int(counts.StudentCount),
			"evaluatedCount": int(counts.EvaluatedCount),
		}
		
		scheduleMap[dateStr] = append(scheduleMap[dateStr], course)
	}

	// 构建完整日程（包含没有课程的日期）
	schedule := []map[string]interface{}{}
	for d := startTime; !d.After(endTime); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		dayOfWeek := weekDays[d.Weekday()]
		courses := scheduleMap[dateStr]
		if courses == nil {
			courses = []map[string]interface{}{}
		}

		schedule = append(schedule, map[string]interface{}{
			"date":      dateStr,
			"dayOfWeek": dayOfWeek,
			"courses":   courses,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"startDate":    startDate,
			"endDate":      endDate,
			"totalCourses": len(courseItems),
			"schedule":     schedule,
		},
	})
}
