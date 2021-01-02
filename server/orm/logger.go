package orm

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

type Logger struct {
}

func (logger *Logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	for k, v := range data {
		fmt.Println(k, ": ", v)
	}
}
