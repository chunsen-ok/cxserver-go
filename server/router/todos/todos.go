package todos

import (
	"cxfw/router/internal/router"
	"cxfw/router/todos/dao"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Service struct {
	dao *dao.TodoDao
}

func New(db *pgxpool.Pool) *Service {
	return &Service{
		dao: dao.New(db),
	}
}

func (s *Service) Init(r gin.IRouter) {
	g := r.Group("/todos")

	g1 := g.Group("/items")
	g1.POST("/", router.Route(s.New))
	g1.DELETE("/:id", router.Route(s.Del))
	g1.GET("/:id", router.Route(s.Get))
	g1.GET("/", router.Route(s.GetAll))
	g1.PUT("/", router.Route(s.Update))

	g2 := g.Group("/tasks")
	g2.POST("/", router.Route(s.NewTask))
	g2.DELETE("/:id", router.Route(s.DelTask))
	g2.GET("/", router.Route(s.GetAllTask))
}
