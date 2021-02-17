package tasks

import (
	"cxfw/service/internal/router"
	"cxfw/service/todos/tasks/dao"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Init(r gin.IRouter) {
	g := r.Group("/tasks")
	g.POST("/", router.Route(add))
	g.DELETE("/:id", router.Route(del))
	g.GET("/", router.Route(getAll))
}

// url: [POST] /api/todos/tasks/
// param: dao.NewTodoTaskParam body
func add(c *gin.Context) (int, interface{}, error) {
	var m dao.NewTodoTaskParam
	if err := c.ShouldBindJSON(&m); err != nil {
		return http.StatusBadRequest, nil, err
	}

	code, data, err := dao.Add(&m)

	return code, data, err
}

// url: [DELETE] /api/todos/tasks/{id}
func del(c *gin.Context) (int, interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	code, err := dao.Del(id)

	return code, nil, err
}

// url: [GET] /api/todos/tasks/
func getAll(c *gin.Context) (int, interface{}, error) {
	code, data, err := dao.GetAll()

	return code, data, err
}
