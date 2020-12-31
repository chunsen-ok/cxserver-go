package router

import (
	"cxfw/conf"
	"cxfw/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RouteHandler ...
type RouteHandler = func(c *gin.Context) (int, interface{}, error)

func route(h RouteHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		code, data, err := h(c)
		c.JSON(code, types.Response{Err: err, Body: data})
	}
}

// Router ...
type Router struct {
	db *gorm.DB
}

// Init ...
func Init(db *gorm.DB) *Router {
	return &Router{
		db: db,
	}
}

// Routes .
func (r *Router) Routes(router gin.IRouter) {
	apiRouter := router.Group("/api", gin.BasicAuth(conf.Instance().BasicAuth))

	r.postsRoutes(apiRouter)
	r.tagsRoutes(apiRouter)
	r.badgesRoutes(apiRouter)
	// r.snRoutes(apiRouter)
}
