package todos

import (
	"cxfw/router/todos/dao"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// url: [POST] /api/todos/tasks/
// param: dao.NewTodoTaskParam body
func (s *Service) NewTask(c *gin.Context) (int, interface{}, error) {
	var m dao.NewTodoTaskParam
	if err := c.ShouldBindJSON(&m); err != nil {
		return http.StatusBadRequest, nil, err
	}

	code, data, err := s.dao.NewTask(&m)

	return code, data, err
}

// url: [DELETE] /api/todos/tasks/{id}
func (s *Service) DelTask(c *gin.Context) (int, interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	code, err := s.dao.DelTask(id)

	return code, nil, err
}

// url: [GET] /api/todos/tasks/
func (s *Service) GetAllTask(c *gin.Context) (int, interface{}, error) {
	code, data, err := s.dao.GetAllTask()

	return code, data, err
}
