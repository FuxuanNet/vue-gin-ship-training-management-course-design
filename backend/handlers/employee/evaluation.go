package employee

import (
	"net/http"
	"time"

	"backend/database"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

// PendingEvaluationResponse 待自评课程响应
type PendingEvaluationResponse struct {
	ItemID         int    `json:"itemId"`
	CourseID       int    `json:"courseId"`
	CourseName     string `json:"courseName"`
	CourseDesc     string `json:"courseDesc"`
	CourseClass    string `json:"courseClass"`
	ClassDate      string `json:"classDate"`
	ClassBeginTime string `json:"classBeginTime"`
	ClassEndTime   string `json:"classEndTime"`
	Location       string `json:"location"`
	PlanID         int    `json:"planId"`
	PlanName       string `json:"planName"`
	TeacherID      int    `json:"teacherId"`
	TeacherName    string `json:"teacherName"`
}

// EvaluationRequest 提交自评请求
type EvaluationRequest struct {
	ItemID        int    `json:"itemId" binding:"required"`
	SelfComment   string `json:"selfComment" binding:"required"`
	Understanding int    `json:"understanding" binding:"required,min=1,max=5"`
	Difficulty    int    `json:"difficulty" binding:"required,min=1,max=5"`
	Satisfaction  int    `json:"satisfaction" binding:"required,min=1,max=5"`
}

// GetPendingEvaluations 获取待自评课程列表
func GetPendingEvaluations(c *gin.Context) {
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

	var courses []PendingEvaluationResponse

	// 查询已完成但未自评的课程
	now := time.Now()
	today := now.Format("2006-01-02")
	currentTime := now.Format("15:04:05")

	err := database.DB.Table("plan_course_item pci").
		Select(`
			pci.item_id,
			pci.course_id,
			c.course_name,
			c.course_desc,
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
		Joins("JOIN plan_employee pe ON tp.plan_id = pe.plan_id AND pe.person_id = ?", userID).
		Joins("LEFT JOIN attendance_evaluation ae ON pci.item_id = ae.item_id AND ae.person_id = ?", userID).
		Where("(pci.class_date < ? OR (pci.class_date = ? AND pci.class_end_time < ?))", 
			today, today, currentTime).
		Where("(ae.self_comment IS NULL OR ae.self_comment = '')").
		Order("pci.class_date DESC, pci.class_begin_time DESC").
		Scan(&courses).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询待自评课程失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"count":   len(courses),
			"courses": courses,
		},
	})
}

// SubmitEvaluation 提交自评
func SubmitEvaluation(c *gin.Context) {
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

	var req EvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 检查课程是否存在且是否已完成
	var item database.PlanCourseItem
	if err := database.DB.Where("item_id = ?", req.ItemID).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "课程不存在",
			"data":    nil,
		})
		return
	}

	// 检查课程是否已结束
	now := time.Now()
	classDateTime := time.Date(
		item.ClassDate.Year(),
		item.ClassDate.Month(),
		item.ClassDate.Day(),
		item.ClassEndTime.Hour(),
		item.ClassEndTime.Minute(),
		item.ClassEndTime.Second(),
		0,
		time.Local,
	)
	
	if now.Before(classDateTime) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "课程尚未结束，无法提交自评",
			"data":    nil,
		})
		return
	}

	// 检查是否已经自评过
	var existingEval database.AttendanceEvaluation
	err := database.DB.Where("item_id = ? AND person_id = ?", req.ItemID, userID).
		First(&existingEval).Error
	
	if err == nil && existingEval.SelfComment != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "已提交过自评，无需重复提交",
			"data":    nil,
		})
		return
	}

	// 获取课程名称（用于AI评分）
	var course database.Course
	if err := database.DB.Where("course_id = ?", item.CourseID).First(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取课程信息失败",
			"data":    nil,
		})
		return
	}

	// 调用AI生成分数
	score, aiErr := utils.GenerateEvaluationScore(
		req.SelfComment,
		req.Understanding,
		req.Difficulty,
		req.Satisfaction,
		course.CourseName,
	)
	
	if aiErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "AI评分失败: " + aiErr.Error(),
			"data":    nil,
		})
		return
	}

	// 创建或更新自评记录
	evaluation := database.AttendanceEvaluation{
		ItemID:      int64(req.ItemID),
		PersonID:    userID,
		SelfComment: req.SelfComment,
		SelfScore:   score,
	}

	// 使用事务确保数据一致性
	tx := database.DB.Begin()
	
	if err == nil {
		// 更新已存在的记录
		if err := tx.Model(&database.AttendanceEvaluation{}).
			Where("person_id = ? AND item_id = ?", userID, req.ItemID).
			Updates(map[string]interface{}{
				"self_comment": req.SelfComment,
				"self_score":   score,
			}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "更新自评失败",
				"data":    nil,
			})
			return
		}
	} else {
		// 创建新记录
		if err := tx.Create(&evaluation).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "提交自评失败",
				"data":    nil,
			})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "提交成功",
		"data": gin.H{
			"personId":     userID,
			"itemId":       req.ItemID,
			"aiScore":      score,
			"submittedAt":  time.Now(),
		},
	})
}
