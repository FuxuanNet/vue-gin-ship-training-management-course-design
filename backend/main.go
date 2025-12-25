package main

import (
	"log"
	"backend/config"
	"backend/database"
	"backend/handlers/auth"
	"backend/handlers/employee"
	"backend/handlers/home"
	"backend/handlers/planner"
	"backend/handlers/teacher"
	"backend/middleware"

	"github.com/gin-gonic/gin"

)

func main() {
	// 1. 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	// 2. 初始化数据库
	if err := database.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer database.CloseDB()

	// 3. 插入测试账号（首次运行时自动插入，已存在则跳过）
	if err := database.SeedTestAccounts(); err != nil {
		log.Printf("测试账号插入失败: %v", err)
	}

	// 4. 创建 Gin 引擎
	r := gin.Default()

	// 5. 应用全局中间件
	r.Use(middleware.CORS()) // CORS 跨域

	// 6. 注册路由
	setupRoutes(r)

	// 7. 启动服务器
	port := ":" + config.AppConfig.ServerPort
	log.Printf("服务器启动在端口 %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// setupRoutes 设置所有路由
func setupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// ==================== 认证相关接口 ====================
	authGroup := api.Group("/auth")
	{
		// POST /api/auth/login - 用户登录
		authGroup.POST("/login", auth.Login)

		// POST /api/auth/register - 用户注册
		authGroup.POST("/register", auth.Register)

		// POST /api/auth/logout - 退出登录（需要鉴权）
		authGroup.POST("/logout", middleware.AuthRequired(), auth.Logout)

		// GET /api/auth/current-user - 获取当前用户信息（需要鉴权）
		authGroup.GET("/current-user", middleware.AuthRequired(), auth.GetCurrentUser)
	}

	// ==================== 主页相关接口 ====================
	homeGroup := api.Group("/home")
	{
		// GET /api/home/statistics - 获取平台统计数据（可选鉴权）
		homeGroup.GET("/statistics", home.GetStatistics)
	}

	// ==================== 讲师端接口 ====================
	teacherGroup := api.Group("/teacher")
	teacherGroup.Use(middleware.AuthRequired(), middleware.RoleRequired("讲师"))
	{

		// GET /api/teacher/schedule - 获取讲师授课表
		teacherGroup.GET("/schedule", teacher.GetSchedule)

		// GET /api/teacher/pending-evaluations - 获取待评分学员列表
		teacherGroup.GET("/pending-evaluations", teacher.GetPendingEvaluations)

		// POST /api/teacher/submit-grading - 提交学员评分
		teacherGroup.POST("/submit-grading", teacher.SubmitGrading)

		// GET /api/teacher/course-statistics - 获取课程成绩统计
		teacherGroup.GET("/course-statistics", teacher.GetCourseStatistics)

		// GET /api/teacher/teaching-statistics - 获取讲师授课统计
		teacherGroup.GET("/teaching-statistics", teacher.GetTeachingStatistics)
	}

	// ==================== 员工端接口 ====================
	employeeGroup := api.Group("/employee")
	employeeGroup.Use(middleware.AuthRequired(), middleware.RoleRequired("员工"))
	{
		// GET /api/employee/schedule - 获取员工课程表
		employeeGroup.GET("/schedule", employee.GetSchedule)

		// GET /api/employee/pending-evaluations - 获取待自评课程列表
		employeeGroup.GET("/pending-evaluations", employee.GetPendingEvaluations)

		// POST /api/employee/submit-evaluation - 提交课程自评
		employeeGroup.POST("/submit-evaluation", employee.SubmitEvaluation)

		// GET /api/employee/scores - 获取员工成绩列表
		employeeGroup.GET("/scores", employee.GetScores)

		// GET /api/employee/course-type-scores - 获取课程类型成绩分析
		employeeGroup.GET("/course-type-scores", employee.GetCourseTypeScores)

		// GET /api/employee/learning-progress - 获取员工学习进度
		employeeGroup.GET("/learning-progress", employee.GetLearningProgress)
	}

	// ==================== 课程大纲制定者端接口 ====================
	plannerGroup := api.Group("/planner")
	plannerGroup.Use(middleware.AuthRequired(), middleware.RoleRequired("课程大纲制定者"))
	{
		// GET /api/planner/teachers - 获取讲师列表（用于选择）
		plannerGroup.GET("/teachers", planner.GetTeachersList)

		// GET /api/planner/employees - 获取员工列表（用于选择）
		plannerGroup.GET("/employees", planner.GetEmployeesList)

		// GET /api/planner/plans - 获取培训计划列表
		plannerGroup.GET("/plans", planner.GetPlansList)

		// POST /api/planner/plans - 创建培训计划
		plannerGroup.POST("/plans", planner.CreatePlan)

		// GET /api/planner/plans/:planId - 获取培训计划详情
		plannerGroup.GET("/plans/:planId", planner.GetPlanDetail)

		// PUT /api/planner/plans/:planId - 修改培训计划
		plannerGroup.PUT("/plans/:planId", planner.UpdatePlan)

		// DELETE /api/planner/plans/:planId - 删除培训计划
		plannerGroup.DELETE("/plans/:planId", planner.DeletePlan)

		// POST /api/planner/plans/:planId/employees - 为培训计划添加员工
		plannerGroup.POST("/plans/:planId/employees", planner.AddEmployeesToPlan)

		// DELETE /api/planner/plans/:planId/employees/:employeeId - 从培训计划移除员工
		plannerGroup.DELETE("/plans/:planId/employees/:employeeId", planner.RemoveEmployeeFromPlan)

		// GET /api/planner/courses - 获取课程列表
		plannerGroup.GET("/courses", planner.GetCoursesList)

		// POST /api/planner/courses - 创建课程
		plannerGroup.POST("/courses", planner.CreateCourse)

		// PUT /api/planner/courses/:courseId - 修改课程
		plannerGroup.PUT("/courses/:courseId", planner.UpdateCourse)

		// DELETE /api/planner/courses/:courseId - 删除课程
		plannerGroup.DELETE("/courses/:courseId", planner.DeleteCourse)

		// GET /api/planner/course-items - 获取课程安排列表
		plannerGroup.GET("/course-items", planner.GetCourseItemsList)

		// POST /api/planner/course-items - 创建课程安排
		plannerGroup.POST("/course-items", planner.CreateCourseItem)

		// PUT /api/planner/course-items/:itemId - 修改课程安排
		plannerGroup.PUT("/course-items/:itemId", planner.UpdateCourseItem)

		// DELETE /api/planner/course-items/:itemId - 删除课程安排
		plannerGroup.DELETE("/course-items/:itemId", planner.DeleteCourseItem)

		// GET /api/planner/analytics - 获取平台数据分析
		plannerGroup.GET("/analytics", planner.GetAnalytics)

		// GET /api/planner/employees/:employeeId/scores - 获取员工成绩详情
		plannerGroup.GET("/employees/:employeeId/scores", planner.GetEmployeeScores)

		// GET /api/planner/courses/:courseId/evaluations - 获取课程评价详情
		plannerGroup.GET("/courses/:courseId/evaluations", planner.GetCourseEvaluations)
	}

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "服务运行正常",
		})
	})
}
