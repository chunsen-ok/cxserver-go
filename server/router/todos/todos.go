package todos

import (
	"cxfw/router/internal/router"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Service struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Init(r gin.IRouter) {
	g := r.Group("/todos")

	g.POST("/", router.Route(s.New))
}

func (s *Service) New(c *gin.Context) (int, interface{}, error) {
	return http.StatusOK, nil, nil
}
