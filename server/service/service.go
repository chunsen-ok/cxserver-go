package service

import (
	"cxfw/conf"
	"cxfw/middlewares"
	"cxfw/service/fragments"
	"cxfw/service/login"
	"cxfw/service/sys"
	"cxfw/service/writer"

	"github.com/gin-gonic/gin"
)

// Init ...
func Init(r gin.IRouter) {
	api := r.Group("/api")

	login.Init(api)

	g := api.
		Use(gin.BasicAuth(conf.Instance().BasicAuth)).
		Use(middlewares.SessionAuth).(gin.IRouter)
	sys.Init(g)
	writer.Init(g)
	// todos.Init(g)
	fragments.Init(g)
}
