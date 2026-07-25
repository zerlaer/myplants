package controller

import (
	"github.com/gin-gonic/gin"

	"myplants/internal/database"
	"myplants/internal/model"
	"myplants/internal/response"
)

// ListNotes 笔记列表
func ListNotes(c *gin.Context) {
	plantID := c.Query("plant_id")
	var notes []model.Note
	query := database.DB.Order("created_at DESC")
	if plantID != "" {
		query = query.Where("plant_id = ?", plantID)
	}
	if err := query.Find(&notes).Error; err != nil {
		response.Fail(c, "查询失败")
		return
	}
	response.OK(c, notes)
}

// CreateNote 创建笔记
func CreateNote(c *gin.Context) {
	var req model.NoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误: "+err.Error())
		return
	}
	note := model.Note{
		PlantID: req.PlantID,
		Content: req.Content,
	}
	if err := database.DB.Create(&note).Error; err != nil {
		response.Fail(c, "创建失败: "+err.Error())
		return
	}
	response.OKWithMsg(c, "添加成功", note)
}

// DeleteNote 删除笔记
func DeleteNote(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&model.Note{}, id).Error; err != nil {
		response.Fail(c, "删除失败")
		return
	}
	response.OKWithMsg(c, "删除成功", nil)
}
