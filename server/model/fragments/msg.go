package fragments

import "time"

type Msg struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

const MsgSQL = `
CREATE TABLE IF NOT EXISTS fragments.msgs (
	id serial primary key,
	content text not null,
	created_at timestamp(0) not null default now()
);
`
