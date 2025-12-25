package planner

import (
	"net/http"
	"strconv"
	"backend/database"

	"github.com/gin-gonic/gin"
)

// GetCourseEvaluations 获取课程评价详情（接口5.18）
func GetCourseEvaluations(c *gin.Context) {
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

	// 1. 查询课程统计信息
	type CourseStats struct {
		CourseAvgScore float64 `json:"courseAvgScore"`
		StudentCount   int64   `json:"studentCount"`
		MaxScore       float64 `json:"maxScore"`
		MinScore       float64 `json:"minScore"`
	}
	var stats CourseStats
	database.DB.Raw(`
		SELECT 
			COALESCE(AVG(ae.weighted_score), 0) AS course_avg_score,
			COUNT(DISTINCT ae.person_id) AS student_count,
			COALESCE(MAX(ae.weighted_score), 0) AS max_score,
			COALESCE(MIN(ae.weighted_score), 0) AS min_score
		FROM plan_course_item pci
		LEFT JOIN attendance_evaluation ae ON pci.item_id = ae.item_id
		WHERE pci.course_id = ?
	`, courseId).Scan(&stats)

	// 2. 查询所有评价详情
	type Evaluation struct {
		PersonID       int64   `json:"personId"`
		PersonName     string  `json:"personName"`
		ItemID         int64   `json:"itemId"`
		ClassDate      string  `json:"classDate"`
		SelfScore      float64 `json:"selfScore"`
		SelfComment    string  `json:"selfComment"`
		TeacherScore   float64 `json:"teacherScore"`
		TeacherComment string  `json:"teacherComment"`
		WeightedScore  float64 `json:"weightedScore"`
	}
	var evaluations []Evaluation
	database.DB.Raw(`
		SELECT 
			p.person_id AS person_id,
			p.name AS person_name,
			ae.item_id AS item_id,
			pci.class_date AS class_date,
			COALESCE(ae.self_score, 0) AS self_score,
			COALESCE(ae.self_comment, '') AS self_comment,
			COALESCE(ae.teacher_score, 0) AS teacher_score,
			COALESCE(ae.teacher_comment, '') AS teacher_comment,
			COALESCE(ae.weighted_score, 0) AS weighted_score
		FROM plan_course_item pci
		INNER JOIN attendance_evaluation ae ON pci.item_id = ae.item_id
		INNER JOIN person p ON ae.person_id = p.person_id
		WHERE pci.course_id = ?
		ORDER BY pci.class_date DESC, ae.weighted_score DESC
	`, courseId).Scan(&evaluations)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"courseId":       courseId,
			"courseName":     course.CourseName,
			"courseAvgScore": stats.CourseAvgScore,
			"studentCount":   stats.StudentCount,
			"maxScore":       stats.MaxScore,
			"minScore":       stats.MinScore,
			"evaluations":    evaluations,
		},
	})
}
