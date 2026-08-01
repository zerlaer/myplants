package controller

import (
	"github.com/gin-gonic/gin"

	"myplants/internal/config"
	"myplants/internal/response"
)

// GetConfig 返回前端需要的配置项(来自 config.yaml)
func GetConfig(c *gin.Context) {
	cfg := config.Get()
	r := cfg.Reminder
	response.OK(c, gin.H{
		"default_water_days":     r.DefaultWaterDays,
		"default_fertilize_days": r.DefaultFertilizeDays,
		"default_spray_days":     r.DefaultSprayDays,
	})
}
