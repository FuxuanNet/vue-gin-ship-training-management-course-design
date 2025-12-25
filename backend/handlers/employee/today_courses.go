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
	today := time.Now().Format("2006-01-02")

	var courses []TodayCourseResponse

	// 查询今日课程
	err := database.DB.Table("plan_course_item pci").
		Select(`
			pci.item_id,
			pci.course_id,
			c.course_name,
			c.course_desc,
			c.course_require,
			c.course_class,
			pci.class_date,
			pci.class_begin_time,
			pci.class_end_time,
			pci.location,
			pci.plan_id,
			tp.plan_name,
			c.teacher_id,
			p.name as teacher_name
		`).
		Joins("JOIN course c ON pci.course_id = c.course_id").
		Joins("JOIN training_plan tp ON pci.plan_id = tp.plan_id").
		Joins("JOIN person p ON c.teacher_id = p.person_id").
		Joins("JOIN plan_employee pe ON tp.plan_id = pe.plan_id").
		Where("pe.person_id = ? AND pci.class_date = ?", userID, today).
		Order("pci.class_begin_time ASC").
		Scan(&courses).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询课程失败",
			"data":    nil,
		})
		return
	}

	// 检查每个课程是否已自评
	for i := range courses {
		var evalCount int64
		database.DB.Model(&database.AttendanceEvaluation{}).
			Where("item_id = ? AND person_id = ? AND self_comment IS NOT NULL AND self_comment != ''", 
				courses[i].ItemID, userID).
			Count(&evalCount)
		
		courses[i].HasEvaluated = evalCount > 0

		// 判断课程状态
		now := time.Now()
		classDateStr := courses[i].ClassDate
		classEndTimeStr := courses[i].ClassEndTime
		classDateTime, err := time.Parse("2006-01-02 15:04:05", 
			classDateStr+" "+classEndTimeStr)
		if err != nil {
			// 如果解析失败，默认为待上课
			courses[i].Status = "待上课"
			continue
		}
		
		if now.After(classDateTime) {
			courses[i].Status = "已完成"
		} else {
			courses[i].Status = "待上课"
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"date":        today,
			"courseCount": len(courses),
			"courses":     courses,
		},
	})
}
