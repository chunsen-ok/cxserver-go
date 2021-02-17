package sys

import (
	"github.com/gin-gonic/gin"
)

func Init(r gin.IRouter) {
	_ = r.Group("/sys")
}
