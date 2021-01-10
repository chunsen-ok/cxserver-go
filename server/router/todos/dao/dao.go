package dao

import (
	"context"
	"cxfw/model/todos"
	"cxfw/orm"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Dimen
const (
	DimenImportance = 0
	DimenUrgency    = 1
	DimenDeadLine   = 2
)

type TodoDao struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *TodoDao {
	return &TodoDao{
		db: db,
	}
}

func (d *TodoDao) New(m *todos.TodoItem) (int, *todos.TodoItem, error) {
	err := orm.NewTx(d.db, func(tx pgx.Tx) error {
		err := tx.QueryRow(context.Background(),
			`insert into todos.todo_items values (default, $1,$2,$3,$4,$5,now() at time zone 'utc',$6) returning *;`).
			Scan(&m.ID, nil, nil, nil, nil, nil, &m.CreatedAt, &m.DeadLine)
		return err
	})

	return http.StatusOK, m, err
}

func (d *TodoDao) Del(itemID int) (int, error) {
	err := orm.NewTx(d.db, func(tx pgx.Tx) error {
		_, err := tx.Exec(context.Background(), `delete from todos.todo_items where id = $1;`, itemID)
		return err
	})

	return http.StatusOK, err
}

func (d *TodoDao) GetAll(dimen, taskID int) (int, []todos.TodoItem, error) {
	sb := new(strings.Builder)
	sb.WriteString(`select * from todos.todo_items `)
	if taskID > 0 {
		sb.WriteString(fmt.Sprintf(`where task_id = %d `, taskID))
	}
	if dimen == DimenImportance {
		sb.WriteString(`order by importance asc`)
	} else if dimen == DimenUrgency {
		sb.WriteString(`order by urgency asc`)
	} else if dimen == DimenDeadLine {
		sb.WriteString(`order by dead_line asc`)
	}
	sb.WriteString(";")

	rows, err := d.db.Query(context.Background(), sb.String())
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	ms := make([]todos.TodoItem, 0)
	for rows.Next() {
		var m todos.TodoItem
		err := rows.Scan(&m.ID, &m.TaskID, &m.Title, &m.Content, &m.Importance, &m.Urgency, &m.CreatedAt, &m.DeadLine)
		if err != nil {
			rows.Close()
			return http.StatusInternalServerError, nil, err
		}

		ms = append(ms, m)
	}
	rows.Close()

	return http.StatusOK, ms, nil
}
