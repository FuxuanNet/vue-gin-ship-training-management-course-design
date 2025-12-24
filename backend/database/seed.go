package database

import (
	"log"
	"time"

)

// SeedTestData 插入测试数据（可选）
func SeedTestData() error {
	log.Println("开始插入测试数据...")

	// 检查是否已有数据
	var count int64
	DB.Model(&Person{}).Count(&count)
	if count > 0 {
		log.Println("数据已存在，跳过测试数据插入")
		return nil
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

	// 插入账号数据（明文密码，实际使用时需在登录时哈希验证）
	accounts := []Account{
		{PersonID: 1, LoginName: "admin", PasswordHash: "123456"},
		{PersonID: 2, LoginName: "teacher1", PasswordHash: "123456"},
		{PersonID: 3, LoginName: "teacher2", PasswordHash: "123456"},
		{PersonID: 4, LoginName: "employee1", PasswordHash: "123456"},
		{PersonID: 5, LoginName: "employee2", PasswordHash: "123456"},
	}
	if err := DB.Create(&accounts).Error; err != nil {
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

	log.Println("测试数据插入完成！")
	return nil
}
