package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config 配置结构体
type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	ServerPort    string
	SessionSecret string
}

var AppConfig *Config

// LoadConfig 加载配置
func LoadConfig() error {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("未找到 .env 文件，使用默认配置")
	}

	AppConfig = &Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "3306"),
		DBUser:        getEnv("DB_USER", "root"),
		DBPassword:    getEnv("DB_PASSWORD", "GreatSQL@2025"),
		DBName:        getEnv("DB_NAME", "training_system"),
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		SessionSecret: getEnv("SESSION_SECRET", "default-secret-key"),
	}

	log.Println("配置加载成功")
	return nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
