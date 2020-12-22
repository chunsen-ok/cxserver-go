package router

import (
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
	apiRouter := router.Group("/api")

	snRouter := apiRouter.Group("/sn")
	snRouter.POST("/", r.genSerialNumber)

	r.postsRoutes(apiRouter)

	tagRouter := apiRouter.Group("/tags")
	tagRouter.POST("/", r.newTag)
	tagRouter.DELETE("/:id", r.delTag)
	tagRouter.GET("/", r.getTags)
	tagRouter.GET("/:id", r.getTag)
	tagRouter.PUT("/", r.updateTag)
}
