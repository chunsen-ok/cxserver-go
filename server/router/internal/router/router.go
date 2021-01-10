package router

import (
	"cxfw/types"
	"fmt"

	"github.com/gin-gonic/gin"
)

// RouteHandler ...
type RouteHandler = func(c *gin.Context) (int, interface{}, error)

func Route(h RouteHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		code, data, err := h(c)
		if err != nil {
			fmt.Printf("Route error: %v\n", err)
		}
		c.JSON(code, types.Response{Err: err, Body: data})
	}
}
