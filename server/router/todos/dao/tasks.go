package dao

import (
	"context"
	"cxfw/model/todos"
	"cxfw/orm"
	"cxfw/types"
	"net/http"

	"github.com/jackc/pgx/v4"
)

type NewTodoTaskParam struct {
	Title  string  `json:"title"`
	Remark *string `json:"remark"`
}

func (d *TodoDao) NewTask(p *NewTodoTaskParam) (int, *todos.TodoTask, error) {
	var m todos.TodoTask
	err := orm.NewTx(d.db, func(tx pgx.Tx) error {
		err := tx.QueryRow(context.Background(),
			`insert into todos.todo_tasks values (default, $1, now() at time zone 'utc', $2, $3) returning *;`,
			p.Title, p.Remark, types.StatusActive).
			Scan(&m.ID, &m.Title, &m.CreatedAt, &m.Remark, &m.Status)
		return err
	})
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, &m, nil
}

func (d *TodoDao) DelTask(id int) (int, error) {
	err := orm.NewTx(d.db, func(tx pgx.Tx) error {
		_, err := tx.Exec(context.Background(),
			`update todos.todo_tasks set status = $1 where id = $2`, types.StatusTrash, id)
		return err
	})
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (d *TodoDao) GetAllTask() (int, []todos.TodoTask, error) {
	rows, err := d.db.Query(context.Background(), `select * from todos.todo_tasks where status <> $1;`, types.StatusTrash)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	ms := make([]todos.TodoTask, 0)
	for rows.Next() {
		var m todos.TodoTask
		err := rows.Scan(&m.ID, &m.Title, &m.CreatedAt, &m.Remark, &m.Status)
		if err != nil {
			rows.Close()
			return http.StatusInternalServerError, nil, err
		}

		ms = append(ms, m)
	}
	rows.Close()

	return http.StatusOK, ms, nil
}
