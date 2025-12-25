package planner

import (
	"net/http"
	"strconv"
	"backend/database"

	"github.com/gin-gonic/gin"
)

// GetAnalytics 获取平台数据分析（接口5.16）
func GetAnalytics(c *gin.Context) {
	// 获取查询参数
	topNStr := c.DefaultQuery("topN", "10")
	topN, err := strconv.Atoi(topNStr)
	if err != nil || topN < 1 {
		topN = 10
	}

	// 1. 课程排名（按平均分）
	type CourseRanking struct {
		CourseID       int64   `json:"courseId"`
		CourseName     string  `json:"courseName"`
		CourseAvgScore float64 `json:"courseAvgScore"`
		StudentCount   int64   `json:"studentCount"`
	}
	var courseRankings []CourseRanking
	database.DB.Raw(`
		SELECT 
			c.course_id AS course_id,
			c.course_name AS course_name,
			COALESCE(AVG(ae.self_score * (1 - ae.score_ratio) + ae.teacher_score * ae.score_ratio), 0) AS course_avg_score,
			COUNT(DISTINCT ae.person_id) AS student_count
		FROM course c
		LEFT JOIN plan_course_item pci ON c.course_id = pci.course_id
		LEFT JOIN attendance_evaluation ae ON pci.item_id = ae.item_id
		WHERE ae.teacher_score != 0 OR ae.teacher_comment != ''
		GROUP BY c.course_id, c.course_name
		HAVING COUNT(DISTINCT ae.person_id) > 0
		ORDER BY course_avg_score DESC
		LIMIT ?
	`, topN).Scan(&courseRankings)

	// 2. 培训计划排名（按平均分）
	type PlanRanking struct {
		PlanID       int64   `json:"planId"`
		PlanName     string  `json:"planName"`
		PlanAvgScore float64 `json:"planAvgScore"`
	}
	var planRankings []PlanRanking
	database.DB.Raw(`
		SELECT 
			tp.plan_id AS plan_id,
			tp.plan_name AS plan_name,
			COALESCE(AVG(ae.self_score * (1 - ae.score_ratio) + ae.teacher_score * ae.score_ratio), 0) AS plan_avg_score
		FROM training_plan tp
		LEFT JOIN plan_course_item pci ON tp.plan_id = pci.plan_id
		LEFT JOIN attendance_evaluation ae ON pci.item_id = ae.item_id
		WHERE ae.teacher_score != 0 OR ae.teacher_comment != ''
		GROUP BY tp.plan_id, tp.plan_name
		HAVING COUNT(DISTINCT ae.person_id) > 0
		ORDER BY plan_avg_score DESC
		LIMIT ?
	`, topN).Scan(&planRankings)

	// 3. 课程类型分布
	type CourseClassDist struct {
		CourseClass  string  `json:"courseClass"`
		CourseCount  int64   `json:"courseCount"`
		AvgScore     float64 `json:"avgScore"`
	}
	var courseClassDistribution []CourseClassDist
	database.DB.Raw(`
		SELECT 
			c.course_class AS course_class,
			COUNT(DISTINCT c.course_id) AS course_count,
			COALESCE(AVG(ae.self_score * (1 - ae.score_ratio) + ae.teacher_score * ae.score_ratio), 0) AS avg_score
		FROM course c
		LEFT JOIN plan_course_item pci ON c.course_id = pci.course_id
		LEFT JOIN attendance_evaluation ae ON pci.item_id = ae.item_id AND (ae.teacher_score != 0 OR ae.teacher_comment != '')
		GROUP BY c.course_class
		ORDER BY course_count DESC
	`).Scan(&courseClassDistribution)

	// 4. 计划状态统计
	type PlanStatusStat struct {
		PlanStatus string `json:"planStatus"`
		Count      int64  `json:"count"`
	}
	var planStatusStatistics []PlanStatusStat
	database.DB.Raw(`
		SELECT plan_status AS plan_status, COUNT(*) AS count
		FROM training_plan
		GROUP BY plan_status
	`).Scan(&planStatusStatistics)

	// 5. 员工排名（按平均分）
	type EmployeeRanking struct {
		PersonID    int64   `json:"personId"`
		PersonName  string  `json:"personName"`
		AvgScore    float64 `json:"avgScore"`
		CourseCount int64   `json:"courseCount"`
	}
	var employeeRankings []EmployeeRanking
	database.DB.Raw(`
		SELECT 
			p.person_id AS person_id,
			p.name AS person_name,
			COALESCE(AVG(ae.self_score * (1 - ae.score_ratio) + ae.teacher_score * ae.score_ratio), 0) AS avg_score,
			COUNT(DISTINCT ae.item_id) AS course_count
		FROM person p
		INNER JOIN attendance_evaluation ae ON p.person_id = ae.person_id
		WHERE p.role = '员工' AND (ae.teacher_score != 0 OR ae.teacher_comment != '')
		GROUP BY p.person_id, p.name
		ORDER BY avg_score DESC
		LIMIT ?
	`, topN).Scan(&employeeRankings)

	// 6. 讲师统计
	type TeacherStat struct {
		TeacherID    int64   `json:"teacherId"`
		TeacherName  string  `json:"teacherName"`
		CourseCount  int64   `json:"courseCount"`
		AvgScore     float64 `json:"avgScore"`
		StudentCount int64   `json:"studentCount"`
	}
	var teacherStatistics []TeacherStat
	database.DB.Raw(`
		SELECT 
			p.person_id AS teacher_id,
			p.name AS teacher_name,
			COUNT(DISTINCT pci.item_id) AS course_count,
			COALESCE(AVG(ae.self_score * (1 - ae.score_ratio) + ae.teacher_score * ae.score_ratio), 0) AS avg_score,
			COUNT(DISTINCT ae.person_id) AS student_count
		FROM person p
		INNER JOIN course c ON p.person_id = c.teacher_id
		LEFT JOIN plan_course_item pci ON c.course_id = pci.course_id
		LEFT JOIN attendance_evaluation ae ON pci.item_id = ae.item_id AND (ae.teacher_score != 0 OR ae.teacher_comment != '')
		WHERE p.role = '讲师'
		GROUP BY p.person_id, p.name
		ORDER BY course_count DESC
	`).Scan(&teacherStatistics)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"courseRankings":          courseRankings,
			"planRankings":            planRankings,
			"courseClassDistribution": courseClassDistribution,
			"planStatusStatistics":    planStatusStatistics,
			"employeeRankings":        employeeRankings,
			"teacherStatistics":       teacherStatistics,
		},
	})
}
