package sys

import (
	"cxfw/service/sys/fs"

	"github.com/gin-gonic/gin"
)

func Init(r gin.IRouter) {
	g := r.Group("/sys")

	fs.Init(g)
}
