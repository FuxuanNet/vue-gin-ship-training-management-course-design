package employee

import (
	"backend/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetScores 获取员工成绩列表
func GetScores(c *gin.Context) {
	// 从中间件获取用户ID
	personID, exists := c.Get("person_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录", "data": nil})
		return
	}

	// 获取查询参数
	planID := c.Query("planId")
	courseClass := c.Query("courseClass")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// 构建WHERE条件
	whereConditions := "WHERE pe.person_id = ? AND pci.class_date < CURDATE()"
	params := []interface{}{personID}

	if planID != "" {
		whereConditions += " AND pci.plan_id = ?"
		params = append(params, planID)
	}
	if courseClass != "" {
		whereConditions += " AND c.course_class = ?"
		params = append(params, courseClass)
	}
	if startDate != "" {
		whereConditions += " AND pci.class_date >= ?"
		params = append(params, startDate)
	}
	if endDate != "" {
		whereConditions += " AND pci.class_date <= ?"
		params = append(params, endDate)
	}

	// 查询成绩数据
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
			ae.self_score,
			ae.self_comment,
			ae.teacher_score,
			ae.teacher_comment,
			ae.score_ratio,
			CASE 
				WHEN ae.teacher_score IS NOT NULL THEN 
					ae.self_score * (1 - ae.score_ratio) + ae.teacher_score * ae.score_ratio
				ELSE NULL
			END AS weighted_score,
			CASE 
				WHEN ae.teacher_score IS NOT NULL THEN true
				ELSE false
			END AS has_teacher_score
		FROM plan_employee pe
		JOIN plan_course_item pci ON pe.plan_id = pci.plan_id
		JOIN course c ON pci.course_id = c.course_id
		JOIN training_plan tp ON pci.plan_id = tp.plan_id
		JOIN person p ON c.teacher_id = p.person_id
		LEFT JOIN attendance_evaluation ae ON ae.item_id = pci.item_id AND ae.person_id = ?
	` + whereConditions + `
		ORDER BY pci.class_date DESC
	`

	rows, err := database.DB.Raw(query, append([]interface{}{personID}, params...)...).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "data": nil})
		return
	}
	defer rows.Close()

	// 解析结果
	var scores []map[string]interface{}
	var totalCourses, completedCourses, pendingEvaluation int
	var sumScore, maxScore, minScore float64
	maxScore = 0
	minScore = 100

	for rows.Next() {
		var (
			itemID          int64
			courseID        int64
			courseName      string
			courseDesc      string
			courseClass     string
			classDate       string
			classBeginTime  string
			classEndTime    string
			location        string
			planID          int64
			planName        string
			teacherID       int64
			teacherName     string
			selfScore       *float64
			selfComment     *string
			teacherScore    *float64
			teacherComment  *string
			scoreRatio      *float64
			weightedScore   *float64
			hasTeacherScore bool
		)

		rows.Scan(
			&itemID, &courseID, &courseName, &courseDesc, &courseClass,
			&classDate, &classBeginTime, &classEndTime, &location,
			&planID, &planName, &teacherID, &teacherName,
			&selfScore, &selfComment, &teacherScore, &teacherComment,
			&scoreRatio, &weightedScore, &hasTeacherScore,
		)

		totalCourses++
		if selfScore != nil {
			completedCourses++
			if weightedScore != nil {
				sumScore += *weightedScore
				if *weightedScore > maxScore {
					maxScore = *weightedScore
				}
				if *weightedScore < minScore {
					minScore = *weightedScore
				}
			} else if selfScore != nil {
				sumScore += *selfScore
				if *selfScore > maxScore {
					maxScore = *selfScore
				}
				if *selfScore < minScore {
					minScore = *selfScore
				}
			}
		} else {
			pendingEvaluation++
		}

		score := map[string]interface{}{
			"itemId":          itemID,
			"courseId":        courseID,
			"courseName":      courseName,
			"courseDesc":      courseDesc,
			"courseClass":     courseClass,
			"classDate":       classDate,
			"classBeginTime":  classBeginTime,
			"classEndTime":    classEndTime,
			"location":        location,
			"planId":          planID,
			"planName":        planName,
			"teacherId":       teacherID,
			"teacherName":     teacherName,
			"selfScore":       selfScore,
			"selfComment":     selfComment,
			"teacherScore":    teacherScore,
			"teacherComment":  teacherComment,
			"scoreRatio":      scoreRatio,
			"weightedScore":   weightedScore,
			"hasTeacherScore": hasTeacherScore,
		}
		scores = append(scores, score)
	}

	// 计算平均分
	var averageScore float64
	if completedCourses > 0 {
		averageScore = sumScore / float64(completedCourses)
	}

	// 如果没有成绩，重置最大最小值
	if completedCourses == 0 {
		maxScore = 0
		minScore = 0
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"statistics": gin.H{
				"totalCourses":      totalCourses,
				"completedCourses":  completedCourses,
				"pendingEvaluation": pendingEvaluation,
				"averageScore":      averageScore,
				"maxScore":          maxScore,
				"minScore":          minScore,
			},
			"scores": scores,
		},
	})
}
