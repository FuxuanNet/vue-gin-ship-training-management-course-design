package employee

import (
	"backend/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetLearningProgress 获取员工学习进度
func GetLearningProgress(c *gin.Context) {
	// 从中间件获取用户ID
	personID, exists := c.Get("person_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录", "data": nil})
		return
	}

	// 查询用户姓名
	var personName string
	database.DB.Raw("SELECT name FROM person WHERE person_id = ?", personID).Scan(&personName)

	// 查询整体进度
	var (
		totalPlans        int
		totalCourses      int
		completedCourses  int
		evaluatedCourses  int
		averageScore      float64
	)

	// 参与的培训计划总数
	database.DB.Raw(`
		SELECT COUNT(DISTINCT plan_id) FROM plan_employee WHERE person_id = ?
	`, personID).Scan(&totalPlans)

	// 总课程数（已上过的）
	database.DB.Raw(`
		SELECT COUNT(*)
		FROM plan_employee pe
		JOIN plan_course_item pci ON pe.plan_id = pci.plan_id
		WHERE pe.person_id = ? AND pci.class_date < CURDATE()
	`, personID).Scan(&totalCourses)

	// 已完成课程数（有自评或讲师评分的）
	database.DB.Raw(`
		SELECT COUNT(*)
		FROM attendance_evaluation ae
		JOIN plan_course_item pci ON ae.item_id = pci.item_id
		WHERE ae.person_id = ? AND (ae.self_score IS NOT NULL OR ae.teacher_score IS NOT NULL)
	`, personID).Scan(&completedCourses)

	// 已自评课程数
	database.DB.Raw(`
		SELECT COUNT(*)
		FROM attendance_evaluation 
		WHERE person_id = ? AND self_score IS NOT NULL
	`, personID).Scan(&evaluatedCourses)

	// 平均分
	database.DB.Raw(`
		SELECT AVG(CASE 
			WHEN teacher_score IS NOT NULL THEN 
				self_score * (1 - score_ratio) + teacher_score * score_ratio
			ELSE self_score
		END)
		FROM attendance_evaluation 
		WHERE person_id = ? AND self_score IS NOT NULL
	`, personID).Scan(&averageScore)

	// 计算完成百分比
	progressPercentage := 0.0
	if totalCourses > 0 {
		progressPercentage = float64(completedCourses) / float64(totalCourses) * 100
	}

	// 查询各计划的进度
	planRows, err := database.DB.Raw(`
		SELECT 
			tp.plan_id,
			tp.plan_name,
			tp.plan_status,
			COUNT(DISTINCT pci.item_id) AS total_courses,
			COUNT(DISTINCT ae.item_id) AS completed_courses,
			AVG(CASE 
				WHEN ae.teacher_score IS NOT NULL THEN 
					ae.self_score * (1 - ae.score_ratio) + ae.teacher_score * ae.score_ratio
				ELSE ae.self_score
			END) AS average_score
		FROM plan_employee pe
		JOIN training_plan tp ON pe.plan_id = tp.plan_id
		LEFT JOIN plan_course_item pci ON tp.plan_id = pci.plan_id AND pci.class_date < CURDATE()
		LEFT JOIN attendance_evaluation ae ON ae.item_id = pci.item_id AND ae.person_id = ?
		WHERE pe.person_id = ?
		GROUP BY tp.plan_id, tp.plan_name, tp.plan_status
		ORDER BY tp.plan_start_datetime DESC
	`, personID, personID).Rows()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询计划进度失败", "data": nil})
		return
	}
	defer planRows.Close()

	var planProgress []map[string]interface{}
	for planRows.Next() {
		var (
			planID           int64
			planName         string
			planStatus       string
			planTotalCourses int
			planCompleted    int
			planAvgScore     *float64
		)

		planRows.Scan(&planID, &planName, &planStatus, &planTotalCourses, &planCompleted, &planAvgScore)

		planProgressPercentage := 0.0
		if planTotalCourses > 0 {
			planProgressPercentage = float64(planCompleted) / float64(planTotalCourses) * 100
		}

		avgScore := 0.0
		if planAvgScore != nil {
			avgScore = *planAvgScore
		}

		plan := map[string]interface{}{
			"planId":              planID,
			"planName":            planName,
			"planStatus":          planStatus,
			"totalCourses":        planTotalCourses,
			"completedCourses":    planCompleted,
			"progressPercentage":  planProgressPercentage,
			"averageScore":        avgScore,
		}
		planProgress = append(planProgress, plan)
	}

	// 查询最近5条学习记录
	recentRows, err := database.DB.Raw(`
		SELECT 
			pci.item_id,
			pci.course_id,
			c.course_name,
			pci.class_date,
			CASE 
				WHEN ae.person_id IS NOT NULL THEN true
				ELSE false
			END AS has_evaluated
		FROM plan_employee pe
		JOIN plan_course_item pci ON pe.plan_id = pci.plan_id
		JOIN course c ON pci.course_id = c.course_id
		LEFT JOIN attendance_evaluation ae ON ae.item_id = pci.item_id AND ae.person_id = ?
		WHERE pe.person_id = ? AND pci.class_date < CURDATE()
		ORDER BY pci.class_date DESC
		LIMIT 5
	`, personID, personID).Rows()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询最近课程失败", "data": nil})
		return
	}
	defer recentRows.Close()

	var recentCourses []map[string]interface{}
	for recentRows.Next() {
		var (
			itemID       int64
			courseID     int64
			courseName   string
			classDate    string
			hasEvaluated bool
		)

		recentRows.Scan(&itemID, &courseID, &courseName, &classDate, &hasEvaluated)

		recent := map[string]interface{}{
			"itemId":       itemID,
			"courseId":     courseID,
			"courseName":   courseName,
			"classDate":    classDate,
			"hasEvaluated": hasEvaluated,
		}
		recentCourses = append(recentCourses, recent)
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"personId":   personID,
			"personName": personName,
			"overallProgress": gin.H{
				"totalPlans":         totalPlans,
				"totalCourses":       totalCourses,
				"completedCourses":   completedCourses,
				"evaluatedCourses":   evaluatedCourses,
				"progressPercentage": progressPercentage,
				"averageScore":       averageScore,
			},
			"planProgress":   planProgress,
			"recentCourses":  recentCourses,
		},
	})
}
