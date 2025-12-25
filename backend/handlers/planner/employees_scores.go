package planner

import (
	"net/http"
	"strconv"
	"backend/database"

	"github.com/gin-gonic/gin"
)

// GetEmployeeScores 获取员工成绩详情（接口5.17）
func GetEmployeeScores(c *gin.Context) {
	// 获取路径参数 employeeId
	employeeIdStr := c.Param("employeeId")
	employeeId, err := strconv.ParseInt(employeeIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的员工ID",
			"data":    nil,
		})
		return
	}

	// 验证员工是否存在且角色正确
	var employee database.Person
	if err := database.DB.Where("person_id = ? AND role = ?", employeeId, "员工").First(&employee).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "员工不存在或角色不正确",
			"data":    nil,
		})
		return
	}

	// 1. 查询整体平均分和课程数
	type OverallStats struct {
		OverallAvgScore float64 `json:"overallAvgScore"`
		CourseCount     int64   `json:"courseCount"`
	}
	var overallStats OverallStats
	database.DB.Raw(`
		SELECT 
			COALESCE(AVG(weighted_score), 0) AS overall_avg_score,
			COUNT(DISTINCT item_id) AS course_count
		FROM attendance_evaluation
		WHERE person_id = ?
	`, employeeId).Scan(&overallStats)

	// 2. 按课程类型的平均分
	type CourseClassScore struct {
		CourseClass       string  `json:"courseClass"`
		AvgWeightedScore  float64 `json:"avgWeightedScore"`
	}
	var courseClassScores []CourseClassScore
	database.DB.Raw(`
		SELECT 
			c.course_class AS course_class,
			COALESCE(AVG(ae.weighted_score), 0) AS avg_weighted_score
		FROM attendance_evaluation ae
		INNER JOIN plan_course_item pci ON ae.item_id = pci.item_id
		INNER JOIN course c ON pci.course_id = c.course_id
		WHERE ae.person_id = ?
		GROUP BY c.course_class
		ORDER BY avg_weighted_score DESC
	`, employeeId).Scan(&courseClassScores)

	// 3. 每节课的详细成绩
	type ItemScore struct {
		ItemID         int64   `json:"itemId"`
		CourseName     string  `json:"courseName"`
		CourseClass    string  `json:"courseClass"`
		ClassDate      string  `json:"classDate"`
		ClassBeginTime string  `json:"classBeginTime"`
		ClassEndTime   string  `json:"classEndTime"`
		SelfScore      float64 `json:"selfScore"`
		TeacherScore   float64 `json:"teacherScore"`
		WeightedScore  float64 `json:"weightedScore"`
	}
	var itemScores []ItemScore
	database.DB.Raw(`
		SELECT 
			ae.item_id AS item_id,
			c.course_name AS course_name,
			c.course_class AS course_class,
			pci.class_date AS class_date,
			pci.class_begin_time AS class_begin_time,
			pci.class_end_time AS class_end_time,
			COALESCE(ae.self_score, 0) AS self_score,
			COALESCE(ae.teacher_score, 0) AS teacher_score,
			COALESCE(ae.weighted_score, 0) AS weighted_score
		FROM attendance_evaluation ae
		INNER JOIN plan_course_item pci ON ae.item_id = pci.item_id
		INNER JOIN course c ON pci.course_id = c.course_id
		WHERE ae.person_id = ?
		ORDER BY pci.class_date DESC, pci.class_begin_time DESC
	`, employeeId).Scan(&itemScores)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"personId":          employeeId,
			"personName":        employee.Name,
			"overallAvgScore":   overallStats.OverallAvgScore,
			"courseCount":       overallStats.CourseCount,
			"courseClassScores": courseClassScores,
			"itemScores":        itemScores,
		},
	})
}
