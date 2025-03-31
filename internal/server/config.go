package server

import (
	"log/slog"

	"github.com/fukaraca/skypiea/internal/storage"
	"github.com/gin-gonic/gin"

	"github.com/fukaraca/skypiea/internal/config"
)

const V1 = "/v1"

type Server struct {
	Config *config.Config
	engine *gin.Engine
	Logger *slog.Logger
	Repo   *storage.Repositories
}

func NewServer(cfg *config.Config, engine *gin.Engine, db *storage.DB, logger *slog.Logger) *Server {
	return &Server{
		Config: cfg,
		Logger: logger,
		engine: engine,
		Repo:   storage.NewRepositories(db),
	}
}
