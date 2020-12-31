package router

import (
	"cxfw/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (r *Router) badgesRoutes(g gin.IRouter) {
	group := g.Group("/badges")
	group.POST("/", route(r.newPostBadge))
	group.DELETE("/", route(r.removePostBadge))
}

// route: [POST] /api/badges/
// param: id path int "post id"
// param: name path int "badge name by badge enums"
// param: value query string "badge value"
func (r *Router) newPostBadge(c *gin.Context) (int, interface{}, error) {
	var m model.PostBadge
	if err := c.ShouldBindJSON(&m); err != nil {
		return http.StatusBadRequest, nil, err
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&m).Error
	})
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, nil, nil
}

// route: [DELETE] /api/badges/
// param: data body model.PostBadge "match data"
func (r *Router) removePostBadge(c *gin.Context) (int, interface{}, error) {
	var m model.PostBadge
	if err := c.ShouldBindJSON(&m); err != nil {
		return http.StatusBadRequest, nil, err
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&m).Error
	})
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, nil, nil
}
