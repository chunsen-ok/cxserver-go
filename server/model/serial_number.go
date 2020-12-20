package model

import "time"

type SerialNumber struct {
	SerialNumber uint64     `json:"serial_number" gorm:"column:serial_number;primaryKey;<-:false"`
	ExpireTime   *time.Time `json:"expire_time" gorm:"column:expire_time;type:timestamp"`
	Status       int        `json:"status" gorm:"column:status;default:1"`
}
