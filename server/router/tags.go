package router

import (
	"cxfw/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//  /usr/pgsql-11/bin

// /var/lib/pgsql/11/data

func (r *Router) tagsRoutes(g gin.IRouter) {
	tagRouter := g.Group("/tags")
	tagRouter.POST("/", route(r.newTag))
	tagRouter.DELETE("/:id", route(r.delTag))
	tagRouter.GET("/", route(r.getTags))
	tagRouter.GET("/:id", route(r.getTag))
	tagRouter.PUT("/", route(r.updateTag))
}

func (r *Router) newTag(c *gin.Context) (int, interface{}, error) {
	var m model.Tag
	if err := c.ShouldBindJSON(&m); err != nil {
		return http.StatusBadRequest, nil, err
	}

	m.ID = 0
	err := r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&m).Error
	})
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, &m, nil
}

func (r *Router) delTag(c *gin.Context) (int, interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return http.StatusOK, nil, err
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&model.Tag{}, id).Error
	})
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, nil, nil
}

func (r *Router) getTags(c *gin.Context) (int, interface{}, error) {
	tags := make([]model.Tag, 0)
	if err := r.db.Omit("content").Find(&tags).Error; err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, tags, nil
}

func (r *Router) getTag(c *gin.Context) (int, interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return http.StatusOK, nil, err
	}

	var m model.Tag
	if err := r.db.Find(&m, id).Error; err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, &m, nil
}

func (r *Router) updateTag(c *gin.Context) (int, interface{}, error) {
	var m model.Tag
	if err := c.ShouldBindJSON(&m); err != nil {
		return http.StatusBadRequest, nil, err
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&m).Omit("id").Updates(&m).Error
		if err != nil {
			return err
		}

		return tx.First(&m, m.ID).Error
	})
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, &m, nil
}
