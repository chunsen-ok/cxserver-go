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
	group.POST("/", r.newPost)
	group.DELETE("/:id", r.delPost)
	group.GET("/", r.getPosts)
	group.GET("/:id", r.getPost)
	group.PUT("/", r.updatePost)
	group.PUT("/status/:id", r.updatePostStatus)
}

type PostWithTagIDs struct {
	model.Post
	TagIDs []int `json:"tags"`
}

type PostWithTags struct {
	model.Post
	Tags []model.Tag `json:"tags"`
}

// param: PostWithTagIDs body
// return: PostWithTags
func (r *Router) newPost(c *gin.Context) {
	var m PostWithTagIDs
	if err := c.ShouldBindJSON(&m); err != nil {
		es := err.Error()
		c.JSON(http.StatusBadRequest, types.Response{Err: &es})
		return
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

		if err := r.db.Where("id in (select tag_id from post_tags where post_id = ?)", m.ID).Find(&tags).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		es := err.Error()
		c.JSON(http.StatusInternalServerError, types.Response{Err: &es})
		return
	}

	rm := PostWithTags{
		Post: m.Post,
		Tags: tags,
	}
	c.JSON(http.StatusOK, types.Response{Body: &rm})
}

// param: id path
// param: del query "?del=1 remove from database, else set state to 'trash'."
// return: null
func (r *Router) delPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		es := err.Error()
		c.JSON(http.StatusOK, types.Response{Err: &es})
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
			if err := tx.Model(&model.Post{}).Where(`id = ?`, id).Update("status", types.Trash).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		es := err.Error()
		c.JSON(http.StatusInternalServerError, types.Response{Err: &es})
		return
	}

	c.JSON(http.StatusOK, types.Response{})
}

// param: tags query "标签ID列表"
// param: status query "默认获取非trash的。?status=0"
// return: []model.Post
func (r *Router) getPosts(c *gin.Context) {
	tags := c.QueryArray("tag")
	status, _ := strconv.Atoi(c.Query("status"))

	posts := make([]model.Post, 0)
	d := r.db.Omit("content").Order("updated_at DESC").Where(`status = ?`, status)
	if len(tags) > 0 {
		d = d.Where("id in (select post_id from post_tags where tag_id in ?)", tags)
	}
	if err := d.Find(&posts).Error; err != nil {
		es := err.Error()
		c.JSON(http.StatusInternalServerError, types.Response{Err: &es})
	}

	c.JSON(http.StatusOK, types.Response{Body: posts})
}

// param: id query "post id"
// return: model.Post
func (r *Router) getPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		es := err.Error()
		c.JSON(http.StatusOK, types.Response{Err: &es})
	}

	var p model.Post
	if err := r.db.Find(&p, id).Error; err != nil {
		es := err.Error()
		c.JSON(http.StatusInternalServerError, types.Response{Err: &es})
		return
	}

	tags := make([]model.Tag, 0)
	if err := r.db.Where(`id in (select tag_id from post_tags where post_id = ?)`, id).Find(&tags).Error; err != nil {
		return
	}

	m := PostWithTags{
		Post: p,
		Tags: tags,
	}

	c.JSON(http.StatusOK, types.Response{Body: &m})
}

// param: PostWithTagIDs body
// return: PostWithTags
func (r *Router) updatePost(c *gin.Context) {
	var m PostWithTagIDs
	if err := c.ShouldBindJSON(&m); err != nil {
		es := err.Error()
		c.JSON(http.StatusBadRequest, types.Response{Err: &es})
		return
	}

	tags := make([]model.Tag, 0)
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Post{}).Omit("id").Where(`id = ?`, m.ID).Updates(&m.Post).Error; err != nil {
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

		if err := tx.Where(`id in (select tag_id from post_tags where post_id = ?)`, m.ID).Find(&tags).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		es := err.Error()
		c.JSON(http.StatusInternalServerError, types.Response{Err: &es})
		return
	}

	rm := PostWithTags{
		Post: m.Post,
		Tags: tags,
	}
	c.JSON(http.StatusOK, types.Response{Body: &rm})
}

// route: /api/posts/status/:id
// param: id path
// param: status query "?status=1"
func (r *Router) updatePostStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		es := err.Error()
		c.JSON(http.StatusOK, types.Response{Err: &es})
		return
	}

	status, _ := strconv.Atoi(c.Query("status"))

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Post{}).Where(`id = ?`, id).Update("status", status).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		es := err.Error()
		c.JSON(http.StatusInternalServerError, types.Response{Err: &es})
		return
	}

	c.JSON(http.StatusOK, types.Response{})
}
