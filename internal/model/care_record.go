package model

import "time"

// CareType 养护类型
const (
	CareTypeWater    = "water"    // 浇水
	CareTypeFertilize = "fertilize" // 施肥
	CareTypeSpray    = "spray"    // 打药
)

// CareRecord 养护记录(浇水/施肥/打药)
type CareRecord struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	PlantID    uint      `json:"plant_id" gorm:"index;not null"`
	Type       string    `json:"type" gorm:"size:20;index;not null"` // water/fertilize/spray
	RecordTime time.Time `json:"record_time"`                        // 记录时间
	Remark     string    `json:"remark" gorm:"size:255"`             // 备注(肥料/药品名称、用量等)
	CreatedAt  time.Time `json:"created_at"`
}

// CareRequest 养护请求
type CareRequest struct {
	PlantID    uint      `json:"plant_id" binding:"required"`
	Type       string    `json:"type" binding:"required"`
	RecordTime *time.Time `json:"record_time"`
	Remark     string    `json:"remark"`
}

// ReminderItem 提醒项
type ReminderItem struct {
	PlantID      uint   `json:"plant_id"`
	PlantName    string `json:"plant_name"`
	Avatar       string `json:"avatar"`
	Type         string `json:"type"`          // water/fertilize/spray
	LastTime     *time.Time `json:"last_time"`
	NextTime     time.Time `json:"next_time"`
	CycleDays    int    `json:"cycle_days"`
	DaysLeft     int    `json:"days_left"`     // 距下次天数, 负数表示已逾期
	Overdue      bool   `json:"overdue"`       // 是否逾期
}
