package todos

import "time"

type TodoItem struct {
	ID         int       `json:"id"`
	TaskID     int       `json:"task_id"`
	Title      string    `json:"title"`
	Content    *string   `json:"content"`
	Importance int       `json:"importance"`
	Urgency    int       `json:"urgency"`
	CreatedAt  time.Time `json:"created_at"`
	DeadLine   time.Time `json:"dead_line"`
}

const TodoItemsSQL = `
CREATE TABLE IF NOT EXISTS todos.todo_items (
	id serial primary key,
	task_id int not null,
	title text not null,
	content text null,
	importance int null,
	urgency int null,
	created_at timestamp(0) not null default now(),
	dead_line timestamp(0) null
);
`
