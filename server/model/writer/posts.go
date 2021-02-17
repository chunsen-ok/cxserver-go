package writer

import (
	"time"
)

type Post struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Status    int16     `json:"status"`
	CreatedAt time.Time `json:"created_at"` // UTC时间
	UpdatedAt time.Time `json:"updated_at"` // UTC时间
}

const PostSQL = `
CREATE TABLE IF NOT EXISTS writer.posts (
	id serial primary key,
	content text null,
	status smallint not null,
	created_at timestamp(0) not null default now(),
	updated_at timestamp(0) not null default now()
);
`
