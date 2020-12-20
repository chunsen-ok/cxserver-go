package router

import (
	"cxfw/model"
	"cxfw/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (r *Router) genSerialNumber(c *gin.Context) {
	var m model.SerialNumber
	if err := c.ShouldBindJSON(&m); err != nil {
		es := err.Error()
		c.JSON(http.StatusBadRequest, types.Response{Err: &es})
		return
	}

	m.SerialNumber = 0
	err := r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&m).Error
	})
	if err != nil {
		es := err.Error()
		c.JSON(http.StatusInternalServerError, types.Response{Err: &es})
		return
	}

	c.JSON(http.StatusOK, types.Response{Body: &m})
}
