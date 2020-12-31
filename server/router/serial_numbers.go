package router

import (
	"cxfw/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	err := r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&m).Error
	})
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, &m, nil
}
