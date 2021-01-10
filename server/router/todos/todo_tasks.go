package todos

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// url: [POST] /api/todos/tasks/
func (s *Service) NewTask(c *gin.Context) (int, interface{}, error) {
	return http.StatusOK, nil, nil
}

// url: [GET] /api/todos/tasks/
func (s *Service) GetAllTask(c *gin.Context) (int, interface{}, error) {
	code, data, err := s.dao.GetAllTask()

	return code, data, err
}
