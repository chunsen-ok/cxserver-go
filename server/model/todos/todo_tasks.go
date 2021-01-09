package todos

import "time"

type TodoTask struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Remark    *string   `json:"remark"`
}

const TodoTasksSQL = `
CREATE TABLE IF NOT EXISTS todos.todo_tasks (
	id serial PRIMARY KEY,
	title text NOT NULL,
	created_at timestamp(0) NOT NULL DEFAULT now(),
	remark text NULL
);
`
