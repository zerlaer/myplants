package model

import "time"

// Pot 花盆
type Pot struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	PlantID   *uint     `json:"plant_id" gorm:"index"`               // 当前关联植物(可空)
	Name      string    `json:"name" gorm:"size:50"`                 // 名称
	Size      string    `json:"size" gorm:"size:30"`                 // 尺寸描述: 直径x高
	Diameter  float64   `json:"diameter"`                            // 直径cm
	Height    float64   `json:"height"`                              // 高度cm
	Material  string    `json:"material" gorm:"size:20"`             // 材质: 塑料/陶土/陶瓷/水泥/木质/其他
	Color     string    `json:"color" gorm:"size:20"`                // 颜色
	Status    string    `json:"status" gorm:"size:20;default:'使用中'"` // 使用中/空闲/已弃用
	Remark    string    `json:"remark" gorm:"size:255"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PotRequest 花盆请求
type PotRequest struct {
	PlantID  *uint   `json:"plant_id"`
	Name     string  `json:"name"`
	Size     string  `json:"size"`
	Diameter float64 `json:"diameter"`
	Height   float64 `json:"height"`
	Material string  `json:"material"`
	Color    string  `json:"color"`
	Status   string  `json:"status"`
	Remark   string  `json:"remark"`
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
