package main

import (
	"context"
	"cxfw/conf"
	"cxfw/db"
	"cxfw/model/fragments"
	"cxfw/model/todos"
	"cxfw/model/writer"
	"cxfw/router"

	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	confPath string
)

func init() {
	flag.StringVar(&confPath, "c", "conf.toml", "configure file path.")
}

func main() {
	flag.Parse()

	conf, err := conf.LoadConf(confPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(conf.DatabaseURL())

	pool := db.Init(conf.DatabaseURL())
	setupDB(pool)

	srv := gin.Default()
	// srv.StaticFile("/", "web/index.html")
	// srv.StaticFile("/favicon.ico", "web/favicon.ico")
	// srv.Static("/static", "web/static")

	router := router.New(pool)
	router.Routes(srv)

	if err := srv.RunTLS(fmt.Sprintf("%s:%d", conf.SrvHost, conf.SrvPort), conf.Cert, conf.PKey); err != nil {
		log.Fatal(err)
	}
}

func setupDB(db *pgxpool.Pool) {
	bt := pgx.Batch{}
	// writer
	bt.Queue("CREATE SCHEMA IF NOT EXISTS writer;")
	bt.Queue(writer.PostBadgeSQL)
	bt.Queue(writer.PostTagSQL)
	bt.Queue(writer.PostSQL)
	bt.Queue(writer.TagSQL)
	// todos
	bt.Queue("CREATE SCHEMA IF NOT EXISTS todos;")
	bt.Queue(todos.TodoTasksSQL)
	bt.Queue(todos.TodoItemsSQL)
	// fragments
	bt.Queue("CREATE SCHEMA IF NOT EXISTS fragments;")
	bt.Queue(fragments.MsgSQL)

	tx, err := db.Begin(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback(context.Background())

	br := tx.SendBatch(context.Background(), &bt)
	err = br.Close()
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit(context.Background())
}
