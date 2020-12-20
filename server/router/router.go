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

	postRouter := apiRouter.Group("/posts")
	postRouter.POST("/", r.newPost)
	postRouter.DELETE("/:id", r.delPost)
	postRouter.GET("/", r.getPosts)
	postRouter.GET("/:id", r.getPost)
	postRouter.PUT("/", r.updatePost)

	tagRouter := apiRouter.Group("/tags")
	tagRouter.POST("/", r.newTag)
	tagRouter.DELETE("/:id", r.delTag)
	tagRouter.GET("/", r.getTags)
	tagRouter.GET("/:id", r.getTag)
	tagRouter.PUT("/", r.updateTag)
}
