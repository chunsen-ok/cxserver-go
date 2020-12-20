package router

import (
	"cxfw/model"
	"cxfw/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (r *Router) newPost(c *gin.Context) {
	var m model.Post
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

func (r *Router) delPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		es := err.Error()
		c.JSON(http.StatusOK, types.Response{Err: &es})
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&model.Post{}, id).Error
	})
	if err != nil {
		es := err.Error()
		c.JSON(http.StatusInternalServerError, types.Response{Err: &es})
		return
	}

	c.JSON(http.StatusOK, types.Response{})
}

func (r *Router) getPosts(c *gin.Context) {
	posts := make([]model.Post, 0)
	if err := r.db.Omit("content").Find(&posts).Error; err != nil {
		es := err.Error()
		c.JSON(http.StatusInternalServerError, types.Response{Err: &es})
	}

	c.JSON(http.StatusOK, types.Response{Body: posts})
}

func (r *Router) getPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		es := err.Error()
		c.JSON(http.StatusOK, types.Response{Err: &es})
	}

	var m model.Post
	if err := r.db.Find(&m, id).Error; err != nil {
		es := err.Error()
		c.JSON(http.StatusInternalServerError, types.Response{Err: &es})
		return
	}

	c.JSON(http.StatusOK, types.Response{Body: &m})
}

func (r *Router) updatePost(c *gin.Context) {
	var m model.Post
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
