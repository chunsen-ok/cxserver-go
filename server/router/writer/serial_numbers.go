package writer

import (
	"context"
	"cxfw/model/writer"
	"cxfw/router/internal/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *WriterRouter) snRoutes(ro gin.IRouter) {
	g := ro.Group("/sn")
	g.POST("/", router.Route(r.genSerialNumber))
}

func (r *WriterRouter) genSerialNumber(c *gin.Context) (int, interface{}, error) {
	var m writer.SerialNumber
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
