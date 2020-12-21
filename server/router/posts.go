package router

import (
	"cxfw/model"
	"cxfw/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostWithTags struct {
	model.Post
	Tags []model.Tag `json:"tags"`
}

func (r *Router) newPost(c *gin.Context) {
	var m PostWithTags
	if err := c.ShouldBindJSON(&m); err != nil {
		es := err.Error()
		c.JSON(http.StatusBadRequest, types.Response{Err: &es})
		return
	}

	m.ID = 0
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&m.Post).Error; err != nil {
			return err
		}

		postTags := make([]model.PostTag, 0)
		for _, tag := range m.Tags {
			pt := model.PostTag{
				PostID: m.ID,
				TagID:  tag.ID,
			}

			postTags = append(postTags, pt)
		}

		return tx.Create(&postTags).Error
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
		if err := tx.Delete(&model.Post{}, id).Error; err != nil {
			return err
		}

		return tx.Where("post_id = ?", id).Delete(&model.PostTag{}).Error
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

func (r *Router) updatePost(c *gin.Context) {
	var m PostWithTags
	if err := c.ShouldBindJSON(&m); err != nil {
		es := err.Error()
		c.JSON(http.StatusBadRequest, types.Response{Err: &es})
		return
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Post{}).Omit("id").Where(`id = ?`, m.ID).Updates(&m.Post).Error; err != nil {
			return err
		}

		if err := tx.Where(`post_id = ?`, m.ID).Delete(&model.PostTag{}).Error; err != nil {
			return err
		}

		postTags := make([]model.PostTag, 0)
		for _, tag := range m.Tags {
			pt := model.PostTag{
				PostID: m.ID,
				TagID:  tag.ID,
			}

			postTags = append(postTags, pt)
		}

		if err := tx.Create(&postTags).Error; err != nil {
			return err
		}

		return tx.First(&m.Post, m.ID).Error
	})
	if err != nil {
		es := err.Error()
		c.JSON(http.StatusInternalServerError, types.Response{Err: &es})
		return
	}

	c.JSON(http.StatusOK, types.Response{Body: &m})
}
