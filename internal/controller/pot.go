package controller

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"myplants/internal/database"
	"myplants/internal/model"
	"myplants/internal/response"
)

// ListPots 花盆列表
func ListPots(c *gin.Context) {
	var pots []model.Pot
	query := database.DB.Order("created_at DESC")
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if plantID := c.Query("plant_id"); plantID != "" {
		query = query.Where("plant_id = ?", plantID)
	}
	if err := query.Find(&pots).Error; err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, pots)
}

// CreatePot 创建花盆
func CreatePot(c *gin.Context) {
	var req model.PotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误: "+err.Error())
		return
	}
	pot := model.Pot{
		PlantID:  req.PlantID,
		Name:     req.Name,
		Size:     req.Size,
		Diameter: req.Diameter,
		Height:   req.Height,
		Material: req.Material,
		Color:    req.Color,
		Status:   req.Status,
		Remark:   req.Remark,
	}
	if pot.Status == "" {
		pot.Status = "使用中"
	}
	if err := database.DB.Create(&pot).Error; err != nil {
		response.Fail(c, "创建失败: "+err.Error())
		return
	}
	response.OKWithMsg(c, "创建成功", pot)
}

// UpdatePot 更新花盆
func UpdatePot(c *gin.Context) {
	id := c.Param("id")
	var pot model.Pot
	if err := database.DB.First(&pot, id).Error; err != nil {
		response.Fail(c, "花盆不存在")
		return
	}
	var req model.PotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	pot.PlantID = req.PlantID
	pot.Name = req.Name
	pot.Size = req.Size
	pot.Diameter = req.Diameter
	pot.Height = req.Height
	pot.Material = req.Material
	pot.Color = req.Color
	pot.Status = req.Status
	pot.Remark = req.Remark
	database.DB.Save(&pot)
	response.OKWithMsg(c, "更新成功", pot)
}

// DeletePot 删除花盆
func DeletePot(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&model.Pot{}, id).Error; err != nil {
		response.Fail(c, "删除失败")
		return
	}
	response.OKWithMsg(c, "删除成功", nil)
}

// ListRepottings 换盆记录列表
func ListRepottings(c *gin.Context) {
	var records []model.Repotting
	query := database.DB.Order("repot_time DESC")
	if plantID := c.Query("plant_id"); plantID != "" {
		query = query.Where("plant_id = ?", plantID)
	}
	if err := query.Find(&records).Error; err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, records)
}

// CreateRepotting 创建换盆记录(同时更新关联花盆状态)
func CreateRepotting(c *gin.Context) {
	var req model.RepottingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误: "+err.Error())
		return
	}
	repotTime := time.Now()
	if req.RepotTime != nil {
		repotTime = *req.RepotTime
	}

	var fromName, toName string
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 获取花盆名称
		if req.FromPotID != nil {
			var p model.Pot
			if err := tx.First(&p, *req.FromPotID).Error; err == nil {
				fromName = p.Name
				if fromName == "" {
					fromName = "花盆#" + strconv.Itoa(int(p.ID))
				}
			}
		}
		if req.ToPotID != nil {
			var p model.Pot
			if err := tx.First(&p, *req.ToPotID).Error; err == nil {
				toName = p.Name
				if toName == "" {
					toName = "花盆#" + strconv.Itoa(int(p.ID))
				}
			}
		}

		record := model.Repotting{
			PlantID:     req.PlantID,
			FromPotID:   req.FromPotID,
			ToPotID:     req.ToPotID,
			FromPotName: fromName,
			ToPotName:   toName,
			RepotTime:   repotTime,
			Remark:      req.Remark,
		}
		if err := tx.Create(&record).Error; err != nil {
			return err
		}

		// 旧花盆设为空闲
		if req.FromPotID != nil {
			tx.Model(&model.Pot{}).Where("id = ?", *req.FromPotID).Updates(map[string]interface{}{"plant_id": nil, "status": "空闲"})
		}
		// 新花盆关联到植物
		if req.ToPotID != nil {
			tx.Model(&model.Pot{}).Where("id = ?", *req.ToPotID).Updates(map[string]interface{}{"plant_id": req.PlantID, "status": "使用中"})
		}
		return nil
	})
	if err != nil {
		response.Fail(c, "记录失败: "+err.Error())
		return
	}
	response.OKWithMsg(c, "换盆记录已保存", nil)
}

// DeleteRepotting 删除换盆记录
func DeleteRepotting(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&model.Repotting{}, id).Error; err != nil {
		response.Fail(c, "删除失败")
		return
	}
	response.OKWithMsg(c, "删除成功", nil)
}
