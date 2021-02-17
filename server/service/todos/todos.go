package todos

import (
	"cxfw/service/todos/items"
	"cxfw/service/todos/tasks"

	"github.com/gin-gonic/gin"
)

func Init(r gin.IRouter) {
	g := r.Group("/todos")

	items.Init(g)
	tasks.Init(g)
}
