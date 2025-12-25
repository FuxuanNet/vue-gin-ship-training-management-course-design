package teacher

import (
	"net/http"
	"strconv"
	"backend/database"

	"github.com/gin-gonic/gin"

)

// GetCourseStatistics 获取课程成绩统计
func GetCourseStatistics(c *gin.Context) {
	// 获取讲师ID
	teacherID, exists := c.Get("personId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
			"data":    nil,
		})
		return
	}

	// 获取课程ID参数
	courseIDStr := c.Query("courseId")
	if courseIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少课程ID参数",
			"data":    nil,
		})
		return
	}

	courseID, err := strconv.ParseInt(courseIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "课程ID格式错误",
			"data":    nil,
		})
		return
	}

	// 验证该课程是否为该讲师负责
	var course database.Course
	if err := database.DB.Where("course_id = ? AND teacher_id = ?", courseID, teacherID).First(&course).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "无权限：该课程不是您负责的课程",
			"data":    nil,
		})
		return
	}

	// 获取讲师姓名
	var teacher database.Person
	database.DB.Where("person_id = ?", teacherID).First(&teacher)

	// 获取该课程的所有课程安排
	var courseItems []database.PlanCourseItem
	database.DB.Preload("TrainingPlan").Where("course_id = ?", courseID).Find(&courseItems)

	if len(courseItems) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "该课程暂无课程安排",
			"data":    nil,
		})
		return
	}

	// 提取所有itemId
	itemIDs := make([]int64, len(courseItems))
	for i, item := range courseItems {
		itemIDs[i] = item.ItemID
	}

	// 查询所有评价记录
	var evaluations []database.AttendanceEvaluation
	database.DB.Where("item_id IN ?", itemIDs).Find(&evaluations)

	// 计算整体统计数据
	var totalScore float64
	var maxScore float64 = 0
	var minScore float64 = 100
	var passCount int
	var excellentCount int
	var scoreDistribution = map[string]int{
		"0-59":    0,
		"60-69":   0,
		"70-79":   0,
		"80-89":   0,
		"90-100":  0,
	}

	// 用于去重学员
	studentMap := make(map[int64]bool)

	// 用于按学员统计
	studentScoresMap := make(map[int64][]float64)
	studentNamesMap := make(map[int64]string)

	for _, eval := range evaluations {
		// 计算加权分数
		weightedScore := eval.SelfScore*(1-eval.ScoreRatio) + eval.TeacherScore*eval.ScoreRatio
		
		totalScore += weightedScore
		if weightedScore > maxScore {
			maxScore = weightedScore
		}
		if weightedScore < minScore {
			minScore = weightedScore
		}

		// 及格率和优秀率
		if weightedScore >= 60 {
			passCount++
		}
		if weightedScore >= 85 {
			excellentCount++
		}

		// 分数分布
		switch {
		case weightedScore < 60:
			scoreDistribution["0-59"]++
		case weightedScore < 70:
			scoreDistribution["60-69"]++
		case weightedScore < 80:
			scoreDistribution["70-79"]++
		case weightedScore < 90:
			scoreDistribution["80-89"]++
		default:
			scoreDistribution["90-100"]++
		}

		// 学员统计
		studentMap[eval.PersonID] = true
		studentScoresMap[eval.PersonID] = append(studentScoresMap[eval.PersonID], weightedScore)
		
		// 获取学员姓名
		if _, exists := studentNamesMap[eval.PersonID]; !exists {
			var person database.Person
			if err := database.DB.Where("person_id = ?", eval.PersonID).First(&person).Error; err == nil {
				studentNamesMap[eval.PersonID] = person.Name
			}
		}
	}

	totalStudents := len(studentMap)
	averageScore := 0.0
	if len(evaluations) > 0 {
		averageScore = totalScore / float64(len(evaluations))
	}
	passRate := 0.0
	if len(evaluations) > 0 {
		passRate = float64(passCount) / float64(len(evaluations)) * 100
	}
	excellentRate := 0.0
	if len(evaluations) > 0 {
		excellentRate = float64(excellentCount) / float64(len(evaluations)) * 100
	}

	// 按课次统计
	classStatMap := make(map[int64]*ClassStat)
	for _, item := range courseItems {
		classStatMap[item.ItemID] = &ClassStat{
			ItemID:         item.ItemID,
			ClassDate:      item.ClassDate.Format("2006-01-02"),
			Location:       item.Location,
			StudentCount:   0,
			EvaluatedCount: 0,
			TotalScore:     0,
		}
	}

	for _, eval := range evaluations {
		if stat, exists := classStatMap[eval.ItemID]; exists {
			stat.StudentCount++
			if eval.TeacherComment != "" {
				stat.EvaluatedCount++
			}
			weightedScore := eval.SelfScore*(1-eval.ScoreRatio) + eval.TeacherScore*eval.ScoreRatio
			stat.TotalScore += weightedScore
		}
	}

	classStat := make([]map[string]interface{}, 0)
	for _, stat := range classStatMap {
		avgScore := 0.0
		if stat.StudentCount > 0 {
			avgScore = stat.TotalScore / float64(stat.StudentCount)
		}
		classStat = append(classStat, map[string]interface{}{
			"itemId":         stat.ItemID,
			"classDate":      stat.ClassDate,
			"location":       stat.Location,
			"studentCount":   stat.StudentCount,
			"evaluatedCount": stat.EvaluatedCount,
			"averageScore":   avgScore,
		})
	}

	// 按学员统计
	studentScores := make([]map[string]interface{}, 0)
	for personID, scores := range studentScoresMap {
		if len(scores) == 0 {
			continue
		}

		// 计算平均分
		var total float64
		for _, score := range scores {
			total += score
		}
		avgScore := total / float64(len(scores))
		latestScore := scores[len(scores)-1]

		// 判断趋势
		trend := "stable"
		if len(scores) >= 2 {
			if latestScore > scores[len(scores)-2] {
				trend = "up"
			} else if latestScore < scores[len(scores)-2] {
				trend = "down"
			}
		}

		studentScores = append(studentScores, map[string]interface{}{
			"personId":     personID,
			"personName":   studentNamesMap[personID],
			"classCount":   len(scores),
			"averageScore": avgScore,
			"latestScore":  latestScore,
			"trend":        trend,
		})
	}

	// 返回统计数据
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"courseId":    courseID,
			"courseName":  course.CourseName,
			"courseClass": course.CourseClass,
			"teacherName": teacher.Name,
			"statistics": gin.H{
				"totalClasses":      len(courseItems),
				"totalStudents":     totalStudents,
				"averageScore":      averageScore,
				"maxScore":          maxScore,
				"minScore":          minScore,
				"passRate":          passRate,
				"excellentRate":     excellentRate,
				"scoreDistribution": scoreDistribution,
			},
			"classStat":     classStat,
			"studentScores": studentScores,
		},
	})
}

// ClassStat 课次统计辅助结构
type ClassStat struct {
	ItemID         int64
	ClassDate      string
	Location       string
	StudentCount   int
	EvaluatedCount int
	TotalScore     float64
}
