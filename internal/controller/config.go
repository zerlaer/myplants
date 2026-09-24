package controller

import (
	"strings"

	"github.com/gin-gonic/gin"

	"myplants/internal/config"
	"myplants/internal/response"
)

// GetConfig 返回前端需要的配置项(来自 config.yaml)
func GetConfig(c *gin.Context) {
	cfg := config.Get()
	r := cfg.Reminder
	storageCfg := gin.H{"driver": cfg.Storage.Driver}
	if cfg.Storage.Driver == "r2" && cfg.Storage.R2.PublicBase != "" {
		base := strings.TrimSuffix(cfg.Storage.R2.PublicBase, "/")
		if !strings.Contains(base, "://") {
			base = "https://" + base // 配置漏写协议头时自动补全
		}
		storageCfg["base"] = base
	}
	response.OK(c, gin.H{
		"default_water_days":     r.DefaultWaterDays,
		"default_fertilize_days": r.DefaultFertilizeDays,
		"default_spray_days":     r.DefaultSprayDays,
		"storage":                storageCfg,
	})
}
