package controller

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"myplants/internal/config"
	"myplants/internal/database"
	"myplants/internal/model"
	"myplants/internal/response"
)

// ListPlants 植物列表
func ListPlants(c *gin.Context) {
	var plants []model.Plant
	query := database.DB.Order("created_at DESC")

	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}
	if status := c.Query("health_status"); status != "" {
		query = query.Where("health_status = ?", status)
	}
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ? OR species LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Find(&plants).Error; err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, plants)
}

// GetPlant 植物详情
func GetPlant(c *gin.Context) {
	id := c.Param("id")
	var plant model.Plant
	if err := database.DB.First(&plant, id).Error; err != nil {
		response.Fail(c, "植物不存在")
		return
	}
	response.OK(c, plant)
}

// CreatePlant 创建植物
func CreatePlant(c *gin.Context) {
	var req model.PlantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误: "+err.Error())
		return
	}
	plant := model.Plant{
		Name:             req.Name,
		Species:          req.Species,
		Category:         req.Category,
		Location:         req.Location,
		Avatar:           req.Avatar,
		AcquiredAt:       req.AcquiredAt,
		HealthStatus:     req.HealthStatus,
		LightRequirement: req.LightRequirement,
		WaterCycle:       req.WaterCycle,
		FertilizeCycle:   req.FertilizeCycle,
		SprayCycle:       req.SprayCycle,
		Price:            req.Price,
		PotID:            req.PotID,
		Description:      req.Description,
	}
	if plant.HealthStatus == "" {
		plant.HealthStatus = "长势良好"
	}
	if plant.WaterCycle == 0 {
		plant.WaterCycle = config.Get().Reminder.DefaultWaterDays
	}
	if plant.FertilizeCycle == 0 {
		plant.FertilizeCycle = config.Get().Reminder.DefaultFertilizeDays
	}
	if plant.SprayCycle == 0 {
		plant.SprayCycle = config.Get().Reminder.DefaultSprayDays
	}
	if err := database.DB.Create(&plant).Error; err != nil {
		response.Fail(c, "创建失败: "+err.Error())
		return
	}
	response.OKWithMsg(c, "创建成功", plant)
}

// UpdatePlant 更新植物
func UpdatePlant(c *gin.Context) {
	id := c.Param("id")
	var plant model.Plant
	if err := database.DB.First(&plant, id).Error; err != nil {
		response.Fail(c, "植物不存在")
		return
	}
	var req model.PlantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误: "+err.Error())
		return
	}
	plant.Name = req.Name
	plant.Species = req.Species
	plant.Category = req.Category
	plant.Location = req.Location
	plant.Avatar = req.Avatar
	plant.AcquiredAt = req.AcquiredAt
	plant.HealthStatus = req.HealthStatus
	plant.LightRequirement = req.LightRequirement
	plant.WaterCycle = req.WaterCycle
	plant.FertilizeCycle = req.FertilizeCycle
	plant.SprayCycle = req.SprayCycle
	plant.Price = req.Price
	plant.PotID = req.PotID
	plant.Description = req.Description
	if err := database.DB.Save(&plant).Error; err != nil {
		response.Fail(c, "更新失败: "+err.Error())
		return
	}
	response.OKWithMsg(c, "更新成功", plant)
}

// DeletePlant 删除植物
func DeletePlant(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.Plant{}, id).Error; err != nil {
			return err
		}
		tx.Where("plant_id = ?", id).Delete(&model.Photo{})
		tx.Where("plant_id = ?", id).Delete(&model.CareRecord{})
		tx.Where("plant_id = ?", id).Delete(&model.Repotting{})
		tx.Where("plant_id = ?", id).Delete(&model.Note{})
		return nil
	}); err != nil {
		response.Fail(c, "删除失败: "+err.Error())
		return
	}
	response.OKWithMsg(c, "删除成功", nil)
}

// GetReminders 养护提醒列表
func GetReminders(c *gin.Context) {
	var plants []model.Plant
	if err := database.DB.Find(&plants).Error; err != nil {
		response.Fail(c, "查询失败")
		return
	}

	now := time.Now()
	var reminders []model.ReminderItem
	for _, p := range plants {
		reminders = append(reminders, buildReminder(p, model.CareTypeWater, now))
		reminders = append(reminders, buildReminder(p, model.CareTypeFertilize, now))
		reminders = append(reminders, buildReminder(p, model.CareTypeSpray, now))
	}
	response.OK(c, reminders)
}

func buildReminder(p model.Plant, careType string, now time.Time) model.ReminderItem {
	var lastTime *time.Time
	var cycle int
	switch careType {
	case model.CareTypeWater:
		lastTime = p.LastWateredAt
		cycle = p.WaterCycle
	case model.CareTypeFertilize:
		lastTime = p.LastFertilizedAt
		cycle = p.FertilizeCycle
	case model.CareTypeSpray:
		lastTime = p.LastSprayedAt
		cycle = p.SprayCycle
	}

	item := model.ReminderItem{
		PlantID:   p.ID,
		PlantName: p.Name,
		Avatar:    p.Avatar,
		Type:      careType,
		LastTime:  lastTime,
		CycleDays: cycle,
	}

	if lastTime == nil {
		// 没记录过,视为需要立即处理
		item.NextTime = p.CreatedAt
		item.DaysLeft = int(now.Sub(p.CreatedAt).Hours() / 24)
		item.Overdue = true
	} else {
		next := lastTime.AddDate(0, 0, cycle)
		item.NextTime = next
		item.DaysLeft = int(next.Sub(now).Hours() / 24)
		if now.After(next) {
			item.Overdue = true
		}
	}
	return item
}

// GetPlantStats 单个植物统计
func GetPlantStats(c *gin.Context) {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)

	var waterCount, fertilizeCount, sprayCount, photoCount, repotCount, noteCount int64
	database.DB.Model(&model.CareRecord{}).Where("plant_id = ? AND type = ?", idInt, model.CareTypeWater).Count(&waterCount)
	database.DB.Model(&model.CareRecord{}).Where("plant_id = ? AND type = ?", idInt, model.CareTypeFertilize).Count(&fertilizeCount)
	database.DB.Model(&model.CareRecord{}).Where("plant_id = ? AND type = ?", idInt, model.CareTypeSpray).Count(&sprayCount)
	database.DB.Model(&model.Photo{}).Where("plant_id = ?", idInt).Count(&photoCount)
	database.DB.Model(&model.Repotting{}).Where("plant_id = ?", idInt).Count(&repotCount)
	database.DB.Model(&model.Note{}).Where("plant_id = ?", idInt).Count(&noteCount)

	daysKept := 0
	var plant model.Plant
	if err := database.DB.First(&plant, id).Error; err == nil {
		base := plant.CreatedAt
		if plant.AcquiredAt != nil {
			base = *plant.AcquiredAt
		}
		daysKept = int(time.Since(base).Hours() / 24)
	}

	response.OK(c, gin.H{
		"water_count":     waterCount,
		"fertilize_count": fertilizeCount,
		"spray_count":     sprayCount,
		"photo_count":     photoCount,
		"repot_count":     repotCount,
		"note_count":      noteCount,
		"days_kept":       daysKept,
	})
}
