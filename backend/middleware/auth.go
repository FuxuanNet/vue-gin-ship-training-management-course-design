package middleware

import (
	"net/http"
	"time"
	"backend/database"

	"github.com/gin-gonic/gin"

)

// AuthRequired 简单鉴权中间件（验证 Session-ID）
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.GetHeader("Session-ID")
		if sessionID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未登录或登录已过期",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 查询会话
		var session database.Session
		if err := database.DB.Where("session_id = ?", sessionID).First(&session).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "会话无效或已过期",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 检查会话是否过期
		if time.Now().After(session.ExpiresAt) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "会话已过期，请重新登录",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("personId", session.PersonID)
		c.Set("role", session.Role)

		c.Next()
	}
}

// RoleRequired 角色验证中间件
func RoleRequired(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "无权限访问",
				"data":    nil,
			})
			c.Abort()
			return
		}

		roleStr := role.(string)
		allowed := false
		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "无权限访问",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
