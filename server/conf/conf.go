package conf

import (
	"errors"
	"fmt"

	"github.com/pelletier/go-toml"
)

// Conf ...
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

// DatabaseURL ...
func (c Conf) DatabaseURL() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", c.User, c.password, c.Host, c.Port, c.Database)
}

// LoadConf ...
func LoadConf(path string) (*Conf, error) {
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
