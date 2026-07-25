package model

import "time"

// Photo 相册照片
type Photo struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	PlantID   uint      `json:"plant_id" gorm:"index;not null"`
	Path      string    `json:"path" gorm:"size:255;not null"` // 文件相对路径
	Remark    string    `json:"remark" gorm:"size:255"`        // 备注
	TakenAt   time.Time `json:"taken_at"`                      // 记录时间
	CreatedAt time.Time `json:"created_at"`
}

// PhotoRequest 上传/更新照片请求
type PhotoRequest struct {
	PlantID uint   `json:"plant_id" binding:"required"`
	Remark  string `json:"remark"`
	TakenAt *time.Time `json:"taken_at"`
}
