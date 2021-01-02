package router

import (
	"context"
	"cxfw/model"
	"net/http"
	"strconv"

	"cxfw/types"

	"github.com/gin-gonic/gin"
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

	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	defer tx.Rollback(context.Background())

	if m.BadgeName == types.BadgeRank {
		_, err := tx.Exec(context.Background(), `delete from post_badges where badge_name = $1 and post_id = $2`, types.BadgeRank, m.PostID)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}

		rows, err := tx.Query(context.Background(), `select * from post_badges where badge_name = $1 order by badge_value asc;`, types.BadgeRank)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}

		pbs := make([]model.PostBadge, 0)
		for rows.Next() {
			var pb model.PostBadge
			err := rows.Scan(&pb.BadgeName, &pb.BadgeValue, &pb.PostID)
			if err != nil {
				rows.Close()
				return http.StatusInternalServerError, nil, err
			}
			pbs = append(pbs, pb)
		}
		rows.Close()

		_, err = tx.Exec(context.Background(), `delete from post_badges where badge_name = $1`, types.BadgeRank)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}

		pbs = append(pbs, m)
		for i := 0; i < len(pbs); i++ {
			v := strconv.Itoa(i)
			_, err := tx.Exec(context.Background(), `insert into post_badges values ($1, $2, $3);`, types.BadgeRank, v, pbs[i].PostID)
			if err != nil {
				return http.StatusInternalServerError, nil, err
			}
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
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

	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	defer tx.Rollback(context.Background())

	if m.BadgeName == types.BadgeRank {
		_, err := tx.Exec(context.Background(), `delete from post_badges where badge_name = $1 and post_id = $2`, types.BadgeRank, m.PostID)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}

		rows, err := tx.Query(context.Background(), `select * from post_badges where badge_name = $1 order by badge_value asc;`, types.BadgeRank)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}

		pbs := make([]model.PostBadge, 0)
		for rows.Next() {
			var pb model.PostBadge
			err := rows.Scan(&pb.BadgeName, &pb.BadgeValue, &pb.PostID)
			if err != nil {
				rows.Close()
				return http.StatusInternalServerError, nil, err
			}
			pbs = append(pbs, pb)
		}
		rows.Close()

		_, err = tx.Exec(context.Background(), `delete from post_badges where badge_name = $1`, types.BadgeRank)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}

		for i := 0; i < len(pbs); i++ {
			v := strconv.Itoa(i)
			_, err := tx.Exec(context.Background(), `insert into post_badges values ($1, $2, $3);`, types.BadgeRank, v, pbs[i].PostID)
			if err != nil {
				return http.StatusInternalServerError, nil, err
			}
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, nil, nil
}
