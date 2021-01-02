package router

import (
	"context"
	"cxfw/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Router) snRoutes(g gin.IRouter) {
	group := g.Group("/sn")
	group.POST("/", route(r.genSerialNumber))
}

func (r *Router) genSerialNumber(c *gin.Context) (int, interface{}, error) {
	var m model.SerialNumber
	if err := c.ShouldBindJSON(&m); err != nil {
		return http.StatusBadRequest, nil, err
	}
	m.SerialNumber = 0

	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, &m, nil
}
