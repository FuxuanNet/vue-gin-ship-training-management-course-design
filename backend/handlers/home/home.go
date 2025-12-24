package home

import (
	"backend/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetStatistics 获取平台统计数据
func GetStatistics(c *gin.Context) {
	// 获取全局统计数据
	globalStats := getGlobalStatistics()

	// 检查是否已登录
	personID, exists := c.Get("person_id")
	if !exists {
		// 未登录，返回全局统计
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取成功",
			"data":    globalStats,
		})
		return
	}

	// 已登录，查询用户角色
	var person database.Person
	if err := database.DB.First(&person, personID).Error; err != nil {
		// 查询失败，返回全局统计
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取成功",
			"data":    globalStats,
		})
		return
	}

	// 根据角色返回个性化统计
	switch person.Role {
	case "员工":
		personalStats := getEmployeePersonalStats(person.PersonID)
		data := globalStats
		data["personalStats"] = personalStats
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "获取成功", "data": data})
	case "讲师":
		personalStats := getTeacherPersonalStats(person.PersonID)
		data := globalStats
		data["personalStats"] = personalStats
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "获取成功", "data": data})
	default:
		// 课程大纲制定者返回全局统计
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "获取成功", "data": globalStats})
	}
}

// getGlobalStatistics 获取全局统计数据
func getGlobalStatistics() gin.H {
	var courseCount int64
	var teacherCount int64
	var planCount int64
	var totalStudentCount int64
	var totalClassCount int64
	var ongoingPlanCount int64
	var completedPlanCount int64
	var averageSatisfaction float64

	// 培训课程数
	database.DB.Model(&database.Course{}).Count(&courseCount)

	// 讲师团队数
	database.DB.Model(&database.Person{}).Where("role = ?", "讲师").Count(&teacherCount)

	// 培训计划数
	database.DB.Model(&database.TrainingPlan{}).Count(&planCount)

	// 员工总数
	database.DB.Model(&database.Person{}).Where("role = ?", "员工").Count(&totalStudentCount)

	// 总课次（plan_course_item 表总记录数）
	database.DB.Model(&database.PlanCourseItem{}).Count(&totalClassCount)

	// 进行中的计划数
	database.DB.Model(&database.TrainingPlan{}).Where("plan_status = ?", "进行中").Count(&ongoingPlanCount)

	// 已完成的计划数
	database.DB.Model(&database.TrainingPlan{}).Where("plan_status = ?", "已完成").Count(&completedPlanCount)

	// 平均满意度（加权得分的平均值，转换为百分比）
	database.DB.Model(&database.AttendanceEvaluation{}).
		Select("AVG(COALESCE(self_score, 0) * (1 - COALESCE(score_ratio, 0.5)) + COALESCE(teacher_score, 0) * COALESCE(score_ratio, 0.5))").
		Scan(&averageSatisfaction)

	return gin.H{
		"courseCount":         courseCount,
		"teacherCount":        teacherCount,
		"planCount":           planCount,
		"averageSatisfaction": int(averageSatisfaction), // 转为整数百分比
		"totalStudentCount":   totalStudentCount,
		"totalClassCount":     totalClassCount,
		"ongoingPlanCount":    ongoingPlanCount,
		"completedPlanCount":  completedPlanCount,
	}
}

// getEmployeePersonalStats 获取员工个性化统计
func getEmployeePersonalStats(personID int64) gin.H {
	var myPlanCount int64
	var myTotalCourseCount int64
	var myCompletedCourseCount int64
	var myAverageScore float64
	var myTodayCourseCount int64
	var myWeekCourseCount int64

	today := time.Now().Format("2006-01-02")
	weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday())+1).Format("2006-01-02")
	weekEnd := time.Now().AddDate(0, 0, 7-int(time.Now().Weekday())).Format("2006-01-02")

	// 该员工参与的计划数
	database.DB.Model(&database.PlanEmployee{}).Where("person_id = ?", personID).Count(&myPlanCount)

	// 该员工需要上的总课程数
	database.DB.Table("plan_employee").
		Joins("JOIN plan_course_item ON plan_employee.plan_id = plan_course_item.plan_id").
		Where("plan_employee.person_id = ?", personID).
		Count(&myTotalCourseCount)

	// 该员工已完成的课程数（有评分记录的）
	database.DB.Model(&database.AttendanceEvaluation{}).
		Where("person_id = ? AND (self_score IS NOT NULL OR teacher_score IS NOT NULL)", personID).
		Count(&myCompletedCourseCount)

	// 该员工的平均得分
	database.DB.Model(&database.AttendanceEvaluation{}).
		Select("AVG(COALESCE(self_score, 0) * (1 - COALESCE(score_ratio, 0.5)) + COALESCE(teacher_score, 0) * COALESCE(score_ratio, 0.5))").
		Where("person_id = ?", personID).
		Scan(&myAverageScore)

	// 该员工今日课程数
	database.DB.Table("plan_employee").
		Joins("JOIN plan_course_item ON plan_employee.plan_id = plan_course_item.plan_id").
		Where("plan_employee.person_id = ? AND plan_course_item.class_date = ?", personID, today).
		Count(&myTodayCourseCount)

	// 该员工本周课程数
	database.DB.Table("plan_employee").
		Joins("JOIN plan_course_item ON plan_employee.plan_id = plan_course_item.plan_id").
		Where("plan_employee.person_id = ? AND plan_course_item.class_date >= ? AND plan_course_item.class_date <= ?",
			personID, weekStart, weekEnd).
		Count(&myWeekCourseCount)

	return gin.H{
		"myPlanCount":           myPlanCount,
		"myTotalCourseCount":    myTotalCourseCount,
		"myCompletedCourseCount": myCompletedCourseCount,
		"myAverageScore":        myAverageScore,
		"myTodayCourseCount":    myTodayCourseCount,
		"myWeekCourseCount":     myWeekCourseCount,
	}
}

// getTeacherPersonalStats 获取讲师个性化统计
func getTeacherPersonalStats(personID int64) gin.H {
	var myCourseCount int64
	var myClassCount int64
	var myStudentCount int64
	var myAverageTeachingScore float64
	var myTodayClassCount int64
	var myWeekClassCount int64

	today := time.Now().Format("2006-01-02")
	weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday())+1).Format("2006-01-02")
	weekEnd := time.Now().AddDate(0, 0, 7-int(time.Now().Weekday())).Format("2006-01-02")

	// 该讲师主讲的课程数
	database.DB.Model(&database.Course{}).Where("teacher_id = ?", personID).Count(&myCourseCount)

	// 该讲师的授课次数
	database.DB.Table("plan_course_item").
		Joins("JOIN course ON plan_course_item.course_id = course.course_id").
		Where("course.teacher_id = ?", personID).
		Count(&myClassCount)

	// 该讲师教授的学员总数（去重）
	database.DB.Table("attendance_evaluation").
		Select("COUNT(DISTINCT attendance_evaluation.person_id)").
		Joins("JOIN plan_course_item ON attendance_evaluation.item_id = plan_course_item.item_id").
		Joins("JOIN course ON plan_course_item.course_id = course.course_id").
		Where("course.teacher_id = ?", personID).
		Scan(&myStudentCount)

	// 该讲师的平均教学评分
	database.DB.Table("attendance_evaluation").
		Select("AVG(COALESCE(attendance_evaluation.teacher_score, 0))").
		Joins("JOIN plan_course_item ON attendance_evaluation.item_id = plan_course_item.item_id").
		Joins("JOIN course ON plan_course_item.course_id = course.course_id").
		Where("course.teacher_id = ? AND attendance_evaluation.teacher_score IS NOT NULL", personID).
		Scan(&myAverageTeachingScore)

	// 该讲师今日授课数
	database.DB.Table("plan_course_item").
		Joins("JOIN course ON plan_course_item.course_id = course.course_id").
		Where("course.teacher_id = ? AND plan_course_item.class_date = ?", personID, today).
		Count(&myTodayClassCount)

	// 该讲师本周授课数
	database.DB.Table("plan_course_item").
		Joins("JOIN course ON plan_course_item.course_id = course.course_id").
		Where("course.teacher_id = ? AND plan_course_item.class_date >= ? AND plan_course_item.class_date <= ?",
			personID, weekStart, weekEnd).
		Count(&myWeekClassCount)

	return gin.H{
		"myCourseCount":          myCourseCount,
		"myClassCount":           myClassCount,
		"myStudentCount":         myStudentCount,
		"myAverageTeachingScore": myAverageTeachingScore,
		"myTodayClassCount":      myTodayClassCount,
		"myWeekClassCount":       myWeekClassCount,
	}
}
