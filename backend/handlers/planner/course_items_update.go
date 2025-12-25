package planner

import (
	"net/http"
	"strconv"
	"time"
	"backend/database"

	"github.com/gin-gonic/gin"
)

// UpdateCourseItem 修改课程安排（接口5.14）
func UpdateCourseItem(c *gin.Context) {
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

	// 解析请求体
	var req struct {
		ClassDate      *string `json:"classDate"`
		ClassBeginTime *string `json:"classBeginTime"`
		ClassEndTime   *string `json:"classEndTime"`
		Location       *string `json:"location"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误：" + err.Error(),
			"data":    nil,
		})
		return
	}

	// 构建更新映射
	updates := make(map[string]interface{})

	if req.ClassDate != nil {
		classDate, err := time.Parse("2006-01-02", *req.ClassDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "日期格式错误，请使用 YYYY-MM-DD 格式",
				"data":    nil,
			})
			return
		}
		updates["class_date"] = classDate
	}

	if req.ClassBeginTime != nil {
		// 验证时间格式
		if _, err := time.Parse("15:04:05", *req.ClassBeginTime); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "开始时间格式错误，请使用 HH:mm:ss 格式",
				"data":    nil,
			})
			return
		}
		updates["class_begin_time"] = *req.ClassBeginTime
	}

	if req.ClassEndTime != nil {
		// 验证时间格式
		if _, err := time.Parse("15:04:05", *req.ClassEndTime); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "结束时间格式错误，请使用 HH:mm:ss 格式",
				"data":    nil,
			})
			return
		}
		updates["class_end_time"] = *req.ClassEndTime
	}

	if req.Location != nil {
		if len(*req.Location) > 100 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "上课地点长度不能超过100字符",
				"data":    nil,
			})
			return
		}
		updates["location"] = *req.Location
	}

	// 如果没有更新内容，直接返回
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "没有提供更新内容",
			"data":    nil,
		})
		return
	}

	// 执行更新
	if err := database.DB.Model(&item).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新课程安排失败",
			"data":    nil,
		})
		return
	}

	// 重新查询更新后的课程安排（带关联数据）
	database.DB.Preload("Plan").Preload("Course").Where("item_id = ?", itemId).First(&item)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "修改成功",
		"data": gin.H{
			"itemId":         item.ItemID,
			"planId":         item.PlanID,
			"planName":       item.Plan.PlanName,
			"courseId":       item.CourseID,
			"courseName":     item.Course.CourseName,
			"classDate":      item.ClassDate.Format("2006-01-02"),
			"classBeginTime": item.ClassBeginTime,
			"classEndTime":   item.ClassEndTime,
			"location":       item.Location,
		},
	})
}
