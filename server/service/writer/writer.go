package writer

import (
	"cxfw/service/writer/badges"
	"cxfw/service/writer/posts"
	"cxfw/service/writer/tags"

	"github.com/gin-gonic/gin"
)

func Init(r gin.IRouter) {
	g := r.Group("/writer")
	badges.Init(g)
	tags.Init(g)
	posts.Init(g)
}
