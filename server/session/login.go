package session

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	c.String(http.StatusOK, "login ok")
}

func Logout(c *gin.Context) {
	c.String(http.StatusOK, "logout ok")
}
