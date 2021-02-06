package session

import "time"

type ISession interface {
	Set(k string, v interface{})
	Get(k string) interface{}
	Del(k string) bool
	ID() string
	StartTime() time.Time
	Update() bool
}
