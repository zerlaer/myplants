package model

import "time"

// Plant 植物
type Plant struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	Name            string     `json:"name" gorm:"size:50;not null"`             // 名称
	Species         string     `json:"species" gorm:"size:100"`                  // 品种/学名
	Category        string     `json:"category" gorm:"size:30;index"`            // 分类: 绿植/多肉/花卉/草本/木本/其他
	Location        string     `json:"location" gorm:"size:100"`                 // 摆放位置
	Avatar          string     `json:"avatar" gorm:"size:255"`                   // 头像路径
	AcquiredAt      *time.Time `json:"acquired_at"`                              // 获得日期
	HealthStatus    string     `json:"health_status" gorm:"size:20;default:'长势良好'"` // 长势良好/正在缓苗/生长缓慢/状态不佳/生病枯萎/含苞待放/已经开花/已经结果
	LightRequirement string    `json:"light_requirement" gorm:"size:20"`         // 喜阳/半阴/喜阴
	WaterCycle      int        `json:"water_cycle" gorm:"default:7"`             // 浇水周期(天)
	FertilizeCycle  int        `json:"fertilize_cycle" gorm:"default:30"`        // 施肥周期(天)
	SprayCycle      int        `json:"spray_cycle" gorm:"default:45"`            // 打药周期(天)
	Price           float64    `json:"price" gorm:"default:0"`                   // 购买价格
	PotID           *uint      `json:"pot_id"`                                   // 关联花盆ID
	Description     string     `json:"description" gorm:"type:text"`            // 描述
	LastWateredAt   *time.Time `json:"last_watered_at"`                          // 最后浇水时间
	LastFertilizedAt *time.Time `json:"last_fertilized_at"`                     // 最后施肥时间
	LastSprayedAt   *time.Time `json:"last_sprayed_at"`                         // 最后打药时间
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// PlantRequest 创建/更新植物请求
type PlantRequest struct {
	Name             string     `json:"name" binding:"required"`
	Species          string     `json:"species"`
	Category         string     `json:"category"`
	Location         string     `json:"location"`
	Avatar           string     `json:"avatar"`
	AcquiredAt       *time.Time `json:"acquired_at"`
	HealthStatus     string     `json:"health_status"`
	LightRequirement string     `json:"light_requirement"`
	WaterCycle       int        `json:"water_cycle"`
	FertilizeCycle   int        `json:"fertilize_cycle"`
	SprayCycle       int        `json:"spray_cycle"`
	Price            float64    `json:"price"`
	PotID            *uint      `json:"pot_id"`
	Description      string     `json:"description"`
}
