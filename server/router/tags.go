package router

import (
	"context"
	"cxfw/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

	err = tx.QueryRow(context.Background(),
		`insert into tags values (default, $1, now() at time zone 'utc', $2) returning id;`,
		m.Title, m.Parent).Scan(&m.ID)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, &m, nil
}

// 删除一个标签时，将其子标签放到其父标签下
// 关联的 post直接删除关联关系即可
func (r *Router) delTag(c *gin.Context) (int, interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return http.StatusOK, nil, err
	}

	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	p := 0
	err = tx.QueryRow(context.Background(), `delete from tags where id = $1 returning parent;`, id).Scan(&p)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	_, err = tx.Exec(context.Background(), `update tags set parent = $1 where parent = $2`, p, id)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, nil, nil
}

func (r *Router) getTags(c *gin.Context) (int, interface{}, error) {
	rows, err := r.db.Query(context.Background(), `select * from tags;`)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	tags := make([]model.Tag, 0)
	for rows.Next() {
		var tag model.Tag
		err := rows.Scan(&tag.ID, &tag.Title, &tag.CreatedAt, &tag.Parent)
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
