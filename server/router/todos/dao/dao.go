package dao

import (
	"context"
	"cxfw/model/todos"
	"cxfw/orm"
	"cxfw/types"
	"fmt"
	"net/http"
	"strings"
	"time"

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

type NewTodoItemParam struct {
	ID         int        `json:"id"`
	TaskID     int        `json:"task_id"`
	Title      string     `json:"title"`
	Content    *string    `json:"content"`
	Importance int        `json:"importance"`
	Urgency    int        `json:"urgency"`
	DeadLine   *time.Time `json:"dead_line"`
}

func (d *TodoDao) New(p *NewTodoItemParam) (int, *todos.TodoItem, error) {
	var m todos.TodoItem
	err := orm.NewTx(d.db, func(tx pgx.Tx) error {
		err := tx.QueryRow(context.Background(),
			`insert into todos.todo_items values (default, $1,$2,$3,$4,$5,now() at time zone 'utc',$6 at time zone 'utc', default) returning *;`,
			p.TaskID, p.Title, p.Content, p.Importance, p.Urgency, p.DeadLine).
			Scan(&m.ID, &m.TaskID, &m.Title, &m.Content, &m.Importance, &m.Urgency, &m.CreatedAt, &m.DeadLine, &m.Status)
		return err
	})

	return http.StatusOK, &m, err
}

func (d *TodoDao) Del(itemID int) (int, error) {
	err := orm.NewTx(d.db, func(tx pgx.Tx) error {
		_, err := tx.Exec(context.Background(), `update todos.todo_items set status = $1 where id = $2;`, types.StatusTrash, itemID)
		return err
	})

	return http.StatusOK, err
}

func (d *TodoDao) GetAll(dimen, taskID int) (int, []todos.TodoItem, error) {
	sql := `with t as (select *, case "%s" when 0 then 1234 else "%s" end as rk from todos.todo_items where status <> $1) select * from t `
	sb := new(strings.Builder)
	rank := true
	if dimen == DimenImportance {
		sb.WriteString(fmt.Sprintf(sql, "importance", "importance"))
	} else if dimen == DimenUrgency {
		sb.WriteString(fmt.Sprintf(sql, "urgency", "urgency"))
	} else if dimen == DimenDeadLine {
		sb.WriteString(fmt.Sprintf(sql, "dead_line", "dead_line"))
	} else {
		sb.WriteString(`select *, '' from todos.todo_items`)
		rank = false
	}

	if taskID > 0 {
		sb.WriteString(fmt.Sprintf(`where task_id = %d `, taskID))
	}

	if rank {
		sb.WriteString(`order by rk asc `)
	}

	sb.WriteString(";")

	rows, err := d.db.Query(context.Background(), sb.String(), types.StatusTrash)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	ms := make([]todos.TodoItem, 0)
	for rows.Next() {
		var m todos.TodoItem
		err := rows.Scan(&m.ID, &m.TaskID, &m.Title, &m.Content, &m.Importance, &m.Urgency, &m.CreatedAt, &m.DeadLine, &m.Status, nil)
		if err != nil {
			rows.Close()
			return http.StatusInternalServerError, nil, err
		}

		ms = append(ms, m)
	}
	rows.Close()

	return http.StatusOK, ms, nil
}
