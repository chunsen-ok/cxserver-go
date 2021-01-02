package router

import (
	"context"
	"cxfw/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	err = tx.QueryRow(context.Background(), `insert into tags values (default, $1, now()) returing id;`, m.Title).Scan(&m.ID)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, &m, nil
}

func (r *Router) delTag(c *gin.Context) (int, interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return http.StatusOK, nil, err
	}

	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	_, err = tx.Exec(context.Background(), `delete from tags where id = $1`, id)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, nil, nil
}

func (r *Router) getTags(c *gin.Context) (int, interface{}, error) {
	rows, err := r.db.Query(context.Background(), `select * from tags order by created_at asc;`)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	tags := make([]model.Tag, 0)
	for rows.Next() {
		var tag model.Tag
		err := rows.Scan(&tag.ID, &tag.Title, &tag.CreatedAt)
		if err != nil {
			rows.Close()
			return http.StatusInternalServerError, nil, err
		}
		tags = append(tags, tag)
	}
	rows.Close()

	return http.StatusOK, tags, nil
}

func (r *Router) getTag(c *gin.Context) (int, interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return http.StatusOK, nil, err
	}

	var m model.Tag
	err = r.db.QueryRow(context.Background(), `select * from tags where id = $1`, id).
		Scan(&m.ID, &m.Title, &m.CreatedAt)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, &m, nil
}

func (r *Router) updateTag(c *gin.Context) (int, interface{}, error) {
	var m model.Tag
	if err := c.ShouldBindJSON(&m); err != nil {
		return http.StatusBadRequest, nil, err
	}

	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	_, err = tx.Exec(context.Background(), `update tags set title = $1 where id = $1`, m.Title, m.ID)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, &m, nil
}
