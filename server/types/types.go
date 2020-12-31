package types

// Status
const (
	StatusActive   = 0 // 激活
	StatusInActive = 1 // 未激活
	StatusInValid  = 2 // 失效
	StatusTrash    = 3 // 删除
)

// Badges
const (
	BadgeRank = 0 // 置顶排序
)

type Response struct {
	Err  error       `json:"err"`
	Body interface{} `json:"body"`
}
