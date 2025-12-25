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
	today := time.Now()
	todayStr := today.Format("2006-01-02")
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 59, 999999999, today.Location())

	// 获取讲师负责的所有课程ID
	var courses []database.Course
	if err := database.DB.Where("teacher_id = ?", teacherID).Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
			"data":    nil,
		})
		return
	}

	if len(courses) == 0 {
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

	courseIDs := make([]int64, len(courses))
	for i, course := range courses {
		courseIDs[i] = course.CourseID
	}

	// 查询今日课程安排（使用与schedule相同的查询方式）
	var courseItems []database.PlanCourseItem
	if err := database.DB.
		Preload("Course").
		Preload("Plan").
		Where("course_id IN ?", courseIDs).
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

	// 统计每个课程的学员数和已评分数
	itemIDs := make([]int64, len(courseItems))
	for i, item := range courseItems {
		itemIDs[i] = item.ItemID
	}

	type CountResult struct {
		ItemID         int64
		StudentCount   int64
		EvaluatedCount int64
	}
	var countResults []CountResult
	if len(itemIDs) > 0 {
		database.DB.Raw(`
			SELECT 
				item_id,
				COUNT(*) as student_count,
				SUM(CASE WHEN teacher_comment IS NOT NULL AND teacher_comment != '' THEN 1 ELSE 0 END) as evaluated_count
			FROM attendance_evaluation
			WHERE item_id IN ?
			GROUP BY item_id
		`, itemIDs).Scan(&countResults)
	}

	countMap := make(map[int64]CountResult)
	for _, result := range countResults {
		countMap[result.ItemID] = result
	}

	// 构建响应数据
	response := make([]TodayCourseResponse, 0, len(courseItems))
	for _, item := range courseItems {
		counts := countMap[item.ItemID]
		response = append(response, TodayCourseResponse{
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
			StudentCount:   int(counts.StudentCount),
			EvaluatedCount: int(counts.EvaluatedCount),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"date":        todayStr,
			"courseCount": len(response),
			"courses":     response,
		},
	})
}
