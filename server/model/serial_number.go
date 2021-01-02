package model

import "time"

type SerialNumber struct {
	SerialNumber int        `json:"serial_number"`
	ExpireTime   *time.Time `json:"expire_time"`
	Status       int16      `json:"status"`
}
