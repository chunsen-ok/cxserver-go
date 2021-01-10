package todos

import (
	"cxfw/router/todos/dao"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// url: [POST] /api/todos/items/
func (s *Service) New(c *gin.Context) (int, interface{}, error) {
	var m dao.NewTodoItemParam
	if err := c.ShouldBindJSON(&m); err != nil {
		return http.StatusBadRequest, nil, err
	}

	code, data, err := s.dao.New(&m)

	return code, data, err
}

// url: [DELETE] /api/todos/items/{id}
func (s *Service) Del(c *gin.Context) (int, interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	code, err := s.dao.Del(id)

	return code, nil, err
}

// url: [GET] /api/todos/items/{id}
func (s *Service) Get(c *gin.Context) (int, interface{}, error) {
	return http.StatusOK, nil, nil
}

// url: [GET] /api/todos/items/
// param: dimen query int "0:重要性, 1:紧急性, 2:截止时间"
// param: task query int "task id; 0:所有"
// response: []todos.TodoItem
func (s *Service) GetAll(c *gin.Context) (int, interface{}, error) {
	dimen, err := strconv.Atoi(c.Query("dimen"))
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	taskID, _ := strconv.Atoi(c.Query("task"))

	fmt.Println(dimen, taskID)
	code, data, err := s.dao.GetAll(dimen, taskID)

	return code, data, err
}

// url: [PUT] /api/todos/items/
func (s *Service) Update(c *gin.Context) (int, interface{}, error) {
	return http.StatusOK, nil, nil
}
