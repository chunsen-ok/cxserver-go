package login

import (
	"cxfw/model"
	"cxfw/service/internal/router"
	"cxfw/service/login/dao"
	"net/http"

	"cxfw/session"

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

	code, data, err := dao.Login(&m)
	if code != http.StatusOK || err != nil {
		return code, data, err
	}

	se := session.S().GetSession(c)
	if se == nil {
		se = session.S().StartSession(c)
		if se == nil {
			return http.StatusInternalServerError, nil, nil
		}
	} else {
		se.Update()
	}

	se.Set("info", data)

	return code, data, err
}

func out(c *gin.Context) (int, interface{}, error) {
	session.S().StopSession(c)
	ok := dao.Logout(nil)
	return http.StatusOK, ok, nil
}
