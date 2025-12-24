package employee

import (
	"backend/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPendingEvaluations 获取待自评课程列表
func GetPendingEvaluations(c *gin.Context) {
	// 从中间件获取用户ID
	personID, exists := c.Get("person_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录", "data": nil})
		return
	}

	// 获取查询参数
	status := c.DefaultQuery("status", "pending") // pending 或 all
	limitStr := c.DefaultQuery("limit", "0")      // 0表示不限制
	limit, _ := strconv.Atoi(limitStr)

	// 构建WHERE条件
	whereClause := ""
	if status == "pending" {
		whereClause = "AND ae.person_id IS NULL" // 未自评
	}

	// 构建LIMIT条件
	limitClause := ""
	if limit > 0 {
		limitClause = " LIMIT " + limitStr
	}

	// 查询课程（已上完的课程）
	query := `
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
				WHEN ae.person_id IS NOT NULL THEN 'evaluated'
				ELSE 'pending'
			END AS evaluation_status
		FROM plan_employee pe
		JOIN plan_course_item pci ON pe.plan_id = pci.plan_id
		JOIN course c ON pci.course_id = c.course_id
		JOIN training_plan tp ON pci.plan_id = tp.plan_id
		JOIN person p ON c.teacher_id = p.person_id
		LEFT JOIN attendance_evaluation ae ON ae.item_id = pci.item_id AND ae.person_id = ?
		WHERE pe.person_id = ? 
			AND pci.class_date < CURDATE()
			` + whereClause + `
		ORDER BY pci.class_date DESC
	` + limitClause

	rows, err := database.DB.Raw(query, personID, personID).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "data": nil})
		return
	}
	defer rows.Close()

	// 解析结果
	var courses []map[string]interface{}
	pendingCount := 0

	for rows.Next() {
		var (
			itemID           int64
			courseID         int64
			courseName       string
			courseDesc       string
			courseClass      string
			classDate        string
			classBeginTime   string
			classEndTime     string
			location         string
			planID           int64
			planName         string
			teacherID        int64
			teacherName      string
			evaluationStatus string
		)

		rows.Scan(
			&itemID, &courseID, &courseName, &courseDesc, &courseClass,
			&classDate, &classBeginTime, &classEndTime, &location,
			&planID, &planName, &teacherID, &teacherName, &evaluationStatus,
		)

		if evaluationStatus == "pending" {
			pendingCount++
		}

		course := map[string]interface{}{
			"itemId":           itemID,
			"courseId":         courseID,
			"courseName":       courseName,
			"courseDesc":       courseDesc,
			"courseClass":      courseClass,
			"classDate":        classDate,
			"classBeginTime":   classBeginTime,
			"classEndTime":     classEndTime,
			"location":         location,
			"planId":           planID,
			"planName":         planName,
			"teacherId":        teacherID,
			"teacherName":      teacherName,
			"evaluationStatus": evaluationStatus,
		}
		courses = append(courses, course)
	}

	// 如果没有设置limit，totalCount等于courses长度
	totalCount := len(courses)
	if limit > 0 && status == "all" {
		// 需要查询总数
		var count int64
		database.DB.Raw(`
			SELECT COUNT(*)
			FROM plan_employee pe
			JOIN plan_course_item pci ON pe.plan_id = pci.plan_id
			WHERE pe.person_id = ? AND pci.class_date < CURDATE()
		`, personID).Scan(&count)
		totalCount = int(count)
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"totalCount":   totalCount,
			"pendingCount": pendingCount,
			"courses":      courses,
		},
	})
}
