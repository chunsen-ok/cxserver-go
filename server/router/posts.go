package router

import (
	"cxfw/model"
	"cxfw/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (r *Router) postsRoutes(router gin.IRouter) {
	group := router.Group("/posts")
	group.POST("/", route(r.newPost))
	group.DELETE("/:id", route(r.delPost))
	group.GET("/", route(r.getPosts))
	group.GET("/:id", route(r.getPost))
	group.PUT("/", route(r.updatePost))
	group.PUT("/status/:id", route(r.updatePostStatus))
}

type PostWithTagIDs struct {
	model.Post
	TagIDs []int `json:"tags"`
}

type PostWithTags struct {
	model.Post
	Tags []model.Tag `json:"tags"`
}

// route: [POST] /api/posts/
// param: data body PostWithTagIDs ""
// response: PostWithTags
func (r *Router) newPost(c *gin.Context) (int, interface{}, error) {
	var m PostWithTagIDs
	if err := c.ShouldBindJSON(&m); err != nil {
		return http.StatusBadRequest, nil, err
	}
	m.ID = 0

	tags := make([]model.Tag, 0)
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&m.Post).Error; err != nil {
			return err
		}

		postTags := make([]model.PostTag, 0)
		for _, tagID := range m.TagIDs {
			pt := model.PostTag{
				PostID: m.ID,
				TagID:  tagID,
			}

			postTags = append(postTags, pt)
		}

		if len(postTags) > 0 {
			if err := tx.Create(&postTags).Error; err != nil {
				return err
			}
		}

		if err := r.db.Where("id in (select tag_id from post_tags where post_id = ?)", m.ID).
			Find(&tags).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	rm := PostWithTags{
		Post: m.Post,
		Tags: tags,
	}

	return http.StatusOK, &rm, nil
}

// param: id path
// param: del query "?del=1 remove from database, else set state to 'trash'."
// return: null
func (r *Router) delPost(c *gin.Context) (int, interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return http.StatusOK, nil, err
	}

	del, _ := strconv.Atoi(c.Query("del"))

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if del == 1 {
			if err := tx.Delete(&model.Post{}, id).Error; err != nil {
				return err
			}

			if err := tx.Where("post_id = ?", id).Delete(&model.PostTag{}).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&model.Post{}).Where(`id = ?`, id).
				Update("status", types.StatusTrash).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusNoContent, nil, nil
}

// param: tags query "标签ID列表"
// param: status query "默认获取非trash的。?status=0"
// return: []model.Post
func (r *Router) getPosts(c *gin.Context) (int, interface{}, error) {
	tags := c.QueryArray("tag")
	status, _ := strconv.Atoi(c.Query("status"))

	posts := make([]model.Post, 0)
	d := r.db.Omit("content").Order("updated_at DESC").Where(`status = ?`, status)
	if len(tags) > 0 {
		d = d.Where("id in (select post_id from post_tags where tag_id in ?)", tags)
	}
	if err := d.Find(&posts).Error; err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, posts, nil
}

// param: id query "post id"
// return: model.Post
func (r *Router) getPost(c *gin.Context) (int, interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	var p model.Post
	if err := r.db.Find(&p, id).Error; err != nil {
		return http.StatusInternalServerError, nil, err
	}

	tags := make([]model.Tag, 0)
	if err := r.db.Where(`id in (select tag_id from post_tags where post_id = ?)`, id).
		Find(&tags).Error; err != nil {
		return http.StatusInternalServerError, nil, err
	}

	m := PostWithTags{
		Post: p,
		Tags: tags,
	}

	return http.StatusOK, &m, nil
}

// param: PostWithTagIDs body
// return: PostWithTags
func (r *Router) updatePost(c *gin.Context) (int, interface{}, error) {
	var m PostWithTagIDs
	if err := c.ShouldBindJSON(&m); err != nil {
		return http.StatusBadRequest, nil, err
	}

	tags := make([]model.Tag, 0)
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Post{}).Omit("id").Where(`id = ?`, m.ID).
			Updates(&m.Post).Error; err != nil {
			return err
		}

		if err := tx.Where(`post_id = ?`, m.ID).Delete(&model.PostTag{}).Error; err != nil {
			return err
		}

		postTags := make([]model.PostTag, 0)
		for _, tagID := range m.TagIDs {
			pt := model.PostTag{
				PostID: m.ID,
				TagID:  tagID,
			}

			postTags = append(postTags, pt)
		}

		if len(postTags) > 0 {
			if err := tx.Create(&postTags).Error; err != nil {
				return err
			}
		}

		if err := tx.Where(`id in (select tag_id from post_tags where post_id = ?)`, m.ID).
			Find(&tags).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	rm := PostWithTags{
		Post: m.Post,
		Tags: tags,
	}

	return http.StatusOK, &rm, nil
}

// route: /api/posts/status/:id
// param: id path
// param: status query "?status=1"
func (r *Router) updatePostStatus(c *gin.Context) (int, interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return http.StatusOK, nil, err
	}

	status, _ := strconv.Atoi(c.Query("status"))

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Post{}).Where(`id = ?`, id).Update("status", status).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, nil, nil
}
