package database

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// SeedTestAccounts 仅插入测试账号数据（用于开发测试）
func SeedTestAccounts() error {
	log.Println("开始插入测试账号...")

	// 检查是否已有账号数据
	var count int64
	DB.Model(&Account{}).Count(&count)
	if count > 0 {
		log.Println("账号数据已存在，跳过测试账号插入")
		return nil
	}

	// 密码统一为 "123456"，需要进行 bcrypt 哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 插入人员数据
	persons := []Person{
		{PersonID: 1, Name: "张主管", Role: "课程大纲制定者"},
		{PersonID: 2, Name: "李老师", Role: "讲师"},
		{PersonID: 3, Name: "王老师", Role: "讲师"},
		{PersonID: 4, Name: "赵员工", Role: "员工"},
		{PersonID: 5, Name: "钱员工", Role: "员工"},
	}
	if err := DB.Create(&persons).Error; err != nil {
		return err
	}

	// 插入账号数据（使用哈希后的密码）
	accounts := []Account{
		{PersonID: 1, LoginName: "planner", PasswordHash: string(hashedPassword)},
		{PersonID: 2, LoginName: "teacher", PasswordHash: string(hashedPassword)},
		{PersonID: 3, LoginName: "teacher2", PasswordHash: string(hashedPassword)},
		{PersonID: 4, LoginName: "employee", PasswordHash: string(hashedPassword)},
		{PersonID: 5, LoginName: "employee2", PasswordHash: string(hashedPassword)},
	}
	if err := DB.Create(&accounts).Error; err != nil {
		return err
	}

	log.Println("测试账号插入完成！")
	log.Println("测试账号列表（密码均为 123456）：")
	log.Println("  - planner (课程大纲制定者)")
	log.Println("  - teacher (讲师)")
	log.Println("  - teacher2 (讲师)")
	log.Println("  - employee (员工)")
	log.Println("  - employee2 (员工)")
	return nil
}

// SeedTestData 插入完整测试数据（包含计划、课程等，可选）
func SeedTestData() error {
	log.Println("开始插入完整测试数据...")

	// 检查是否已有数据
	var count int64
	DB.Model(&Person{}).Count(&count)
	if count > 0 {
		log.Println("数据已存在，跳过测试数据插入")
		return nil
	}

	// 先插入测试账号
	if err := SeedTestAccounts(); err != nil {
		return err
	}

	// 插入培训计划数据
	plans := []TrainingPlan{
		{
			PlanName:          "2024年新员工培训计划",
			PlanStatus:        "进行中",
			PlanStartDatetime: time.Date(2024, 1, 1, 9, 0, 0, 0, time.Local),
			PlanEndDatetime:   time.Date(2024, 6, 30, 18, 0, 0, 0, time.Local),
			CreatorID:         1,
		},
		{
			PlanName:          "船舶安全专项培训",
			PlanStatus:        "规划中",
			PlanStartDatetime: time.Date(2024, 7, 1, 9, 0, 0, 0, time.Local),
			PlanEndDatetime:   time.Date(2024, 9, 30, 18, 0, 0, 0, time.Local),
			CreatorID:         1,
		},
	}
	if err := DB.Create(&plans).Error; err != nil {
		return err
	}

	// 插入课程数据
	courses := []Course{
		{
			CourseName:    "船舶安全基础",
			CourseDesc:    "介绍船舶安全的基本知识和操作规范",
			CourseRequire: "无特殊要求",
			CourseClass:   "安全培训",
			TeacherID:     2,
		},
		{
			CourseName:    "船舶结构力学",
			CourseDesc:    "学习船舶结构设计和力学分析",
			CourseRequire: "需有基础力学知识",
			CourseClass:   "专业技能",
			TeacherID:     3,
		},
	}
	if err := DB.Create(&courses).Error; err != nil {
		return err
	}

	// 插入课程安排数据
	items := []PlanCourseItem{
		{
			PlanID:         1,
			CourseID:       1,
			ClassDate:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.Local),
			ClassBeginTime: time.Date(0, 1, 1, 9, 0, 0, 0, time.Local),
			ClassEndTime:   time.Date(0, 1, 1, 11, 0, 0, 0, time.Local),
			Location:       "培训室A",
		},
		{
			PlanID:         1,
			CourseID:       2,
			ClassDate:      time.Date(2024, 1, 16, 0, 0, 0, 0, time.Local),
			ClassBeginTime: time.Date(0, 1, 1, 14, 0, 0, 0, time.Local),
			ClassEndTime:   time.Date(0, 1, 1, 16, 0, 0, 0, time.Local),
			Location:       "培训室B",
		},
	}
	if err := DB.Create(&items).Error; err != nil {
		return err
	}

	// 插入规划员工数据
	planEmployees := []PlanEmployee{
		{PlanID: 1, PersonID: 4},
		{PlanID: 1, PersonID: 5},
	}
	if err := DB.Create(&planEmployees).Error; err != nil {
		return err
	}

	log.Println("完整测试数据插入完成！")
	return nil
}
