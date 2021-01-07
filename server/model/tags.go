package model

import "time"

type Tag struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Parent    int       `json:"parent"`
	CreatedAt time.Time `json:"created_at"`
}

const tagSQL = `
CREATE TABLE IF NOT EXISTS tags (
	id serial primary key,
	title text not null,
	parent integer default -1,
	created_at timestamp(0) not null default now()
);`
