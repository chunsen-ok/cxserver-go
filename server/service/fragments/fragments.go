package fragments

import (
	"cxfw/model/fragments"
	"cxfw/service/fragments/dao"
	"cxfw/service/internal/router"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Init(r gin.IRouter) {
	g := r.Group("/fragments")
	g.POST("/", router.Route(add))
	g.DELETE("/:id", router.Route(del))
	g.GET("/:id", router.Route(get))
	g.GET("/", router.Route(all))
}

func add(c *gin.Context) (int, interface{}, error) {
	m := fragments.Msg{}
	if err := c.ShouldBindJSON(&m); err != nil {
		return http.StatusBadRequest, nil, err
	}

	code, data, err := dao.Add(&m)

	return code, data, err
}

func del(c *gin.Context) (int, interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return http.StatusBadRequest, nil, errors.New("invalid id")
	}

	code, err := dao.Del(id)

	return code, nil, err
}

func get(c *gin.Context) (int, interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return http.StatusBadRequest, nil, errors.New("invalid id")
	}

	code, data, err := dao.Get(id)

	return code, data, err
}

func all(c *gin.Context) (int, interface{}, error) {
	code, data, err := dao.All()
	return code, data, err
}
