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

// ListCareRecords 养护记录列表
func ListCareRecords(c *gin.Context) {
	plantID := c.Query("plant_id")
	careType := c.Query("type")
	var records []model.CareRecord
	query := database.DB.Order("record_time DESC")
	if plantID != "" {
		query = query.Where("plant_id = ?", plantID)
	}
	if careType != "" {
		query = query.Where("type = ?", careType)
	}
	if err := query.Find(&records).Error; err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, records)
}

// CreateCareRecord 创建养护记录(同时更新植物的最后养护时间)
func CreateCareRecord(c *gin.Context) {
	var req model.CareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误: "+err.Error())
		return
	}
	if req.Type != model.CareTypeWater && req.Type != model.CareTypeFertilize && req.Type != model.CareTypeSpray {
		response.Fail(c, "无效的养护类型")
		return
	}

	recordTime := time.Now()
	if req.RecordTime != nil {
		recordTime = *req.RecordTime
	}

	var record model.CareRecord
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var plant model.Plant
		if err := tx.First(&plant, req.PlantID).Error; err != nil {
			return err
		}
		record = model.CareRecord{
			PlantID:    req.PlantID,
			Type:       req.Type,
			RecordTime: recordTime,
			Remark:     req.Remark,
		}
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		switch req.Type {
		case model.CareTypeWater:
			plant.LastWateredAt = &recordTime
		case model.CareTypeFertilize:
			plant.LastFertilizedAt = &recordTime
		case model.CareTypeSpray:
			plant.LastSprayedAt = &recordTime
		}
		return tx.Save(&plant).Error
	})
	if err != nil {
		response.Fail(c, "记录失败: "+err.Error())
		return
	}

	msg := "记录成功"
	switch req.Type {
	case model.CareTypeWater:
		msg = "浇水完成,已记录"
	case model.CareTypeFertilize:
		msg = "施肥完成,已记录"
	case model.CareTypeSpray:
		msg = "打药完成,已记录"
	}
	response.OKWithMsg(c, msg, record)
}

// OneClickCare 一键养护(浇水/施肥/打药), 用当前时间
func OneClickCare(c *gin.Context) {
	plantID := c.Param("id")
	careType := c.Param("type")
	pid, err := strconv.Atoi(plantID)
	if err != nil {
		response.Fail(c, "无效的植物ID")
		return
	}
	if careType != model.CareTypeWater && careType != model.CareTypeFertilize && careType != model.CareTypeSpray {
		response.Fail(c, "无效的养护类型")
		return
	}

	now := time.Now()
	var record model.CareRecord
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		var plant model.Plant
		if err := tx.First(&plant, pid).Error; err != nil {
			return err
		}
		record = model.CareRecord{
			PlantID:    uint(pid),
			Type:       careType,
			RecordTime: now,
		}
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		switch careType {
		case model.CareTypeWater:
			plant.LastWateredAt = &now
		case model.CareTypeFertilize:
			plant.LastFertilizedAt = &now
		case model.CareTypeSpray:
			plant.LastSprayedAt = &now
		}
		return tx.Save(&plant).Error
	})
	if err != nil {
		response.Fail(c, "操作失败: "+err.Error())
		return
	}

	msg := "操作完成"
	switch careType {
	case model.CareTypeWater:
		msg = "浇水完成 💧"
	case model.CareTypeFertilize:
		msg = "施肥完成 🌱"
	case model.CareTypeSpray:
		msg = "打药完成 🛡️"
	}
	response.OKWithMsg(c, msg, record)
}

// BatchCare 批量一键养护
func BatchCare(c *gin.Context) {
	var req struct {
		PlantIDs []uint `json:"plant_ids" binding:"required"`
		Type     string `json:"type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误: "+err.Error())
		return
	}
	if req.Type != model.CareTypeWater && req.Type != model.CareTypeFertilize && req.Type != model.CareTypeSpray {
		response.Fail(c, "无效的养护类型")
		return
	}

	now := time.Now()
	var records []model.CareRecord
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		for _, pid := range req.PlantIDs {
			var plant model.Plant
			if err := tx.First(&plant, pid).Error; err != nil {
				continue
			}
			record := model.CareRecord{
				PlantID:    pid,
				Type:       req.Type,
				RecordTime: now,
			}
			if err := tx.Create(&record).Error; err != nil {
				continue
			}
			records = append(records, record)
			switch req.Type {
			case model.CareTypeWater:
				plant.LastWateredAt = &now
			case model.CareTypeFertilize:
				plant.LastFertilizedAt = &now
			case model.CareTypeSpray:
				plant.LastSprayedAt = &now
			}
			tx.Save(&plant)
		}
		return nil
	})
	if err != nil {
		response.Fail(c, "操作失败: "+err.Error())
		return
	}
	response.OKWithMsg(c, "批量操作完成", records)
}

// DeleteCareRecord 删除养护记录
func DeleteCareRecord(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&model.CareRecord{}, id).Error; err != nil {
		response.Fail(c, "删除失败")
		return
	}
	response.OKWithMsg(c, "删除成功", nil)
}
