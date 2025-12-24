package auth

import (
	"backend/database"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Name     string `json:"name" binding:"required,min=2,max=20"`
	Role     string `json:"role" binding:"required,oneof=employee teacher planner"`
}

// 角色映射：英文 -> 中文
var roleMap = map[string]string{
	"employee": "员工",
	"teacher":  "讲师",
	"planner":  "课程大纲制定者",
}

// 角色映射：中文 -> 英文
var reverseRoleMap = map[string]string{
	"员工":       "employee",
	"讲师":       "teacher",
	"课程大纲制定者": "planner",
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误：" + err.Error(), "data": nil})
		return
	}

	// 1. 查询账号
	var account database.Account
	if err := database.DB.Where("login_name = ?", req.Username).First(&account).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误", "data": nil})
		return
	}

	// 2. 验证密码（bcrypt）
	if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误", "data": nil})
		return
	}

	// 3. 查询用户信息
	var person database.Person
	if err := database.DB.First(&person, account.PersonID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询用户信息失败", "data": nil})
		return
	}

	// 4. 创建会话
	sessionID := generateSessionID()
	session := database.Session{
		SessionID: sessionID,
		PersonID:  person.PersonID,
		Role:      person.Role,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24小时过期
	}
	if err := database.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建会话失败", "data": nil})
		return
	}

	// 5. 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data": gin.H{
			"token": sessionID, // 前端存储为Session-ID
			"user": gin.H{
				"id":          person.PersonID,
				"name":        person.Name,
				"role":        reverseRoleMap[person.Role],
				"roleDisplay": person.Role,
				"accountId":   account.AccountID,
			},
		},
	})
}

// Register 用户注册
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误：" + err.Error(), "data": nil})
		return
	}

	// 1. 检查用户名是否已存在
	var count int64
	database.DB.Model(&database.Account{}).Where("login_name = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户名已被使用", "data": nil})
		return
	}

	// 2. 哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "密码加密失败", "data": nil})
		return
	}

	// 3. 创建人员记录
	person := database.Person{
		Name: req.Name,
		Role: roleMap[req.Role],
	}
	if err := database.DB.Create(&person).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建用户信息失败", "data": nil})
		return
	}

	// 4. 创建账号记录
	account := database.Account{
		PersonID:     person.PersonID,
		LoginName:    req.Username,
		PasswordHash: string(hashedPassword),
	}
	if err := database.DB.Create(&account).Error; err != nil {
		// 回滚：删除已创建的 person
		database.DB.Delete(&person)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建账号失败", "data": nil})
		return
	}

	// 5. 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
		"data": gin.H{
			"personId":    person.PersonID,
			"accountId":   account.AccountID,
			"username":    account.LoginName,
			"name":        person.Name,
			"role":        req.Role,
			"roleDisplay": person.Role,
		},
	})
}

// Logout 退出登录
func Logout(c *gin.Context) {
	sessionID := c.GetHeader("Session-ID")
	if sessionID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录或登录已过期", "data": nil})
		return
	}

	// 删除会话
	database.DB.Where("session_id = ?", sessionID).Delete(&database.Session{})

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "退出成功", "data": nil})
}

// GetCurrentUser 获取当前用户信息
func GetCurrentUser(c *gin.Context) {
	// 从中间件设置的上下文获取用户信息
	personID, exists := c.Get("person_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录或登录已过期", "data": nil})
		return
	}

	// 查询用户信息
	var person database.Person
	if err := database.DB.First(&person, personID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在", "data": nil})
		return
	}

	// 查询账号信息
	var account database.Account
	database.DB.Where("person_id = ?", personID).First(&account)

	// 根据角色获取统计数据
	statistics := getStatisticsByRole(person.PersonID, person.Role)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"personId":    person.PersonID,
			"name":        person.Name,
			"role":        reverseRoleMap[person.Role],
			"roleDisplay": person.Role,
			"accountId":   account.AccountID,
			"username":    account.LoginName,
			"statistics":  statistics,
		},
	})
}

// getStatisticsByRole 根据角色获取统计数据
func getStatisticsByRole(personID int64, role string) gin.H {
	switch role {
	case "员工":
		return getEmployeeStatistics(personID)
	case "讲师":
		return getTeacherStatistics(personID)
	case "课程大纲制定者":
		return getPlannerStatistics(personID)
	default:
		return gin.H{}
	}
}

// getEmployeeStatistics 获取员工统计数据
func getEmployeeStatistics(personID int64) gin.H {
	var trainingPlanCount int64
	var totalCourseCount int64
	var completedCourseCount int64
	var avgScore float64

	// 参与的培训计划数量
	database.DB.Model(&database.PlanEmployee{}).Where("person_id = ?", personID).Count(&trainingPlanCount)

	// 总课程数量（通过plan_employee和plan_course_item关联）
	database.DB.Table("plan_employee").
		Joins("JOIN plan_course_item ON plan_employee.plan_id = plan_course_item.plan_id").
		Where("plan_employee.person_id = ?", personID).
		Count(&totalCourseCount)

	// 已完成的课程数量（有评分记录的）
	database.DB.Model(&database.AttendanceEvaluation{}).
		Where("person_id = ? AND (self_score IS NOT NULL OR teacher_score IS NOT NULL)", personID).
		Count(&completedCourseCount)

	// 平均得分（加权计算）
	database.DB.Model(&database.AttendanceEvaluation{}).
		Select("AVG(COALESCE(self_score, 0) * (1 - COALESCE(score_ratio, 0.5)) + COALESCE(teacher_score, 0) * COALESCE(score_ratio, 0.5))").
		Where("person_id = ?", personID).
		Scan(&avgScore)

	return gin.H{
		"trainingPlanCount":   trainingPlanCount,
		"completedCourseCount": completedCourseCount,
		"totalCourseCount":    totalCourseCount,
		"averageScore":        avgScore,
	}
}

// getTeacherStatistics 获取讲师统计数据
func getTeacherStatistics(personID int64) gin.H {
	var courseCount int64
	var classCount int64
	var studentCount int64
	var avgScore float64

	// 主讲的课程数量
	database.DB.Model(&database.Course{}).Where("teacher_id = ?", personID).Count(&courseCount)

	// 授课次数
	database.DB.Table("plan_course_item").
		Joins("JOIN course ON plan_course_item.course_id = course.course_id").
		Where("course.teacher_id = ?", personID).
		Count(&classCount)

	// 教授的学员总数（去重）
	database.DB.Table("attendance_evaluation").
		Select("COUNT(DISTINCT person_id)").
		Joins("JOIN plan_course_item ON attendance_evaluation.item_id = plan_course_item.item_id").
		Joins("JOIN course ON plan_course_item.course_id = course.course_id").
		Where("course.teacher_id = ?", personID).
		Scan(&studentCount)

	// 平均教学评分（学员给的分数）
	database.DB.Table("attendance_evaluation").
		Select("AVG(COALESCE(teacher_score, 0))").
		Joins("JOIN plan_course_item ON attendance_evaluation.item_id = plan_course_item.item_id").
		Joins("JOIN course ON plan_course_item.course_id = course.course_id").
		Where("course.teacher_id = ? AND teacher_score IS NOT NULL", personID).
		Scan(&avgScore)

	return gin.H{
		"courseCount":           courseCount,
		"classCount":            classCount,
		"studentCount":          studentCount,
		"averageTeachingScore":  avgScore,
	}
}

// getPlannerStatistics 获取课程大纲制定者统计数据
func getPlannerStatistics(personID int64) gin.H {
	var planCount int64
	var totalCourseCount int64
	var totalStudentCount int64
	var avgScore float64

	// 创建的培训计划数量
	database.DB.Model(&database.TrainingPlan{}).Where("creator_id = ?", personID).Count(&planCount)

	// 所有培训计划中的课程总数
	database.DB.Table("plan_course_item").
		Joins("JOIN training_plan ON plan_course_item.plan_id = training_plan.plan_id").
		Where("training_plan.creator_id = ?", personID).
		Count(&totalCourseCount)

	// 所有培训计划中的学员总数
	database.DB.Table("plan_employee").
		Joins("JOIN training_plan ON plan_employee.plan_id = training_plan.plan_id").
		Where("training_plan.creator_id = ?", personID).
		Count(&totalStudentCount)

	// 所有培训计划的平均得分
	database.DB.Table("attendance_evaluation").
		Select("AVG(COALESCE(self_score, 0) * (1 - COALESCE(score_ratio, 0.5)) + COALESCE(teacher_score, 0) * COALESCE(score_ratio, 0.5))").
		Joins("JOIN plan_course_item ON attendance_evaluation.item_id = plan_course_item.item_id").
		Joins("JOIN training_plan ON plan_course_item.plan_id = training_plan.plan_id").
		Where("training_plan.creator_id = ?", personID).
		Scan(&avgScore)

	return gin.H{
		"planCount":          planCount,
		"totalCourseCount":   totalCourseCount,
		"totalStudentCount":  totalStudentCount,
		"averagePlanScore":   avgScore,
	}
}

// generateSessionID 生成随机会话ID
func generateSessionID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
