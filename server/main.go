package main

import (
	"context"
	"cxfw/conf"
	"cxfw/orm"
	"cxfw/router"
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
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

	// db, err := gorm.Open(postgres.Open(conf.DatabaseURL()),
	// 	&gorm.Config{
	// 		NowFunc: func() time.Time {
	// 			return time.Now().UTC()
	// 		},
	// 		Logger: logger.Default.LogMode(logger.Info),
	// 	},
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// db.AutoMigrate(&model.SerialNumber{})
	// db.AutoMigrate(&model.Post{})
	// db.AutoMigrate(&model.Tag{})
	// db.AutoMigrate(&model.PostTag{})
	// db.AutoMigrate(&model.PostBadge{})
	connConfig, err := pgxpool.ParseConfig(conf.DatabaseURL())
	if err != nil {
		log.Fatal(err)
	}

	connConfig.ConnConfig.LogLevel = pgx.LogLevelInfo
	connConfig.ConnConfig.Logger = new(orm.Logger)

	dbPool, err := pgxpool.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		log.Fatal(err)
	}

	srv := gin.Default()
	// srv.StaticFile("/", "web/index.html")
	// srv.StaticFile("/favicon.ico", "web/favicon.ico")
	// srv.Static("/static", "web/static")

	router := router.Init(dbPool)
	router.Routes(srv)

	if err := srv.RunTLS(fmt.Sprintf("%s:%d", conf.SrvHost, conf.SrvPort), conf.Cert, conf.PKey); err != nil {
		log.Fatal(err)
	}
}
