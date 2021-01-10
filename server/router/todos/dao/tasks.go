package dao

import (
	"context"
	"cxfw/model/todos"
	"net/http"
)

func (d *TodoDao) GetAllTask() (int, []todos.TodoTask, error) {
	rows, err := d.db.Query(context.Background(), `select * from todos.todo_tasks;`)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	ms := make([]todos.TodoTask, 0)
	for rows.Next() {
		var m todos.TodoTask
		err := rows.Scan(&m.ID, &m.Title, &m.CreatedAt, &m.Remark)
		if err != nil {
			rows.Close()
			return http.StatusInternalServerError, nil, err
		}

		ms = append(ms, m)
	}
	rows.Close()

	return http.StatusOK, ms, nil
}
