package model

import "time"

// Note 植物日记/笔记
type Note struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	PlantID   uint      `json:"plant_id" gorm:"index;not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at"`
}

// NoteRequest 笔记请求
type NoteRequest struct {
	PlantID uint   `json:"plant_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}
