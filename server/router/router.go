package router

import (
	"cxfw/conf"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Router struct {
	db *gorm.DB
}

func Init(db *gorm.DB) *Router {
	return &Router{
		db: db,
	}
}

func (r *Router) Routes(router gin.IRouter) {
	apiRouter := router.Group("/api", gin.BasicAuth(conf.Instance().BasicAuth))

	snRouter := apiRouter.Group("/sn")
	snRouter.POST("/", r.genSerialNumber)

	r.postsRoutes(apiRouter)
	r.tagsRoutes(apiRouter)
}
