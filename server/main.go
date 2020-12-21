//Package main
//
// ## http 会话管理的几种方式。
//
// ### 基于服务器的会话管理
//
// 在服务端生成并管理 session id 。客户端访问时，需上传该 session id 。
// 服务端根据客户端的 session id ，去查找用户登录状态等，判断该用户是否登录了。
// 如果已经登录则可进行数据访问等其他操作。
//
// ### 基于 Cookie 的会话管理
//
// 服务器利用加密算法，将登录成功后的用户信息加密以及生成摘要等处理后，返回给客户端，
// 由客户端保存在 Cookie 中，而服务器不需要保存。每次客户端的访问，都需要带上这些信息，
// 服务端尽心解密，验证。
//
// ### 基于 token 的会话管理
//
// 关于证书：https://www.cnblogs.com/kyrios/p/tls-and-certificates.html
//
package main

import (
	"cxfw/model"
	"cxfw/router"
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/pelletier/go-toml"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	confPath string
)

func init() {
	flag.StringVar(&confPath, "c", "conf.toml", "configure file path.")
}

func main() {
	flag.Parse()

	conf, err := loadConf(confPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(conf.DatabaseUrl())

	db, err := gorm.Open(postgres.Open(conf.DatabaseUrl()),
		&gorm.Config{
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
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
	srv.StaticFile("/", "web/index.html")
	srv.StaticFile("/favicon.ico", "web/favicon.ico")
	srv.Static("/static", "web/static")

	router := router.Init(db)
	router.Routes(srv)

	if err := srv.RunTLS(fmt.Sprintf("%s:%d", conf.SrvHost, conf.SrvPort), conf.Cert, conf.PKey); err != nil {
		log.Fatal(err)
	}
}

type Conf struct {
	SrvHost string
	SrvPort int64

	User     string
	password string
	Database string
	Host     string
	Port     int64

	Cert string
	PKey string
}

func (c Conf) DatabaseUrl() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", c.User, c.password, c.Host, c.Port, c.Database)
}

func loadConf(path string) (*Conf, error) {
	tree, err := toml.LoadFile(path)
	if err != nil {
		return nil, err
	}

	server, ok := tree.Get("server").(*toml.Tree)
	if !ok {
		return nil, errors.New("server conf error")
	}

	database, ok := tree.Get("database").(*toml.Tree)
	if !ok {
		return nil, errors.New("database conf error")
	}

	https, ok := tree.Get("https").(*toml.Tree)
	if !ok {
		return nil, errors.New("https conf error")
	}

	conf := Conf{
		SrvHost:  server.Get("host").(string),
		SrvPort:  server.Get("port").(int64),
		User:     database.Get("user").(string),
		password: database.Get("password").(string),
		Database: database.Get("db").(string),
		Host:     database.Get("host").(string),
		Port:     database.Get("port").(int64),
		Cert:     https.Get("cert").(string),
		PKey:     https.Get("pkey").(string),
	}

	return &conf, nil
}
