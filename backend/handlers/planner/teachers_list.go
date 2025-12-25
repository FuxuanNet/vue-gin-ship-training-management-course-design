package planner

import (
	"net/http"
	"backend/database"

	"github.com/gin-gonic/gin"
)

// GetTeachersList 获取讲师列表（用于课程管理选择讲师）
func GetTeachersList(c *gin.Context) {
	// 查询所有讲师
	var teachers []database.Person
	if err := database.DB.Where("role = ?", "讲师").Find(&teachers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询讲师列表失败",
			"data":    nil,
		})
		return
	}

	// 构建响应数据
	type TeacherResponse struct {
		PersonID int64  `json:"personId"`
		Name     string `json:"name"`
	}

	list := make([]TeacherResponse, 0, len(teachers))
	for _, teacher := range teachers {
		list = append(list, TeacherResponse{
			PersonID: teacher.PersonID,
			Name:     teacher.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"total": len(list),
			"list":  list,
		},
	})
}
