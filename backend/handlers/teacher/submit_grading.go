package teacher

import (
	"backend/database"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SubmitGradingRequest 提交评分请求
type SubmitGradingRequest struct {
	ItemID         int64    `json:"itemId" binding:"required"`
	PersonID       int64    `json:"personId" binding:"required"`
	TeacherScore   *float64 `json:"teacherScore" binding:"required,min=0,max=100"`
	TeacherComment string   `json:"teacherComment"`
	ScoreRatio     float64  `json:"scoreRatio" binding:"required,min=0,max=1"`
}

// SubmitGradingResponse 提交评分响应
type SubmitGradingResponse struct {
	FinalScore float64 `json:"finalScore"` // 最终综合得分
}

// SubmitGrading 提交讲师评分
func SubmitGrading(c *gin.Context) {
	var req SubmitGradingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 获取当前讲师ID
	teacherID, exists := c.Get("personId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
		})
		return
	}

	// 验证课程是否是该讲师的课程
	var courseItem database.PlanCourseItem
	if err := database.DB.Preload("Course").Where("item_id = ?", req.ItemID).First(&courseItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "课程安排不存在",
		})
		return
	}
	
	// 验证讲师权限
	if courseItem.Course.TeacherID != teacherID.(int64) {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "无权评分该课程",
		})
		return
	}

	// 查询该学员的评价记录
	var evaluation database.AttendanceEvaluation
	err := database.DB.Where("item_id = ? AND person_id = ?", req.ItemID, req.PersonID).First(&evaluation).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "未找到学员评价记录",
		})
		return
	}

	// 如果讲师提供了评语但没有评分，使用AI生成评分
	teacherScore := *req.TeacherScore
	if req.TeacherComment != "" && req.TeacherScore != nil && *req.TeacherScore == 0 {
		// 获取课程名称
		var course database.Course
		database.DB.Where("course_id = ?", courseItem.CourseID).First(&course)
		
		// 调用AI生成讲师评分
		aiScore, err := utils.GenerateTeacherScore(req.TeacherComment, course.CourseName)
		if err == nil && aiScore > 0 {
			teacherScore = aiScore
		}
	}

	// 计算综合得分
	selfScore := evaluation.SelfScore
	finalScore := selfScore*(1-req.ScoreRatio) + teacherScore*req.ScoreRatio

	// 更新评价记录
	evaluation.TeacherScore = teacherScore
	evaluation.TeacherComment = req.TeacherComment
	evaluation.ScoreRatio = req.ScoreRatio

	if err := database.DB.Save(&evaluation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "保存评分失败",
			"error":   err.Error(),
		})
		return
	}

	// 获取学员姓名
	var person database.Person
	database.DB.Where("person_id = ?", req.PersonID).First(&person)
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "评分提交成功",
		"data": gin.H{
			"itemId":         req.ItemID,
			"personId":       req.PersonID,
			"personName":     person.Name,
			"selfScore":      selfScore,
			"teacherScore":   teacherScore,
			"scoreRatio":     req.ScoreRatio,
			"weightedScore":  finalScore,
			"teacherComment": req.TeacherComment,
		},
	})
}
