package model

import "time"

type Tag struct {
	ID        int       `json:"id" gorm:"column:id;primaryKey;<-:false"`
	Title     string    `json:"title" gorm:"column:title;"`
	CreatedAt time.Time `json:"create_time" gorm:"column:created_at;type:timestamp"`
}
