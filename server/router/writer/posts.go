package writer

import (
	"context"
	"cxfw/model/writer"
	"cxfw/router/internal/router"
	"cxfw/types"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func (r *WriterRouter) postsRoutes(ro gin.IRouter) {
	g := ro.Group("/posts")
	g.POST("/", router.Route(r.newPost))
	g.DELETE("/:id", router.Route(r.delPost))
	g.GET("/", router.Route(r.getPosts))
	g.GET("/:id", router.Route(r.getPost))
	g.PUT("/", router.Route(r.updatePost))
	g.PUT("/status/:id", router.Route(r.updatePostStatus))
}

type PostWithTagIDs struct {
	writer.Post
	Title  string `json:"title"`
	TagIDs []int  `json:"tags"`
}

type PostWithTags struct {
	writer.Post
	Title string       `json:"title"`
	Tags  []writer.Tag `json:"tags"`
}

type PostBadge struct {
	BadgeName  int     `json:"badge_name"`
	BadgeValue *string `json:"badge_value"`
}

type PostWithBadges struct {
	writer.Post
	Title  string      `json:"title"`
	Badges []PostBadge `json:"badges"`
}

// route: [POST] /api/posts/
// param: data body PostWithTagIDs ""
// response: PostWithTags
func (r *WriterRouter) newPost(c *gin.Context) (int, interface{}, error) {
	var m PostWithTagIDs
	if err := c.ShouldBindJSON(&m); err != nil {
		return http.StatusBadRequest, nil, err
	}
	m.ID = 0

	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	defer tx.Rollback(context.Background())

	err = tx.QueryRow(context.Background(),
		`insert into posts values (default, $1, $2, now() at time zone 'utc', now() at time zone 'utc') returning id`,
		m.Content, types.StatusActive).Scan(&m.ID)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	tags := make([]writer.Tag, 0)
	if len(m.TagIDs) > 0 {
		bt := pgx.Batch{}
		for _, tagID := range m.TagIDs {
			bt.Queue(`insert into post_tags values ($1, $2);`, m.ID, tagID)
		}
		br := tx.SendBatch(context.Background(), &bt)
		if err := br.Close(); err != nil {
			return http.StatusInternalServerError, nil, err
		}

		rows, err := tx.Query(context.Background(), `select * from tags where id in (select tag_id from post_tags where post_id = $1)`, m.ID)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}

		for rows.Next() {
			var tag writer.Tag
			err := rows.Scan(&tag.ID, &tag.Title, nil, nil)
			if err != nil {
				rows.Close()
				return http.StatusInternalServerError, nil, err
			}
			tags = append(tags, tag)
		}
		rows.Close()
	}

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	rm := PostWithTags{
		Post:  m.Post,
		Title: m.Title,
		Tags:  tags,
	}

	return http.StatusOK, &rm, nil
}

// param: id path
// param: del query "?del=1 remove from database, else set state to 'trash'."
// return: null
func (r *WriterRouter) delPost(c *gin.Context) (int, interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return http.StatusOK, nil, err
	}

	del, _ := strconv.Atoi(c.Query("del"))

	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	defer tx.Rollback(context.Background())

	if del == 1 {
		_, err := tx.Exec(context.Background(), `delete from posts where id = $1`, id)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
		_, err = tx.Exec(context.Background(), `delete from post_tags where post_id = $1`, id)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
	} else {
		_, err := tx.Exec(context.Background(), `update posts set status = $1 where id = $2`, types.StatusTrash, id)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
	}

	_, err = tx.Exec(context.Background(), `delete from post_badges where post_id = $1`, id)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, nil, nil
}

// url: [GET] /api/posts/
// param: tags query []int "标签ID列表, -1 表示获取全部，-2 表示获取没有设置标签的"
// param: status query string "默认获取非trash的。?status=0"
// return: []PostWithBadges
func (r *WriterRouter) getPosts(c *gin.Context) (int, interface{}, error) {
	tags := c.QueryArray("tag")
	status, _ := strconv.Atoi(c.Query("status"))

	sb := strings.Builder{}
	sb.WriteString(`with t1 as (
		select post_id, jsonb_build_object('badge_name', badge_name, 'badge_value', badge_value) as badge from post_badges
	),
	t2 as (
		with t as (
			select id, position(chr(10) in "content") as first_line_end from posts
		) select p.id, p.status, p.created_at, p.updated_at,
			case t.first_line_end
				when 0 then p."content"
				else substring(p."content" from 0 for t.first_line_end) end
			as title
		from posts p left join t on t.id = p.id
	) select t2.*, t1.badge from t2 left join t1 on t1.post_id = t2.id `)
	// shit begin
	if len(tags) > 0 {
		if len(tags) == 1 && tags[0] == "-1" {
			// ... get all
			// do nothing
			sb.WriteString(`where t2.status = $1`)
		} else {
			sort.Strings(tags)

			// get untagged posts
			if idx := sort.SearchStrings(tags, "-2"); tags[idx] == "-2" {
				sb.WriteString(`left join post_tags pt on pt.post_id = t2.id where t2.status = $1 and (pt.tag_id = 0 or pt.tag_id is null)`)
			} else {
				sb.WriteString(fmt.Sprintf(`where t2.status = $1 and t2.id in (select post_id from post_tags where tag_id in (%s))`, strings.Join(tags, ",")))
			}
		}
	} else {
		sb.WriteString(`where t2.status = $1`)
	}
	// shit end

	rows, err := r.db.Query(context.Background(), sb.String(), status)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	postsMap := make(map[int]PostWithBadges, 0)
	for rows.Next() {
		var p PostWithBadges
		var b *PostBadge
		err := rows.Scan(&p.ID, &p.Status, &p.CreatedAt, &p.UpdatedAt, &p.Title, &b)
		if err != nil {
			rows.Close()
			return http.StatusInternalServerError, nil, err
		}
		pb, ok := postsMap[p.ID]
		if !ok {
			if b != nil {
				p.Badges = []PostBadge{*b}
			}
			postsMap[p.ID] = p
		} else {
			if b != nil {
				pb.Badges = append(pb.Badges, *b)
			}
			postsMap[p.ID] = pb
		}
	}
	rows.Close()

	posts := make([]PostWithBadges, 0)
	for _, p := range postsMap {
		posts = append(posts, p)
	}

	sort.Slice(posts, func(i, j int) bool {
		lhs := posts[i].Badges
		rhs := posts[j].Badges

		if lhs == nil && rhs == nil {
			return posts[j].UpdatedAt.Before(posts[i].UpdatedAt)
		}

		if lhs == nil {
			return false
		}

		if rhs == nil {
			return true
		}

		var lv *string
		for _, b := range lhs {
			if b.BadgeName == types.BadgeRank {
				lv = b.BadgeValue
				break
			}
		}

		var rv *string
		for _, b := range rhs {
			if b.BadgeName == types.BadgeRank {
				rv = b.BadgeValue
				break
			}
		}

		if lv == nil {
			return false
		}

		if rv == nil {
			return true
		}

		if *lv > *rv {
			return false
		}

		return true
	})

	return http.StatusOK, posts, nil
}

// param: id query "post id"
// return: writer.Post
func (r *WriterRouter) getPost(c *gin.Context) (int, interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	var p writer.Post
	err = r.db.QueryRow(context.Background(), `select * from posts where id = $1`, id).
		Scan(&p.ID, &p.Content, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	rows, err := r.db.Query(context.Background(),
		`select * from tags where id in (select tag_id from post_tags where post_id = $1)`, id)
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

	m := PostWithTags{
		Post: p,
		Tags: tags,
	}

	return http.StatusOK, &m, nil
}

// param: PostWithTagIDs body
// return: PostWithTags
func (r *WriterRouter) updatePost(c *gin.Context) (int, interface{}, error) {
	var m PostWithTagIDs
	if err := c.ShouldBindJSON(&m); err != nil {
		return http.StatusBadRequest, nil, err
	}

	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	defer tx.Rollback(context.Background())

	bt := pgx.Batch{}
	bt.Queue(`update posts set content = $1, updated_at = now() at time zone 'utc' where id = $2`,
		m.Content, m.ID)
	bt.Queue(`delete from post_tags where post_id = $1`, m.ID)
	for _, tagID := range m.TagIDs {
		bt.Queue(`insert into post_tags values ($1, $2);`, m.ID, tagID)
	}
	br := tx.SendBatch(context.Background(), &bt)
	if err := br.Close(); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	rows, err := tx.Query(context.Background(), `select * from tags where id in (select tag_id from post_tags where post_id = $1)`, m.ID)
	if err := br.Close(); err != nil {
		return http.StatusInternalServerError, nil, err
	}
	tags := make([]writer.Tag, 0)
	for rows.Next() {
		var tag writer.Tag
		err := rows.Scan(&tag.ID, &tag.Title, nil, nil)
		if err != nil {
			rows.Close()
			return http.StatusInternalServerError, nil, err
		}
		tags = append(tags, tag)
	}
	rows.Close()

	if err := tx.Commit(context.Background()); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	rm := PostWithTags{
		Post:  m.Post,
		Title: m.Title,
		Tags:  tags,
	}

	return http.StatusOK, &rm, nil
}

// route: /api/posts/status/:id
// param: id path
// param: status query "?status=1"
func (r *WriterRouter) updatePostStatus(c *gin.Context) (int, interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return http.StatusOK, nil, err
	}

	status, _ := strconv.Atoi(c.Query("status"))

	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	_, err = tx.Exec(context.Background(), `update posts set status = $1 where id = $2`, status, id)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, nil, nil
}
