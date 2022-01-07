package usecase

import (
	"ELKExample/repository"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Handler struct {
	DB repository.Databases
	Logger zerolog.Logger
	ESClient *elasticsearch.Client
}

func New(db repository.Databases, esClient *elasticsearch.Client, logger zerolog.Logger) *Handler {
	return &Handler{
		DB: db,
		Logger: logger,
		ESClient: esClient,
	}
}

func (h *Handler) Register(group *gin.RouterGroup) {

	group.DELETE("/posts/:id", h.DeletePost)
	group.PATCH("/posts/:id", h.UpdatePosts)
	group.GET("/posts", h.GetPosts)
	group.POST("/posts", h.CreatePost)

	group.GET("/search", h.SearchPosts)
}
