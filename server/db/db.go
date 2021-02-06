package db

import (
	"context"
	"cxfw/orm"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var db *pgxpool.Pool

func S() *pgxpool.Pool {
	return db
}

func Init(url string) *pgxpool.Pool {

	connConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Panic(err)
	}

	connConfig.ConnConfig.LogLevel = pgx.LogLevelInfo
	connConfig.ConnConfig.Logger = new(orm.Logger)

	dbPool, err := pgxpool.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		log.Fatal(err)
	}

	db = dbPool

	return db
}
