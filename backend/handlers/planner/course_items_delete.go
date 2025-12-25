package planner

import (
	"net/http"
	"strconv"
	"backend/database"

	"github.com/gin-gonic/gin"
)

// DeleteCourseItem 删除课程安排（接口5.15）
func DeleteCourseItem(c *gin.Context) {
	// 获取路径参数 itemId
	itemIdStr := c.Param("itemId")
	itemId, err := strconv.ParseInt(itemIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的课程安排ID",
			"data":    nil,
		})
		return
	}

	// 获取查询参数 force
	forceStr := c.Query("force")
	force := forceStr == "true"

	// 验证课程安排是否存在
	var item database.PlanCourseItem
	if err := database.DB.Where("item_id = ?", itemId).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "课程安排不存在",
			"data":    nil,
		})
		return
	}

	// 检查是否有评价记录
	var evaluationCount int64
	database.DB.Model(&database.AttendanceEvaluation{}).
		Where("item_id = ?", itemId).
		Count(&evaluationCount)

	// 如果有评价记录且不是强制删除，返回警告
	if evaluationCount > 0 && !force {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "该课程安排已有员工评价记录，建议不要删除。如需强制删除，请添加force=true参数",
			"data": gin.H{
				"evaluationCount": evaluationCount,
			},
		})
		return
	}

	// 如果是强制删除，先删除评价记录
	if force && evaluationCount > 0 {
		database.DB.Where("item_id = ?", itemId).Delete(&database.AttendanceEvaluation{})
	}

	// 删除课程安排
	if err := database.DB.Delete(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除课程安排失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
		"data":    nil,
	})
}
