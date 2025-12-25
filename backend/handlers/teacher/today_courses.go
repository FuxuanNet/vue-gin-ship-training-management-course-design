package teacher

import (
	"net/http"
	"time"

	"backend/database"

	"github.com/gin-gonic/gin"
)

// TodayCourseResponse 今日授课响应结构
type TodayCourseResponse struct {
	ItemID          int    `json:"itemId"`
	CourseID        int    `json:"courseId"`
	CourseName      string `json:"courseName"`
	CourseDesc      string `json:"courseDesc"`
	CourseRequire   string `json:"courseRequire"`
	CourseClass     string `json:"courseClass"`
	ClassDate       string `json:"classDate"`
	ClassBeginTime  string `json:"classBeginTime"`
	ClassEndTime    string `json:"classEndTime"`
	Location        string `json:"location"`
	PlanID          int    `json:"planId"`
	PlanName        string `json:"planName"`
	StudentCount    int    `json:"studentCount"`
	EvaluatedCount  int    `json:"evaluatedCount"`
}

// GetTodayCourses 获取讲师今日授课列表
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
	
	teacherID := personID.(int64)

	// 获取今天的日期
	today := time.Now().Format("2006-01-02")

	var courses []TodayCourseResponse

	// 查询今日授课课程
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
			tp.plan_name
		`).
		Joins("JOIN course c ON pci.course_id = c.course_id").
		Joins("JOIN training_plan tp ON pci.plan_id = tp.plan_id").
		Where("c.teacher_id = ? AND pci.class_date = ?", teacherID, today).
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

	// 统计每个课程的学员数和已评分数
	for i := range courses {
		// 学员总数（该课程安排的参与人数）
		var studentCount int64
		database.DB.Model(&database.AttendanceEvaluation{}).
			Where("item_id = ?", courses[i].ItemID).
			Count(&studentCount)
		courses[i].StudentCount = int(studentCount)

		// 已评分学员数（teacher_comment不为空）
		var evaluatedCount int64
		database.DB.Model(&database.AttendanceEvaluation{}).
			Where("item_id = ? AND teacher_comment IS NOT NULL AND teacher_comment != ''", 
				courses[i].ItemID).
			Count(&evaluatedCount)
		courses[i].EvaluatedCount = int(evaluatedCount)
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
