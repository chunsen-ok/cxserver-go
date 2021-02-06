package router

import (
	"cxfw/conf"
	"cxfw/router/fragments"
	"cxfw/router/todos"
	"cxfw/router/writer"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Router ...
type Router struct {
	writer *writer.WriterRouter
	todos  *todos.Service
}

// Init ...
func New(db *pgxpool.Pool) *Router {
	return &Router{
		writer: writer.New(db),
		todos:  todos.New(db),
	}
}

// Routes .
func (s *Router) Routes(router gin.IRouter) {
	g := router.Group("/api", gin.BasicAuth(conf.Instance().BasicAuth))
	s.writer.Init(g)
	s.todos.Init(g)
	fragments.Init(g)
}
