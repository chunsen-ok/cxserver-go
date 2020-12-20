package router

import (
	"cxfw/model"
	"cxfw/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//  /usr/pgsql-11/bin

// /var/lib/pgsql/11/data

func (r *Router) newTag(c *gin.Context) {
	var m model.Tag
	if err := c.ShouldBindJSON(&m); err != nil {
		es := err.Error()
		c.JSON(http.StatusBadRequest, types.Response{Err: &es})
		return
	}

	m.ID = 0
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

func (r *Router) delTag(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		es := err.Error()
		c.JSON(http.StatusOK, types.Response{Err: &es})
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&model.Tag{}, id).Error
	})
	if err != nil {
		es := err.Error()
		c.JSON(http.StatusInternalServerError, types.Response{Err: &es})
		return
	}

	c.JSON(http.StatusOK, types.Response{})
}

func (r *Router) getTags(c *gin.Context) {
	tags := make([]model.Tag, 0)
	if err := r.db.Omit("content").Find(&tags).Error; err != nil {
		es := err.Error()
		c.JSON(http.StatusInternalServerError, types.Response{Err: &es})
	}

	c.JSON(http.StatusOK, types.Response{Body: tags})
}

func (r *Router) getTag(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		es := err.Error()
		c.JSON(http.StatusOK, types.Response{Err: &es})
	}

	var m model.Tag
	if err := r.db.Find(&m, id).Error; err != nil {
		es := err.Error()
		c.JSON(http.StatusInternalServerError, types.Response{Err: &es})
		return
	}

	c.JSON(http.StatusOK, types.Response{Body: &m})
}

func (r *Router) updateTag(c *gin.Context) {
	var m model.Tag
	if err := c.ShouldBindJSON(&m); err != nil {
		es := err.Error()
		c.JSON(http.StatusBadRequest, types.Response{Err: &es})
		return
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&m).Omit("id").Updates(&m).Error
		if err != nil {
			return err
		}

		return tx.First(&m, m.ID).Error
	})
	if err != nil {
		es := err.Error()
		c.JSON(http.StatusInternalServerError, types.Response{Err: &es})
		return
	}

	c.JSON(http.StatusOK, types.Response{Body: &m})
}
