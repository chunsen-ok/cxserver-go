package writer

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type WriterRouter struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *WriterRouter {
	return &WriterRouter{
		db: db,
	}
}

func (s *WriterRouter) Init(r gin.IRouter) {
	g := r.Group("/writer")
	s.badgesRoutes(g)
	s.tagsRoutes(g)
	// s.snRoutes(g)
	s.postsRoutes(g)
}
