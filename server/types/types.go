package types

const (
	Active   = iota // 激活
	InActive        // 未激活
	InValid         // 失效
	Trash           // 删除
)

type Response struct {
	Err  *string     `json:"err"`
	Body interface{} `json:"body"`
}
