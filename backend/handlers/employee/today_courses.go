package employee

import (
	"backend/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetTodayCourses 获取员工今日课程列表
func GetTodayCourses(c *gin.Context) {
	// 从中间件获取用户ID
	personID, exists := c.Get("person_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录", "data": nil})
		return
	}

	// 获取今天的日期
	today := time.Now().Format("2006-01-02")

	// 查询今日课程
	var courses []map[string]interface{}
	
	// SQL查询：通过plan_employee关联查找该员工今日的课程
	rows, err := database.DB.Raw(`
		SELECT 
			pci.item_id,
			pci.course_id,
			c.course_name,
			c.course_desc,
			c.course_require,
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
			END AS has_evaluated,
			CASE 
				WHEN pci.class_date < CURDATE() THEN '已完成'
				WHEN pci.class_date = CURDATE() THEN '待上课'
				ELSE '待上课'
			END AS status
		FROM plan_employee pe
		JOIN plan_course_item pci ON pe.plan_id = pci.plan_id
		JOIN course c ON pci.course_id = c.course_id
		JOIN training_plan tp ON pci.plan_id = tp.plan_id
		JOIN person p ON c.teacher_id = p.person_id
		LEFT JOIN attendance_evaluation ae ON ae.item_id = pci.item_id AND ae.person_id = ?
		WHERE pe.person_id = ? AND DATE(pci.class_date) = ?
		ORDER BY pci.class_begin_time ASC
	`, personID, personID, today).Rows()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "data": nil})
		return
	}
	defer rows.Close()

	// 解析结果
	for rows.Next() {
		var (
			itemID          int64
			courseID        int64
			courseName      string
			courseDesc      string
			courseRequire   string
			courseClass     string
			classDate       string
			classBeginTime  string
			classEndTime    string
			location        string
			planID          int64
			planName        string
			teacherID       int64
			teacherName     string
			hasEvaluated    bool
			status          string
		)

		rows.Scan(
			&itemID, &courseID, &courseName, &courseDesc, &courseRequire,
			&courseClass, &classDate, &classBeginTime, &classEndTime,
			&location, &planID, &planName, &teacherID, &teacherName,
			&hasEvaluated, &status,
		)

		course := map[string]interface{}{
			"itemId":         itemID,
			"courseId":       courseID,
			"courseName":     courseName,
			"courseDesc":     courseDesc,
			"courseRequire":  courseRequire,
			"courseClass":    courseClass,
			"classDate":      classDate,
			"classBeginTime": classBeginTime,
			"classEndTime":   classEndTime,
			"location":       location,
			"planId":         planID,
			"planName":       planName,
			"teacherId":      teacherID,
			"teacherName":    teacherName,
			"hasEvaluated":   hasEvaluated,
			"status":         status,
		}
		courses = append(courses, course)
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"date":        today,
			"courseCount": len(courses),
			"courses":     courses,
		},
	})
}
