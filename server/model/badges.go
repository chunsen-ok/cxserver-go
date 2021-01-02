package model

type PostBadge struct {
	BadgeName  int     `json:"badge_name"`
	BadgeValue *string `json:"badge_value"`
	PostID     int     `json:"post_id"`
}

const postBadgeSQL = `
CREATE TABLE IF NOT EXISTS post_badges (
	badge_name integer NOT NULL,
	badge_value varchar(200) NULL,
	post_id integer NOT NULL
);`
