package model

import (
	"time"
)

type Post struct {
	ID        int       `json:"id" gorm:"column:id;primaryKey;<-:false"`
	Title     string    `json:"title" gorm:"column:title"`
	Content   string    `json:"content" gorm:"column:content"`
	Status    int       `json:"status" gorm:"column:status;default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamp;<-:create"` // UTC时间
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp"`           // UTC时间
}
