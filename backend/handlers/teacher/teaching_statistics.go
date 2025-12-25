package teacher

import (
	"net/http"
	"time"

	"backend/database"

	"github.com/gin-gonic/gin"
)

// GetTeachingStatistics 获取讲师授课统计
func GetTeachingStatistics(c *gin.Context) {
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

	// 获取日期范围参数（可选）
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "开始日期格式错误",
				"data":    nil,
			})
			return
		}
	} else {
		// 默认从今年开始
		startDate = time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local)
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "结束日期格式错误",
				"data":    nil,
			})
			return
		}
	} else {
		// 默认到今天
		endDate = time.Now()
	}

	// 获取讲师姓名
	var teacher database.Person
	if err := database.DB.Where("person_id = ?", teacherID).First(&teacher).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取讲师信息失败",
			"data":    nil,
		})
		return
	}

	// 获取该讲师负责的所有课程
	var courses []database.Course
	database.DB.Where("teacher_id = ?", teacherID).Find(&courses)

	if len(courses) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取成功",
			"data": gin.H{
				"teacherId":   teacherID,
				"teacherName": teacher.Name,
				"period": gin.H{
					"startDate": startDate.Format("2006-01-02"),
					"endDate":   endDate.Format("2006-01-02"),
				},
				"statistics": gin.H{
					"courseCount":              0,
					"classCount":               0,
					"studentCount":             0,
					"totalHours":               0,
					"averageTeachingScore":     0,
					"evaluationRate":           0,
					"courseTypeDistribution":   []map[string]interface{}{},
				},
				"recentClasses": []map[string]interface{}{},
			},
		})
		return
	}

	// 提取课程ID
	courseIDs := make([]int64, len(courses))
	for i, course := range courses {
		courseIDs[i] = course.CourseID
	}

	// 获取时间范围内的所有课程安排
	var courseItems []database.PlanCourseItem
	database.DB.Preload("Course").
		Where("course_id IN ? AND class_date BETWEEN ? AND ?", courseIDs, startDate, endDate).
		Order("class_date DESC").
		Find(&courseItems)

	// 统计总授课次数和时长
	classCount := len(courseItems)
	var totalHours float64
	for _, item := range courseItems {
		// 计算课时
		beginTime := item.ClassBeginTime
		endTime := item.ClassEndTime
		duration := endTime.Sub(beginTime).Hours()
		totalHours += duration
	}

	// 获取所有课程安排的ID
	itemIDs := make([]int64, len(courseItems))
	for i, item := range courseItems {
		itemIDs[i] = item.ItemID
	}

	// 查询所有评价记录
	var evaluations []database.AttendanceEvaluation
	if len(itemIDs) > 0 {
		database.DB.Where("item_id IN ?", itemIDs).Find(&evaluations)
	}

	// 统计学员数（去重）
	studentMap := make(map[int64]bool)
	var totalScore float64
	var evaluatedCount int
	for _, eval := range evaluations {
		studentMap[eval.PersonID] = true
		weightedScore := eval.SelfScore*(1-eval.ScoreRatio) + eval.TeacherScore*eval.ScoreRatio
		totalScore += weightedScore
		if eval.TeacherScore > 0 {
			evaluatedCount++
		}
	}

	studentCount := len(studentMap)
	averageTeachingScore := 0.0
	if len(evaluations) > 0 {
		averageTeachingScore = totalScore / float64(len(evaluations))
	}

	evaluationRate := 0.0
	if len(evaluations) > 0 {
		evaluationRate = float64(evaluatedCount) / float64(len(evaluations)) * 100
	}

	// 按课程类型统计
	courseTypeMap := make(map[string]*CourseTypeStat)
	for _, course := range courses {
		if _, exists := courseTypeMap[course.CourseClass]; !exists {
			courseTypeMap[course.CourseClass] = &CourseTypeStat{
				CourseClass: course.CourseClass,
				Courses:     make(map[int64]bool),
				TotalScore:  0,
				Count:       0,
			}
		}
		courseTypeMap[course.CourseClass].Courses[course.CourseID] = true
	}

	// 统计每种类型的平均分
	for _, eval := range evaluations {
		// 找到该评价对应的课程类型
		for _, item := range courseItems {
			if item.ItemID == eval.ItemID {
				if item.Course.CourseClass != "" {
					if stat, exists := courseTypeMap[item.Course.CourseClass]; exists {
						weightedScore := eval.SelfScore*(1-eval.ScoreRatio) + eval.TeacherScore*eval.ScoreRatio
						stat.TotalScore += weightedScore
						stat.Count++
					}
				}
				break
			}
		}
	}

	courseTypeDistribution := make([]map[string]interface{}, 0)
	for _, stat := range courseTypeMap {
		avgScore := 0.0
		if stat.Count > 0 {
			avgScore = stat.TotalScore / float64(stat.Count)
		}
		courseTypeDistribution = append(courseTypeDistribution, map[string]interface{}{
			"courseClass":  stat.CourseClass,
			"courseCount":  len(stat.Courses),
			"averageScore": avgScore,
		})
	}

	// 最近授课（取最近5次）
	recentClasses := make([]map[string]interface{}, 0)
	recentLimit := 5
	if len(courseItems) < recentLimit {
		recentLimit = len(courseItems)
	}

	for i := 0; i < recentLimit; i++ {
		item := courseItems[i]
		
		// 统计该课次的学员数和评分情况
		var itemEvaluations []database.AttendanceEvaluation
		database.DB.Where("item_id = ?", item.ItemID).Find(&itemEvaluations)
		
		studentCount := len(itemEvaluations)
		evaluatedCount := 0
		var totalScore float64
		for _, eval := range itemEvaluations {
			if eval.TeacherScore > 0 {
				evaluatedCount++
			}
			weightedScore := eval.SelfScore*(1-eval.ScoreRatio) + eval.TeacherScore*eval.ScoreRatio
			totalScore += weightedScore
		}
		
		avgScore := 0.0
		if studentCount > 0 {
			avgScore = totalScore / float64(studentCount)
		}

		recentClasses = append(recentClasses, map[string]interface{}{
			"itemId":         item.ItemID,
			"courseId":       item.CourseID,
			"courseName":     item.Course.CourseName,
			"classDate":      item.ClassDate.Format("2006-01-02"),
			"studentCount":   studentCount,
			"evaluatedCount": evaluatedCount,
			"averageScore":   avgScore,
		})
	}

	// 返回统计数据
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"teacherId":   teacherID,
			"teacherName": teacher.Name,
			"period": gin.H{
				"startDate": startDate.Format("2006-01-02"),
				"endDate":   endDate.Format("2006-01-02"),
			},
			"statistics": gin.H{
				"courseCount":              len(courses),
				"classCount":               classCount,
				"studentCount":             studentCount,
				"totalHours":               totalHours,
				"averageTeachingScore":     averageTeachingScore,
				"evaluationRate":           evaluationRate,
				"courseTypeDistribution":   courseTypeDistribution,
			},
			"recentClasses": recentClasses,
		},
	})
}

// CourseTypeStat 课程类型统计辅助结构
type CourseTypeStat struct {
	CourseClass string
	Courses     map[int64]bool
	TotalScore  float64
	Count       int
}
