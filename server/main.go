package main

import (
	"context"
	"cxfw/conf"
	"cxfw/db"
	"cxfw/model"
	"cxfw/model/fragments"
	"cxfw/model/todos"
	"cxfw/model/writer"
	"cxfw/orm"
	"cxfw/service"
	"strings"

	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
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
	initDB(pool)

	srv := gin.Default()
	// srv.StaticFile("/", "web/index.html")
	// srv.StaticFile("/favicon.ico", "web/favicon.ico")
	// srv.Static("/static", "web/static")

	service.Init(srv)

	if err := srv.RunTLS(fmt.Sprintf("%s:%d", conf.SrvHost, conf.SrvPort), conf.Cert, conf.PKey); err != nil {
		log.Fatal(err)
	}
}

func setupDB(db *pgxpool.Pool) {
	bt := pgx.Batch{}
	// system
	bt.Queue(model.UserSQL)
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

func initDB(db *pgxpool.Pool) {
	_ = orm.NewTx(db, func(tx pgx.Tx) error {
		bt := pgx.Batch{}

		// init data
		pass, _ := bcrypt.GenerateFromPassword([]byte("lcs1996"), bcrypt.DefaultCost)
		sb := strings.Builder{}
		sb.Write(pass)
		bt.Queue("INSERT INTO users (account, name, password) values ('lcs','lcs',$1) ON CONFLICT DO NOTHING;", sb.String())

		br := tx.SendBatch(context.Background(), &bt)
		err := br.Close()
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})
}
