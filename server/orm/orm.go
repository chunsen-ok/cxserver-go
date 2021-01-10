package orm

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// SqlIn .
// param: a1,a2,...
// return: "('a1', 'a2', ...)"
func SqlIn(arr interface{}) string {
	if arr == nil {
		return ""
	}

	sn := make([]string, 0)
	switch arr.(type) {
	case []int:
		nv, _ := arr.([]int)

		for _, v := range nv {
			sv := strconv.Itoa(v)
			sn = append(sn, fmt.Sprintf("'%s'", sv))
		}
	case []string:
		sv, _ := arr.([]string)
		for _, v := range sv {
			sn = append(sn, fmt.Sprintf("'%s'", v))
		}
	}

	if len(sn) > 0 {
		return fmt.Sprintf("(%s)", strings.Join(sn, ","))
	}

	return "(null)"
}

type TransactionFunc = func(tx pgx.Tx) error

// NewTx .
func NewTx(db *pgxpool.Pool, fn TransactionFunc) error {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	err = fn(tx)
	if err != nil {
		return err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}

	return nil
}
