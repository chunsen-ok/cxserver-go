package items

import (
	"cxfw/service/internal/router"
	"cxfw/service/todos/items/dao"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Init(r gin.IRouter) {
	g := r.Group("/items")
	g.POST("/", router.Route(add))
	g.DELETE("/:id", router.Route(del))
	g.GET("/:id", router.Route(get))
	g.GET("/", router.Route(getAll))
	g.PUT("/", router.Route(update))

}

// url: [POST] /api/todos/items/
func add(c *gin.Context) (int, interface{}, error) {
	var m dao.NewTodoItemParam
	if err := c.ShouldBindJSON(&m); err != nil {
		return http.StatusBadRequest, nil, err
	}

	code, data, err := dao.Add(&m)

	return code, data, err
}

// url: [DELETE] /api/todos/items/{id}
func del(c *gin.Context) (int, interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	code, err := dao.Del(id)

	return code, nil, err
}

// url: [GET] /api/todos/items/{id}
func get(c *gin.Context) (int, interface{}, error) {
	return http.StatusOK, nil, nil
}

// url: [GET] /api/todos/items/
// param: dimen query int "0:重要性, 1:紧急性, 2:截止时间"
// param: task query int "task id; 0:所有"
// response: []todos.TodoItem
func getAll(c *gin.Context) (int, interface{}, error) {
	dimen, err := strconv.Atoi(c.Query("dimen"))
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	taskID, _ := strconv.Atoi(c.Query("task"))

	fmt.Println(dimen, taskID)
	code, data, err := dao.GetAll(dimen, taskID)

	return code, data, err
}

// url: [PUT] /api/todos/items/
func update(c *gin.Context) (int, interface{}, error) {
	return http.StatusOK, nil, nil
}
