package writer

import (
	"context"
	"cxfw/model/writer"
	"cxfw/router/internal/router"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *WriterRouter) tagsRoutes(ro gin.IRouter) {
	g := ro.Group("/tags")
	g.POST("/", router.Route(r.newTag))
	g.DELETE("/:id", router.Route(r.delTag))
	g.GET("/", router.Route(r.getTags))
	g.GET("/:id", router.Route(r.getTag))
	g.PUT("/", router.Route(r.updateTag))
}

func (r *WriterRouter) newTag(c *gin.Context) (int, interface{}, error) {
	var m writer.Tag
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
func (r *WriterRouter) delTag(c *gin.Context) (int, interface{}, error) {
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

func (r *WriterRouter) getTags(c *gin.Context) (int, interface{}, error) {
	rows, err := r.db.Query(context.Background(), `select * from tags order by created_at asc;`)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	tags := make([]writer.Tag, 0)
	for rows.Next() {
		var tag writer.Tag
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

func (r *WriterRouter) getTag(c *gin.Context) (int, interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return http.StatusOK, nil, err
	}

	var m writer.Tag
	err = r.db.QueryRow(context.Background(), `select * from tags where id = $1`, id).
		Scan(&m.ID, &m.Title, &m.CreatedAt)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, &m, nil
}

func (r *WriterRouter) updateTag(c *gin.Context) (int, interface{}, error) {
	var m writer.Tag
	if err := c.ShouldBindJSON(&m); err != nil {
		return http.StatusBadRequest, nil, err
	}

	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	_, err = tx.Exec(context.Background(), `update tags set title = $1 where id = $2`, m.Title, m.ID)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, &m, nil
}
