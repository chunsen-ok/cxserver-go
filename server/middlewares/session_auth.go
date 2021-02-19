package middlewares

import (
	"cxfw/session"

	"github.com/gin-gonic/gin"
)

func SessionAuth(c *gin.Context) {
	se := session.S().GetSession(c)
	if se == nil {
		c.Abort()
	} else {
		se.Update()
		c.Next()
	}
}
