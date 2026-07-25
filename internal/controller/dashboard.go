package controller

import (
	"time"

	"github.com/gin-gonic/gin"

	"myplants/internal/database"
	"myplants/internal/model"
	"myplants/internal/response"
)

// GetDashboard 仪表盘统计
func GetDashboard(c *gin.Context) {
	var plantCount, potCount, photoCount, waterCount, fertilizeCount, sprayCount int64
	database.DB.Model(&model.Plant{}).Count(&plantCount)
	database.DB.Model(&model.Pot{}).Count(&potCount)
	database.DB.Model(&model.Photo{}).Count(&photoCount)
	database.DB.Model(&model.CareRecord{}).Where("type = ?", model.CareTypeWater).Count(&waterCount)
	database.DB.Model(&model.CareRecord{}).Where("type = ?", model.CareTypeFertilize).Count(&fertilizeCount)
	database.DB.Model(&model.CareRecord{}).Where("type = ?", model.CareTypeSpray).Count(&sprayCount)

	// 今日养护统计
	today := time.Now().Format("2006-01-02")
	var todayWater, todayFertilize, todaySpray int64
	database.DB.Model(&model.CareRecord{}).Where("type = ? AND date(record_time) = ?", model.CareTypeWater, today).Count(&todayWater)
	database.DB.Model(&model.CareRecord{}).Where("type = ? AND date(record_time) = ?", model.CareTypeFertilize, today).Count(&todayFertilize)
	database.DB.Model(&model.CareRecord{}).Where("type = ? AND date(record_time) = ?", model.CareTypeSpray, today).Count(&todaySpray)

	// 健康状态分布
	type statusStat struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	var statusStats []statusStat
	database.DB.Model(&model.Plant{}).Select("health_status as status, count(*) as count").Group("health_status").Scan(&statusStats)

	// 分类分布
	var categoryStats []statusStat
	database.DB.Model(&model.Plant{}).Select("category as status, count(*) as count").Group("category").Scan(&categoryStats)

	// 近7天养护趋势
	type dayStat struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	var trend []dayStat
	for i := 6; i >= 0; i-- {
		d := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var count int64
		database.DB.Model(&model.CareRecord{}).Where("date(record_time) = ?", d).Count(&count)
		trend = append(trend, dayStat{Date: d, Count: count})
	}

	// 需要关注的植物(逾期未浇水)
	now := time.Now()
	var plants []model.Plant
	database.DB.Find(&plants)
	var overdueCount int64
	var totalPrice float64
	for _, p := range plants {
		if p.LastWateredAt == nil {
			overdueCount++
		} else if now.After(p.LastWateredAt.AddDate(0, 0, p.WaterCycle)) {
			overdueCount++
		}
		totalPrice += p.Price
	}

	// 分类价格统计
	type categoryPriceStat struct {
		Category string  `json:"category"`
		Count    int64   `json:"count"`
		TotalPrice float64 `json:"total_price"`
	}
	var categoryPriceStats []categoryPriceStat
	database.DB.Model(&model.Plant{}).Select("category, count(*) as count, sum(price) as total_price").Group("category").Scan(&categoryPriceStats)

	response.OK(c, gin.H{
		"plant_count":         plantCount,
		"pot_count":           potCount,
		"photo_count":         photoCount,
		"water_count":         waterCount,
		"fertilize_count":     fertilizeCount,
		"spray_count":         sprayCount,
		"today_water":         todayWater,
		"today_fertilize":     todayFertilize,
		"today_spray":         todaySpray,
		"status_stats":        statusStats,
		"category_stats":      categoryStats,
		"category_price_stats": categoryPriceStats,
		"trend":               trend,
		"overdue_count":       overdueCount,
		"total_price":         totalPrice,
		"plants":              plants,
	})
}
