package model

import "time"

// Pot 花盆
type Pot struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	PlantID   *uint     `json:"plant_id" gorm:"index"`               // 当前关联植物(可空)
	Name      string    `json:"name" gorm:"size:50"`                 // 名称
	Size      string    `json:"size" gorm:"size:30"`                 // 尺寸描述: 直径x高 或 加仑
	Diameter  float64   `json:"diameter"`                            // 直径cm
	Height    float64   `json:"height"`                              // 高度cm
	Gallon    float64   `json:"gallon"`                              // 加仑数(青山盆/加仑盆使用)
	Type      string    `json:"type" gorm:"column:material;size:20"` // 类型: 塑料盆/青山盆/透气盆/自吸盆/陶盆/加仑盆
	Generation string   `json:"generation" gorm:"size:10"`          // 代际: 一代/二代/三代(青山盆专用)
	Color     string    `json:"color" gorm:"size:20"`                // 颜色
	Status    string    `json:"status" gorm:"size:20;default:'使用中'"` // 使用中/空闲/已弃用
	Remark    string    `json:"remark" gorm:"size:255"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PotRequest 花盆请求
type PotRequest struct {
	PlantID    *uint   `json:"plant_id"`
	Name       string  `json:"name"`
	Size       string  `json:"size"`
	Diameter   float64 `json:"diameter"`
	Height     float64 `json:"height"`
	Gallon     float64 `json:"gallon"`
	Type       string  `json:"type"`
	Generation string  `json:"generation"`
	Color      string  `json:"color"`
	Status     string  `json:"status"`
	Remark     string  `json:"remark"`
}

// Repotting 换盆记录
type Repotting struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	PlantID    uint      `json:"plant_id" gorm:"index;not null"`
	FromPotID  *uint     `json:"from_pot_id"`  // 原花盆
	ToPotID    *uint     `json:"to_pot_id"`    // 新花盆
	FromPotName string   `json:"from_pot_name" gorm:"size:50"` // 冗余名称
	ToPotName   string   `json:"to_pot_name" gorm:"size:50"`   // 冗余名称
	RepotTime  time.Time `json:"repot_time"`   // 换盆时间
	Remark     string    `json:"remark" gorm:"size:255"`
	CreatedAt  time.Time `json:"created_at"`
}

// RepottingRequest 换盆请求
type RepottingRequest struct {
	PlantID   uint       `json:"plant_id" binding:"required"`
	FromPotID *uint      `json:"from_pot_id"`
	ToPotID   *uint      `json:"to_pot_id"`
	RepotTime *time.Time `json:"repot_time"`
	Remark    string     `json:"remark"`
}
