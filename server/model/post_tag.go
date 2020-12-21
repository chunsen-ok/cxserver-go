package model

type PostTag struct {
	PostID int `json:"post_id" gorm:"column:post_id;not null"`
	TagID int `json:"tag_id" gorm:"column:tag_id;not null"`
}
