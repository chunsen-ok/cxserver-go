package service

import (
	"cxfw/conf"
	"cxfw/service/fragments"
	"cxfw/service/login"
	"cxfw/service/writer"

	"github.com/gin-gonic/gin"
)

// Init ...
func Init(r gin.IRouter) {
	api := r.Group("/api")

	login.Init(api)

	g := api.Use(gin.BasicAuth(conf.Instance().BasicAuth)).(gin.IRouter)
	writer.Init(g)
	// todos.Init(g)
	fragments.Init(g)
}
