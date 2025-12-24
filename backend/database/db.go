package database

import (
	"fmt"
	"log"
	"time"
	"backend/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	cfg := config.AppConfig

	// 构建 DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	// 配置 GORM 日志
	gormLogger := logger.Default.LogMode(logger.Info)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	// 获取底层的 *sql.DB 对象
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期

	DB = db
	log.Println("数据库连接成功")

	// 自动迁移表结构
	if err := AutoMigrate(); err != nil {
		return fmt.Errorf("数据库表迁移失败: %v", err)
	}

	return nil
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() error {
	log.Println("开始自动迁移数据库表...")

	err := DB.AutoMigrate(
		&Person{},
		&Account{},
		&TrainingPlan{},
		&Course{},
		&PlanCourseItem{},
		&AttendanceEvaluation{},
		&PlanEmployee{},
	)

	if err != nil {
		return err
	}

	log.Println("数据库表迁移完成")
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
