package planner

import (
	"net/http"
	"strconv"
	"backend/database"

	"github.com/gin-gonic/gin"
)

// CourseItemDetail 课程安排详情
type CourseItemDetail struct {
	ItemID         int64  `json:"itemId"`
	CourseID       int64  `json:"courseId"`
	CourseName     string `json:"courseName"`
	CourseClass    string `json:"courseClass"`
	TeacherID      int64  `json:"teacherId"`
	TeacherName    string `json:"teacherName"`
	ClassDate      string `json:"classDate"`
	ClassBeginTime string `json:"classBeginTime"`
	ClassEndTime   string `json:"classEndTime"`
	Location       string `json:"location"`
}

// EmployeeDetail 员工详情
type EmployeeDetail struct {
	PersonID int64  `json:"personId"`
	Name     string `json:"name"`
}

// PlanDetailResponse 培训计划详情响应
type PlanDetailResponse struct {
	PlanID            int64              `json:"planId"`
	PlanName          string             `json:"planName"`
	PlanStatus        string             `json:"planStatus"`
	PlanStartDatetime string             `json:"planStartDatetime"`
	PlanEndDatetime   string             `json:"planEndDatetime"`
	CreatorID         int64              `json:"creatorId"`
	CreatorName       string             `json:"creatorName"`
	CourseItems       []CourseItemDetail `json:"courseItems"`
	Employees         []EmployeeDetail   `json:"employees"`
}

// GetPlanDetail 获取培训计划详情（5.5接口）
func GetPlanDetail(c *gin.Context) {
	// 获取路径参数
	planIDStr := c.Param("planId")
	planID, err := strconv.ParseInt(planIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的计划ID",
			"data":    nil,
		})
		return
	}

	// 查询培训计划基本信息
	var plan database.TrainingPlan
	if err := database.DB.Preload("Creator").Where("plan_id = ?", planID).First(&plan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "培训计划不存在",
			"data":    nil,
		})
		return
	}

	// 查询课程安排
	var courseItems []CourseItemDetail
	err = database.DB.Table("plan_course_item pci").
		Select(`
			pci.item_id,
			pci.course_id,
			c.course_name,
			c.course_class,
			c.teacher_id,
			p.name as teacher_name,
			DATE_FORMAT(pci.class_date, '%Y-%m-%d') as class_date,
			TIME_FORMAT(pci.class_begin_time, '%H:%i:%s') as class_begin_time,
			TIME_FORMAT(pci.class_end_time, '%H:%i:%s') as class_end_time,
			pci.location
		`).
		Joins("JOIN course c ON pci.course_id = c.course_id").
		Joins("JOIN person p ON c.teacher_id = p.person_id").
		Where("pci.plan_id = ?", planID).
		Order("pci.class_date ASC, pci.class_begin_time ASC").
		Scan(&courseItems).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询课程安排失败",
			"data":    nil,
		})
		return
	}

	// 查询关联员工
	var employees []EmployeeDetail
	err = database.DB.Table("plan_employee pe").
		Select("pe.person_id, p.name").
		Joins("JOIN person p ON pe.person_id = p.person_id").
		Where("pe.plan_id = ?", planID).
		Order("p.name ASC").
		Scan(&employees).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询关联员工失败",
			"data":    nil,
		})
		return
	}

	// 构建响应
	response := PlanDetailResponse{
		PlanID:            plan.PlanID,
		PlanName:          plan.PlanName,
		PlanStatus:        plan.PlanStatus,
		PlanStartDatetime: plan.PlanStartDatetime.Format("2006-01-02 15:04:05"),
		PlanEndDatetime:   plan.PlanEndDatetime.Format("2006-01-02 15:04:05"),
		CreatorID:         plan.CreatorID,
		CreatorName:       plan.Creator.Name,
		CourseItems:       courseItems,
		Employees:         employees,
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    response,
	})
}
