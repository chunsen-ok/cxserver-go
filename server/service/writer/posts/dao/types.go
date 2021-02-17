package dao

import "cxfw/model/writer"

type PostWithTagIDs struct {
	writer.Post
	Title  string `json:"title"`
	TagIDs []int  `json:"tags"`
}

type PostWithTags struct {
	writer.Post
	Title string       `json:"title"`
	Tags  []writer.Tag `json:"tags"`
}

type PostBadge struct {
	BadgeName  int     `json:"badge_name"`
	BadgeValue *string `json:"badge_value"`
}

type PostWithBadges struct {
	writer.Post
	Title  string      `json:"title"`
	Badges []PostBadge `json:"badges"`
}
