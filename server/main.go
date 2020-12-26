package main

import (
	"cxfw/conf"
	"cxfw/model"
	"cxfw/router"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	db, err := gorm.Open(postgres.Open(conf.DatabaseURL()),
		&gorm.Config{
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
			Logger: logger.Default.LogMode(logger.Info),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&model.SerialNumber{})
	db.AutoMigrate(&model.Post{})
	db.AutoMigrate(&model.Tag{})
	db.AutoMigrate(&model.PostTag{})

	srv := gin.Default()
	// srv.StaticFile("/", "web/index.html")
	// srv.StaticFile("/favicon.ico", "web/favicon.ico")
	// srv.Static("/static", "web/static")

	router := router.Init(db)
	router.Routes(srv)

	if err := srv.RunTLS(fmt.Sprintf("%s:%d", conf.SrvHost, conf.SrvPort), conf.Cert, conf.PKey); err != nil {
		log.Fatal(err)
	}
}
