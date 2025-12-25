package employee

import (
	"net/http"
	"time"
	"backend/database"

	"github.com/gin-gonic/gin"

)

// TodayCourseResponse 今日课程响应结构
type TodayCourseResponse struct {
	ItemID          int       `json:"itemId"`
	CourseID        int       `json:"courseId"`
	CourseName      string    `json:"courseName"`
	CourseDesc      string    `json:"courseDesc"`
	CourseRequire   string    `json:"courseRequire"`
	CourseClass     string    `json:"courseClass"`
	ClassDate       string    `json:"classDate"`
	ClassBeginTime  string    `json:"classBeginTime"`
	ClassEndTime    string    `json:"classEndTime"`
	Location        string    `json:"location"`
	PlanID          int       `json:"planId"`
	PlanName        string    `json:"planName"`
	TeacherID       int       `json:"teacherId"`
	TeacherName     string    `json:"teacherName"`
	HasEvaluated    bool      `json:"hasEvaluated"`
	Status          string    `json:"status"`
}

// GetTodayCourses 获取员工今日课程列表
func GetTodayCourses(c *gin.Context) {
	// 获取当前用户信息
	personID, exists := c.Get("personId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
			"data":    nil,
		})
		return
	}
	
	userID := personID.(int64)

	// 获取今天的日期
	today := time.Now()
	todayStr := today.Format("2006-01-02")
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 59, 999999999, today.Location())

	// 1. 查询员工参与的培训计划ID
	var planEmployees []database.PlanEmployee
	if err := database.DB.Where("person_id = ?", userID).Find(&planEmployees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
			"data":    nil,
		})
		return
	}

	planIDs := make([]int64, 0, len(planEmployees))
	for _, pe := range planEmployees {
		planIDs = append(planIDs, pe.PlanID)
	}

	if len(planIDs) == 0 {
		// 没有参与任何培训计划
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取成功",
			"data": gin.H{
				"date":        todayStr,
				"courseCount": 0,
				"courses":     []interface{}{},
			},
		})
		return
	}

	// 2. 查询今日课程安排（使用与schedule相同的查询方式）
	var courseItems []database.PlanCourseItem
	if err := database.DB.
		Preload("Course.Teacher").
		Preload("Plan").
		Where("plan_id IN ?", planIDs).
		Where("class_date >= ? AND class_date <= ?", startOfDay, endOfDay).
		Order("class_begin_time ASC").
		Find(&courseItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询课程失败",
			"data":    nil,
		})
		return
	}

	// 3. 查询员工的评价记录
	itemIDs := make([]int64, 0, len(courseItems))
	for _, item := range courseItems {
		itemIDs = append(itemIDs, item.ItemID)
	}

	evaluatedMap := make(map[int64]bool)
	if len(itemIDs) > 0 {
		var evaluations []database.AttendanceEvaluation
		database.DB.Where("person_id = ? AND item_id IN ? AND self_comment IS NOT NULL AND self_comment != ''", 
			userID, itemIDs).Find(&evaluations)
		for _, eval := range evaluations {
			evaluatedMap[eval.ItemID] = true
		}
	}

	// 4. 构建响应数据
	courses := make([]TodayCourseResponse, 0, len(courseItems))
	for _, item := range courseItems {
		// 判断课程状态
		status := "待上课"
		classEndTime, err := time.Parse("15:04:05", item.ClassEndTime)
		if err == nil {
			classDateTime := time.Date(
				item.ClassDate.Year(),
				item.ClassDate.Month(),
				item.ClassDate.Day(),
				classEndTime.Hour(),
				classEndTime.Minute(),
				classEndTime.Second(),
				0,
				today.Location(),
			)
			if today.After(classDateTime) {
				status = "已完成"
			}
		}

		courses = append(courses, TodayCourseResponse{
			ItemID:         int(item.ItemID),
			CourseID:       int(item.CourseID),
			CourseName:     item.Course.CourseName,
			CourseDesc:     item.Course.CourseDesc,
			CourseRequire:  item.Course.CourseRequire,
			CourseClass:    item.Course.CourseClass,
			ClassDate:      item.ClassDate.Format("2006-01-02"),
			ClassBeginTime: item.ClassBeginTime,
			ClassEndTime:   item.ClassEndTime,
			Location:       item.Location,
			PlanID:         int(item.PlanID),
			PlanName:       item.Plan.PlanName,
			TeacherID:      int(item.Course.TeacherID),
			TeacherName:    item.Course.Teacher.Name,
			HasEvaluated:   evaluatedMap[item.ItemID],
			Status:         status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"date":        todayStr,
			"courseCount": len(courses),
			"courses":     courses,
		},
	})
}
