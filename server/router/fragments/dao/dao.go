package dao

import (
	"context"
	"cxfw/db"
	"cxfw/model/fragments"
	"cxfw/orm"
	"net/http"

	"github.com/jackc/pgx/v4"
)

// url: [POST] /api/fragments/
// param: m fragments.Msg{}
// response: fragments.Msg{}
func Add(m *fragments.Msg) (int, *fragments.Msg, error) {
	err := orm.NewTx(db.S(), func(tx pgx.Tx) error {
		err := tx.QueryRow(context.Background(),
			`insert into fragments.msgs (content, created_at) values ($1, now() at time zone 'utc') returning id, created_at;`,
			m.Content).Scan(&m.ID, &m.CreatedAt)
		return err
	})

	code := http.StatusOK
	if err != nil {
		code = http.StatusInternalServerError
	}

	return code, m, err
}

// url: [DELETE] /api/fragments/:id
// param: id int
// response: null
func Del(id int) (int, error) {
	err := orm.NewTx(db.S(), func(tx pgx.Tx) error {
		_, err := tx.Exec(context.Background(), `delete from fragments.msgs where id = $1;`, id)
		return err
	})

	code := http.StatusOK
	if err != nil {
		code = http.StatusInternalServerError
	}

	return code, err
}

// url: [GET] /api/fragments/:id
// param: id int
// response: fragments.Msg{}
func Get(id int) (int, *fragments.Msg, error) {
	m := fragments.Msg{}
	err := db.S().QueryRow(context.Background(),
		`select * from fragments.msgs where id = $1;`, id).
		Scan(&m.ID, &m.Content, &m.CreatedAt)

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, &m, nil
}

// url: [GET] /api/fragments
// response: []fragments.Msg
func All() (int, []fragments.Msg, error) {
	rows, err := db.S().Query(context.Background(), `select * from fragments.msgs;`)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	ms := make([]fragments.Msg, 0)
	for rows.Next() {
		m := fragments.Msg{}
		err := rows.Scan(&m.ID, &m.Content, &m.CreatedAt)
		if err != nil {
			rows.Close()
			return http.StatusInternalServerError, nil, err
		}

		ms = append(ms, m)
	}

	return http.StatusOK, ms, nil
}
