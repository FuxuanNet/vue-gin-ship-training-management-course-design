package employee

import (
	"net/http"

	"backend/database"

	"github.com/gin-gonic/gin"
)

// GetScores 获取员工成绩列表
func GetScores(c *gin.Context) {
	personID, exists := c.Get("personId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权"})
		return
	}

	// 筛选参数
	// planId := c.Query("planId")
	// courseClass := c.Query("courseClass")
	// startDate := c.Query("startDate")
	// endDate := c.Query("endDate")
	// 暂不实现复杂筛选，先实现基本列表

	var results []struct {
		ItemID         int64   `json:"itemId"`
		ClassDate      string  `json:"classDate"`
		ClassBeginTime string  `json:"classBeginTime"`
		ClassEndTime   string  `json:"classEndTime"`
		Location       string  `json:"location"`
		PlanID         int64   `json:"planId"`
		PlanName       string  `json:"planName"`
		CourseID       int64   `json:"courseId"`
		CourseName     string  `json:"courseName"`
		CourseClass    string  `json:"courseClass"`
		TeacherName    string  `json:"teacherName"`
		SelfScore      float64 `json:"selfScore"`
		TeacherScore   float64 `json:"teacherScore"`
		ScoreRatio     float64 `json:"scoreRatio"`
		SelfComment    string  `json:"selfComment"`
		TeacherComment string  `json:"teacherComment"`
	}

	err := database.DB.Table("attendance_evaluation").
		Select(`
			plan_course_item.item_id,
			DATE_FORMAT(plan_course_item.class_date, '%Y-%m-%d') as class_date,
			TIME_FORMAT(plan_course_item.class_begin_time, '%H:%i:%s') as class_begin_time,
			TIME_FORMAT(plan_course_item.class_end_time, '%H:%i:%s') as class_end_time,
			plan_course_item.location,
			training_plan.plan_id,
			training_plan.plan_name,
			course.course_id,
			course.course_name,
			course.course_class,
			teacher.name as teacher_name,
			attendance_evaluation.self_score,
			attendance_evaluation.teacher_score,
			attendance_evaluation.score_ratio,
			attendance_evaluation.self_comment,
			attendance_evaluation.teacher_comment
		`).
		Joins("JOIN plan_course_item ON attendance_evaluation.item_id = plan_course_item.item_id").
		Joins("JOIN course ON plan_course_item.course_id = course.course_id").
		Joins("JOIN training_plan ON plan_course_item.plan_id = training_plan.plan_id").
		Joins("JOIN person AS teacher ON course.teacher_id = teacher.person_id").
		Where("attendance_evaluation.person_id = ?", personID).
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库查询错误", "error": err.Error()})
		return
	}

	// 处理数据，计算加权分
	var responseData []map[string]interface{}
	var totalScore float64
	var count int
	var maxScore float64 = 0
	var minScore float64 = 100

	for _, r := range results {
		item := map[string]interface{}{
			"itemId":         r.ItemID,
			"classDate":      r.ClassDate,
			"classBeginTime": r.ClassBeginTime,
			"classEndTime":   r.ClassEndTime,
			"location":       r.Location,
			"planId":         r.PlanID,
			"planName":       r.PlanName,
			"courseId":       r.CourseID,
			"courseName":     r.CourseName,
			"courseClass":    r.CourseClass,
			"teacherName":    r.TeacherName,
			"selfScore":      r.SelfScore,
			"teacherScore":   r.TeacherScore,
			"scoreRatio":     r.ScoreRatio,
			"selfComment":    r.SelfComment,
			"teacherComment": r.TeacherComment,
			"hasEvaluated":   true,
		}

		// 计算综合得分（只有当讲师已评分时才计算）
		var weightedScore float64
		hasTeacherGraded := r.TeacherScore > 0 || r.TeacherComment != ""
		item["hasTeacherScore"] = hasTeacherGraded
		
		if hasTeacherGraded {
			// 按权重计算综合得分
			weightedScore = r.SelfScore*(1-r.ScoreRatio) + r.TeacherScore*r.ScoreRatio
			item["weightedScore"] = weightedScore
			totalScore += weightedScore
			count++
			if weightedScore > maxScore {
				maxScore = weightedScore
			}
			if weightedScore < minScore {
				minScore = weightedScore
			}
		} else {
			item["weightedScore"] = nil
		}

		responseData = append(responseData, item)
	}

	avgScore := 0.0
	if count > 0 {
		avgScore = totalScore / float64(count)
	} else {
		minScore = 0 // No scores
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "获取成功",
		"data": gin.H{
			"statistics": gin.H{
				"totalCourses":      len(results), // 简化：只统计了已评价的
				"completedCourses":  count,
				"pendingEvaluation": 0, // 需要另外查询
				"averageScore":      avgScore,
				"maxScore":          maxScore,
				"minScore":          minScore,
			},
			"scores": responseData,
		},
	})
}

// GetCourseTypeScores 获取课程类型成绩分析
func GetCourseTypeScores(c *gin.Context) {
	personID, exists := c.Get("personId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权"})
		return
	}

	type TypeStat struct {
		CourseClass string  `json:"courseClass"`
		AvgScore    float64 `json:"avgWeightedScore"`
		Count       int     `json:"courseCount"`
		MaxScore    float64 `json:"maxScore"`
		MinScore    float64 `json:"minScore"`
	}

	// 这是一个比较复杂的聚合查询，为了简化，我们先查出所有记录在内存中计算
	// 或者使用 SQL Group By
	// 注意：weighted_score 不是数据库字段，是计算出来的。
	// SQL: SUM( (self_score * (1-score_ratio) + teacher_score * score_ratio) ) ...
	// 前提是 teacher_comment 不为空 (表示已评分)

	var stats []TypeStat
	err := database.DB.Table("attendance_evaluation").
		Select(`
			course.course_class,
			AVG(attendance_evaluation.self_score * (1 - attendance_evaluation.score_ratio) + attendance_evaluation.teacher_score * attendance_evaluation.score_ratio) as avg_score,
			COUNT(*) as count,
			MAX(attendance_evaluation.self_score * (1 - attendance_evaluation.score_ratio) + attendance_evaluation.teacher_score * attendance_evaluation.score_ratio) as max_score,
			MIN(attendance_evaluation.self_score * (1 - attendance_evaluation.score_ratio) + attendance_evaluation.teacher_score * attendance_evaluation.score_ratio) as min_score
		`).
		Joins("JOIN plan_course_item ON attendance_evaluation.item_id = plan_course_item.item_id").
		Joins("JOIN course ON plan_course_item.course_id = course.course_id").
		Where("attendance_evaluation.person_id = ?", personID).
		Where("attendance_evaluation.teacher_comment != ''"). // 只统计已完成评分的
		Group("course.course_class").
		Scan(&stats).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库查询错误", "error": err.Error()})
		return
	}

	// 构建雷达图数据
	var indicators []map[string]interface{}
	var values []float64

	for _, s := range stats {
		indicators = append(indicators, map[string]interface{}{
			"name": s.CourseClass,
			"max":  100,
		})
		values = append(values, s.AvgScore)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "获取成功",
		"data": gin.H{
			"personId":         personID,
			"courseTypeScores": stats,
			"radarData": gin.H{
				"indicators": indicators,
				"values":     values,
			},
		},
	})
}

// GetLearningProgress 获取员工学习进度
func GetLearningProgress(c *gin.Context) {
	personID, exists := c.Get("personId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权"})
		return
	}

	// 1. 获取总体进度
	// 需要知道总课程数（从 plan_employee -> training_plan -> plan_course_item）
	// 已完成课程数（attendance_evaluation 中有记录且 teacher_comment != ''）

	var totalCourses int64
	database.DB.Table("plan_course_item").
		Joins("JOIN plan_employee ON plan_course_item.plan_id = plan_employee.plan_id").
		Where("plan_employee.person_id = ?", personID).
		Count(&totalCourses)

	var completedCourses int64
	database.DB.Table("attendance_evaluation").
		Where("person_id = ? AND teacher_comment != ''", personID).
		Count(&completedCourses)

	var evaluatedCourses int64
	database.DB.Table("attendance_evaluation").
		Where("person_id = ?", personID).
		Count(&evaluatedCourses)

	// 2. 获取各计划进度
	type PlanProgress struct {
		PlanID             int64   `json:"planId"`
		PlanName           string  `json:"planName"`
		PlanStatus         string  `json:"planStatus"`
		StartDate          string  `json:"startDate"`
		EndDate            string  `json:"endDate"`
		TotalCourses       int     `json:"totalCourses"`
		CompletedCourses   int     `json:"completedCourses"`
		EvaluatedCourses   int     `json:"evaluatedCourses"`
		ProgressPercentage float64 `json:"progressPercentage"`
		AverageScore       float64 `json:"averageScore"`
	}

	var plans []PlanProgress
	// 这里需要先查出员工参与的所有计划，然后遍历计算（或者复杂SQL）
	// 简单起见，先查计划
	database.DB.Table("training_plan").
		Select("training_plan.plan_id, training_plan.plan_name, training_plan.plan_status, training_plan.plan_start_datetime, training_plan.plan_end_datetime").
		Joins("JOIN plan_employee ON training_plan.plan_id = plan_employee.plan_id").
		Where("plan_employee.person_id = ?", personID).
		Scan(&plans)

	for i := range plans {
		// 查该计划的总课程数
		var pTotal int64
		database.DB.Table("plan_course_item").Where("plan_id = ?", plans[i].PlanID).Count(&pTotal)
		plans[i].TotalCourses = int(pTotal)

		// 查该计划已完成（有讲师评分）
		var pCompleted int64
		database.DB.Table("attendance_evaluation").
			Joins("JOIN plan_course_item ON attendance_evaluation.item_id = plan_course_item.item_id").
			Where("attendance_evaluation.person_id = ? AND plan_course_item.plan_id = ? AND attendance_evaluation.teacher_comment != ''", personID, plans[i].PlanID).
			Count(&pCompleted)
		plans[i].CompletedCourses = int(pCompleted)

		// 查该计划已自评
		var pEvaluated int64
		database.DB.Table("attendance_evaluation").
			Joins("JOIN plan_course_item ON attendance_evaluation.item_id = plan_course_item.item_id").
			Where("attendance_evaluation.person_id = ? AND plan_course_item.plan_id = ?", personID, plans[i].PlanID).
			Count(&pEvaluated)
		plans[i].EvaluatedCourses = int(pEvaluated)

		if pTotal > 0 {
			plans[i].ProgressPercentage = float64(pCompleted) / float64(pTotal) * 100
		}

		// 计算平均分
		var avgScore float64
		database.DB.Table("attendance_evaluation").
			Select("AVG(self_score * (1 - score_ratio) + teacher_score * score_ratio)").
			Joins("JOIN plan_course_item ON attendance_evaluation.item_id = plan_course_item.item_id").
			Where("attendance_evaluation.person_id = ? AND plan_course_item.plan_id = ? AND attendance_evaluation.teacher_comment != ''", personID, plans[i].PlanID).
			Scan(&avgScore)
		plans[i].AverageScore = avgScore
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "获取成功",
		"data": gin.H{
			"personId": personID,
			"overallProgress": gin.H{
				"totalPlans":       len(plans),
				"totalCourses":     totalCourses,
				"completedCourses": completedCourses,
				"evaluatedCourses": evaluatedCourses,
				"progressPercentage": 0, // Calculate if needed
			},
			"planProgress": plans,
			"recentCourses": []interface{}{}, // Placeholder
		},
	})
}
