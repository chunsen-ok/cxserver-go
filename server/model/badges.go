package model

type PostBadge struct {
	BadgeName  int     `json:"badge_name" gorm:"column:badge_name"`
	BadgeValue *string `json:"badge_value" gorm:"column:badge_value"`
	PostID     int     `json:"post_id" gorm:"column:post_id"`
}
