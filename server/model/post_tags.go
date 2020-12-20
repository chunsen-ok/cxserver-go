package model

type PostTags struct {
	ID     int `json:"id"`
	PostID int `json:"post_id"`
	TagID  int `json:"tag_id"`
}
