package database

import (
	"time"

)

// Person 人员表
type Person struct {
	PersonID int64  `gorm:"primaryKey;column:person_id" json:"personId"`
	Name     string `gorm:"column:name;size:20;not null" json:"name"`
	Role     string `gorm:"column:role;size:7;not null;comment:角色：员工/讲师/课程大纲制定者" json:"role"`
}

func (Person) TableName() string {
	return "person"
}

// Account 账号表
type Account struct {
	AccountID    int64  `gorm:"primaryKey;column:account_id" json:"accountId"`
	PersonID     int64  `gorm:"column:person_id;not null;index" json:"personId"`
	LoginName    string `gorm:"column:login_name;size:20;not null;uniqueIndex" json:"loginName"`
	PasswordHash string `gorm:"column:password_hash;size:255;not null" json:"-"`
	Person       Person `gorm:"foreignKey:PersonID;references:PersonID;constraint:OnDelete:CASCADE"`
}

func (Account) TableName() string {
	return "account"
}

// TrainingPlan 培训计划表
type TrainingPlan struct {
	PlanID            int64     `gorm:"primaryKey;column:plan_id" json:"planId"`
	PlanName          string    `gorm:"column:plan_name;size:50;not null" json:"planName"`
	PlanStatus        string    `gorm:"column:plan_status;size:3;not null;comment:规划中/进行中/已完成" json:"planStatus"`
	PlanStartDatetime time.Time `gorm:"column:plan_start_datetime;not null" json:"planStartDatetime"`
	PlanEndDatetime   time.Time `gorm:"column:plan_end_datetime;not null" json:"planEndDatetime"`
	CreatorID         int64     `gorm:"column:creator_id;not null;index" json:"creatorId"`
	Creator           Person    `gorm:"foreignKey:CreatorID;references:PersonID"`
}

func (TrainingPlan) TableName() string {
	return "training_plan"
}

// Course 培训课程表
type Course struct {
	CourseID      int64  `gorm:"primaryKey;column:course_id" json:"courseId"`
	CourseName    string `gorm:"column:course_name;size:50;not null" json:"courseName"`
	CourseDesc    string `gorm:"column:course_desc;size:100" json:"courseDesc"`
	CourseRequire string `gorm:"column:course_require;size:500" json:"courseRequire"`
	CourseClass   string `gorm:"column:course_class;size:20;not null" json:"courseClass"`
	TeacherID     int64  `gorm:"column:teacher_id;not null;index" json:"teacherId"`
	Teacher       Person `gorm:"foreignKey:TeacherID;references:PersonID"`
}

func (Course) TableName() string {
	return "course"
}

// PlanCourseItem 培训课程安排表
type PlanCourseItem struct {
	ItemID         int64         `gorm:"primaryKey;column:item_id" json:"itemId"`
	PlanID         int64         `gorm:"column:plan_id;not null;index" json:"planId"`
	CourseID       int64         `gorm:"column:course_id;not null;index" json:"courseId"`
	ClassDate      time.Time     `gorm:"column:class_date;type:date;not null" json:"classDate"`
	ClassBeginTime string        `gorm:"column:class_begin_time;type:varchar(8);not null" json:"classBeginTime"`
	ClassEndTime   string        `gorm:"column:class_end_time;type:varchar(8);not null" json:"classEndTime"`
	Location       string        `gorm:"column:location;size:100;not null" json:"location"`
	Plan           TrainingPlan  `gorm:"foreignKey:PlanID;references:PlanID"`
	Course         Course        `gorm:"foreignKey:CourseID;references:CourseID"`
}

func (PlanCourseItem) TableName() string {
	return "plan_course_item"
}

// AttendanceEvaluation 参与和评价表
type AttendanceEvaluation struct {
	PersonID       int64   `gorm:"primaryKey;column:person_id" json:"personId"`
	ItemID         int64   `gorm:"primaryKey;column:item_id" json:"itemId"`
	SelfScore      float64 `gorm:"column:self_score;type:float" json:"selfScore"`
	SelfComment    string  `gorm:"column:self_comment;type:text" json:"selfComment"`
	TeacherScore   float64 `gorm:"column:teacher_score;type:float" json:"teacherScore"`
	TeacherComment string  `gorm:"column:teacher_comment;type:text" json:"teacherComment"`
	ScoreRatio     float64 `gorm:"column:score_ratio;type:float;comment:讲师评分占比" json:"scoreRatio"`
	Person         Person  `gorm:"foreignKey:PersonID;references:PersonID"`
	Item           PlanCourseItem `gorm:"foreignKey:ItemID;references:ItemID"`
}

func (AttendanceEvaluation) TableName() string {
	return "attendance_evaluation"
}

// PlanEmployee 规划员工表
type PlanEmployee struct {
	PlanID   int64        `gorm:"primaryKey;column:plan_id" json:"planId"`
	PersonID int64        `gorm:"primaryKey;column:person_id" json:"personId"`
	Plan     TrainingPlan `gorm:"foreignKey:PlanID;references:PlanID"`
	Person   Person       `gorm:"foreignKey:PersonID;references:PersonID"`
}

func (PlanEmployee) TableName() string {
	return "plan_employee"
}

// Session 会话表（用于简单鉴权）
type Session struct {
	SessionID string    `gorm:"primaryKey;column:session_id;size:64" json:"sessionId"`
	PersonID  int64     `gorm:"column:person_id;not null;index" json:"personId"`
	Role      string    `gorm:"column:role;size:10;not null" json:"role"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null" json:"expiresAt"`
}

func (Session) TableName() string {
	return "sessions"
}
