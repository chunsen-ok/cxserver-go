package model

import "time"

type Tag struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"create_time"`
}

const tagSQL = `
CREATE TABLE IF NOT EXISTS tags (
	id serial primary key,
	title text not null,
	created_at timestamp(0) not null default now()
);`
