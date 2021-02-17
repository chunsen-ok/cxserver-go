package login

import (
	"cxfw/model"
	"cxfw/service/internal/router"
	"cxfw/service/login/dao"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Init(r gin.IRouter) {
	g := r.Group("/login")
	g.POST("/", router.Route(in))
	g.DELETE("/", router.Route(out))
}

func in(c *gin.Context) (int, interface{}, error) {
	m := model.User{}
	if err := c.ShouldBindJSON(&m); err != nil {
		return http.StatusBadRequest, nil, err
	}

	code, ok, err := dao.Login(&m)

	return code, ok, err
}

func out(c *gin.Context) (int, interface{}, error) {
	ok := dao.Logout(nil)
	return http.StatusOK, ok, nil
}
