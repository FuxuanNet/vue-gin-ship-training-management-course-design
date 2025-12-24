package employee

import (
	"backend/database"
	"backend/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SubmitEvaluationRequest 提交自评请求结构
type SubmitEvaluationRequest struct {
	ItemID        int64  `json:"itemId" binding:"required"`
	SelfComment   string `json:"selfComment" binding:"required,min=50,max=1000"`
	Understanding int    `json:"understanding"` // 1-5
	Difficulty    int    `json:"difficulty"`    // 1-5
	Satisfaction  int    `json:"satisfaction"`  // 1-5
}

// SubmitEvaluation 提交课程自评
func SubmitEvaluation(c *gin.Context) {
	// 从中间件获取用户ID
	personID, exists := c.Get("person_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录", "data": nil})
		return
	}

	// 解析请求
	var req SubmitEvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误：" + err.Error(), "data": nil})
		return
	}

	// 设置默认值
	if req.Understanding == 0 {
		req.Understanding = 3
	}
	if req.Difficulty == 0 {
		req.Difficulty = 3
	}
	if req.Satisfaction == 0 {
		req.Satisfaction = 3
	}

	// 验证课程是否属于该员工
	var count int64
	database.DB.Raw(`
		SELECT COUNT(*)
		FROM plan_employee pe
		JOIN plan_course_item pci ON pe.plan_id = pci.plan_id
		WHERE pe.person_id = ? AND pci.item_id = ?
	`, personID, req.ItemID).Scan(&count)

	if count == 0 {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权限：该课程不属于您的课程安排", "data": nil})
		return
	}

	// 查询课程信息
	var (
		courseID   int64
		courseName string
		classDate  string
	)
	err := database.DB.Raw(`
		SELECT pci.course_id, c.course_name, pci.class_date
		FROM plan_course_item pci
		JOIN course c ON pci.course_id = c.course_id
		WHERE pci.item_id = ?
	`, req.ItemID).Row().Scan(&courseID, &courseName, &classDate)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "课程不存在", "data": nil})
		return
	}

	// 验证课程是否已上完
	classDateParsed, _ := time.Parse("2006-01-02", classDate)
	if classDateParsed.After(time.Now()) || classDateParsed.Equal(time.Now().Truncate(24*time.Hour)) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "课程尚未开始，暂不能自评", "data": nil})
		return
	}

	// 调用AI生成评分
	selfScore, err := utils.GenerateEvaluationScore(
		req.SelfComment,
		req.Understanding,
		req.Difficulty,
		req.Satisfaction,
		courseName,
	)
	if err != nil {
		// AI调用失败，记录日志但不中断流程
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "AI评分服务异常，请稍后重试", "data": nil})
		return
	}

	// 查询是否已有评价记录
	var existingID int64
	database.DB.Raw(`
		SELECT person_id 
		FROM attendance_evaluation 
		WHERE person_id = ? AND item_id = ?
	`, personID, req.ItemID).Scan(&existingID)

	// 插入或更新评价记录
	if existingID > 0 {
		// 更新现有记录
		err = database.DB.Exec(`
			UPDATE attendance_evaluation 
			SET self_score = ?, self_comment = ?
			WHERE person_id = ? AND item_id = ?
		`, selfScore, req.SelfComment, personID, req.ItemID).Error
	} else {
		// 插入新记录（score_ratio默认为0.5，即自评和讲师评分各占50%）
		err = database.DB.Exec(`
			INSERT INTO attendance_evaluation 
			(person_id, item_id, self_score, self_comment, score_ratio) 
			VALUES (?, ?, ?, ?, 0.5)
		`, personID, req.ItemID, selfScore, req.SelfComment).Error
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存失败", "data": nil})
		return
	}

	// 查询讲师评分（如果有）
	var teacherScore *float64
	var weightedScore *float64
	database.DB.Raw(`
		SELECT teacher_score,
			CASE 
				WHEN teacher_score IS NOT NULL THEN 
					self_score * (1 - score_ratio) + teacher_score * score_ratio
				ELSE NULL
			END AS weighted_score
		FROM attendance_evaluation 
		WHERE person_id = ? AND item_id = ?
	`, personID, req.ItemID).Row().Scan(&teacherScore, &weightedScore)

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "自评提交成功",
		"data": gin.H{
			"itemId":        req.ItemID,
			"courseId":      courseID,
			"courseName":    courseName,
			"selfScore":     selfScore,
			"selfComment":   req.SelfComment,
			"teacherScore":  teacherScore,
			"weightedScore": weightedScore,
			"evaluatedAt":   time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}
