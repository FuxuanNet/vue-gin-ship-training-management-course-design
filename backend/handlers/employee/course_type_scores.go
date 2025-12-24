package employee

import (
	"backend/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCourseTypeScores 获取课程类型成绩分析
func GetCourseTypeScores(c *gin.Context) {
	// 从中间件获取用户ID
	personID, exists := c.Get("person_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录", "data": nil})
		return
	}

	// 查询用户姓名
	var personName string
	database.DB.Raw("SELECT name FROM person WHERE person_id = ?", personID).Scan(&personName)

	// 查询各课程类型的成绩统计
	rows, err := database.DB.Raw(`
		SELECT 
			c.course_class,
			COUNT(*) AS course_count,
			AVG(CASE 
				WHEN ae.teacher_score IS NOT NULL THEN 
					ae.self_score * (1 - ae.score_ratio) + ae.teacher_score * ae.score_ratio
				ELSE ae.self_score
			END) AS average_score,
			MAX(CASE 
				WHEN ae.teacher_score IS NOT NULL THEN 
					ae.self_score * (1 - ae.score_ratio) + ae.teacher_score * ae.score_ratio
				ELSE ae.self_score
			END) AS max_score,
			MIN(CASE 
				WHEN ae.teacher_score IS NOT NULL THEN 
					ae.self_score * (1 - ae.score_ratio) + ae.teacher_score * ae.score_ratio
				ELSE ae.self_score
			END) AS min_score
		FROM attendance_evaluation ae
		JOIN plan_course_item pci ON ae.item_id = pci.item_id
		JOIN course c ON pci.course_id = c.course_id
		WHERE ae.person_id = ? AND ae.self_score IS NOT NULL
		GROUP BY c.course_class
		ORDER BY average_score DESC
	`, personID).Rows()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "data": nil})
		return
	}
	defer rows.Close()

	// 解析结果
	var courseTypeScores []map[string]interface{}
	var indicators []map[string]interface{}
	var values []float64

	for rows.Next() {
		var (
			courseClass  string
			courseCount  int
			averageScore float64
			maxScore     float64
			minScore     float64
		)

		rows.Scan(&courseClass, &courseCount, &averageScore, &maxScore, &minScore)

		courseTypeScore := map[string]interface{}{
			"courseClass":  courseClass,
			"courseCount":  courseCount,
			"averageScore": averageScore,
			"maxScore":     maxScore,
			"minScore":     minScore,
		}
		courseTypeScores = append(courseTypeScores, courseTypeScore)

		// 为雷达图准备数据
		indicators = append(indicators, map[string]interface{}{
			"name": courseClass,
			"max":  100,
		})
		values = append(values, averageScore)
	}

	// 如果没有数据
	if len(courseTypeScores) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "暂无成绩数据",
			"data": gin.H{
				"personId":         personID,
				"personName":       personName,
				"courseTypeScores": []interface{}{},
				"radarData": gin.H{
					"indicators": []interface{}{},
					"values":     []interface{}{},
				},
			},
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"personId":         personID,
			"personName":       personName,
			"courseTypeScores": courseTypeScores,
			"radarData": gin.H{
				"indicators": indicators,
				"values":     values,
			},
		},
	})
}
