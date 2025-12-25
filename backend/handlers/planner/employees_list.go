package planner

import (
	"net/http"
	"backend/database"

	"github.com/gin-gonic/gin"
)

// GetEmployeesList 获取所有员工列表（用于选择）
func GetEmployeesList(c *gin.Context) {
	// 查询所有角色为"员工"的人员
	var employees []database.Person
	if err := database.DB.Where("role = ?", "员工").Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
			"data":    nil,
		})
		return
	}

	// 构建响应数据
	type EmployeeResponse struct {
		PersonID   int64  `json:"personId"`
		PersonName string `json:"personName"`
	}

	list := make([]EmployeeResponse, 0, len(employees))
	for _, emp := range employees {
		list = append(list, EmployeeResponse{
			PersonID:   emp.PersonID,
			PersonName: emp.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    list,
	})
}
