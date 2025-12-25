package teacher

import (
	"net/http"
	"strconv"
	"time"

	"backend/database"

	"github.com/gin-gonic/gin"
)

// StudentEvaluation 学员评价信息
type StudentEvaluation struct {
	PersonID       int64   `json:"personId"`
	PersonName     string  `json:"personName"`
	SelfScore      float64 `json:"selfScore"`
	SelfComment    string  `json:"selfComment"`
	TeacherScore   *float64 `json:"teacherScore"`
	TeacherComment string  `json:"teacherComment"`
	ScoreRatio     float64 `json:"scoreRatio"`
	EvaluatedAt    *string `json:"evaluatedAt"`
	Status         string  `json:"status"`
}

// CourseItemWithStudents 课程安排及学员信息
type CourseItemWithStudents struct {
	ItemID         int64               `json:"itemId"`
	CourseID       int64               `json:"courseId"`
	CourseName     string              `json:"courseName"`
	CourseClass    string              `json:"courseClass"`
	ClassDate      string              `json:"classDate"`
	ClassBeginTime string              `json:"classBeginTime"`
	ClassEndTime   string              `json:"classEndTime"`
	Location       string              `json:"location"`
	PlanName       string              `json:"planName"`
	Students       []StudentEvaluation `json:"students"`
}

// GetPendingEvaluations 获取待评分学员列表
func GetPendingEvaluations(c *gin.Context) {
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

	// 获取查询参数
	courseIDStr := c.Query("courseId")
	status := c.Query("status")
	if status == "" {
		status = "pending" // 默认只显示待评分
	}

	// 查询讲师已完成的课程安排（课程已结束）
	now := time.Now()
	today := now.Format("2006-01-02")
	currentTime := now.Format("15:04:05")

	query := database.DB.Table("plan_course_item pci").
		Select(`
			pci.item_id,
			pci.course_id,
			c.course_name,
			c.course_class,
			pci.class_date,
			pci.class_begin_time,
			pci.class_end_time,
			pci.location,
			tp.plan_name
		`).
		Joins("JOIN course c ON pci.course_id = c.course_id").
		Joins("JOIN training_plan tp ON pci.plan_id = tp.plan_id").
		Where("c.teacher_id = ?", teacherID).
		Where("(pci.class_date < ? OR (pci.class_date = ? AND pci.class_end_time < ?))", 
			today, today, currentTime)

	// 如果指定了课程ID
	if courseIDStr != "" {
		courseID, err := strconv.ParseInt(courseIDStr, 10, 64)
		if err == nil {
			query = query.Where("c.course_id = ?", courseID)
		}
	}

	var courseItems []struct {
		ItemID         int64
		CourseID       int64
		CourseName     string
		CourseClass    string
		ClassDate      string
		ClassBeginTime string
		ClassEndTime   string
		Location       string
		PlanName       string
	}

	if err := query.Order("pci.class_date DESC, pci.class_begin_time DESC").
		Scan(&courseItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询课程失败",
			"data":    nil,
		})
		return
	}

	// 查询每个课程安排的学员评价信息
	result := []CourseItemWithStudents{}
	totalCount := 0
	pendingCount := 0

	for _, item := range courseItems {
		// 查询该课程的所有学员评价
		var evaluations []database.AttendanceEvaluation
		evalQuery := database.DB.
			Preload("Person").
			Where("item_id = ?", item.ItemID)
		
		// 如果只查待评分的
		if status == "pending" {
			evalQuery = evalQuery.Where("teacher_comment IS NULL OR teacher_comment = ''")
		}
		
		if err := evalQuery.Find(&evaluations).Error; err != nil {
			continue
		}

		// 如果没有学员，跳过
		if len(evaluations) == 0 {
			continue
		}

		students := make([]StudentEvaluation, 0, len(evaluations))
		for _, eval := range evaluations {
			totalCount++
			
			var teacherScore *float64
			if eval.TeacherScore > 0 {
				teacherScore = &eval.TeacherScore
			}
			
			var evaluatedAt *string
			evalStatus := "pending"
			if eval.TeacherComment != "" {
				evalStatus = "evaluated"
				// 这里简化处理，实际可以从updated_at获取
				timeStr := time.Now().Format("2006-01-02 15:04:05")
				evaluatedAt = &timeStr
			} else {
				pendingCount++
			}

			students = append(students, StudentEvaluation{
				PersonID:       eval.PersonID,
				PersonName:     eval.Person.Name,
				SelfScore:      eval.SelfScore,
				SelfComment:    eval.SelfComment,
				TeacherScore:   teacherScore,
				TeacherComment: eval.TeacherComment,
				ScoreRatio:     eval.ScoreRatio,
				EvaluatedAt:    evaluatedAt,
				Status:         evalStatus,
			})
		}

		result = append(result, CourseItemWithStudents{
			ItemID:         item.ItemID,
			CourseID:       item.CourseID,
			CourseName:     item.CourseName,
			CourseClass:    item.CourseClass,
			ClassDate:      item.ClassDate,
			ClassBeginTime: item.ClassBeginTime,
			ClassEndTime:   item.ClassEndTime,
			Location:       item.Location,
			PlanName:       item.PlanName,
			Students:       students,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"totalCount":   totalCount,
			"pendingCount": pendingCount,
			"courseItems":  result,
		},
	})
}
