package server

import (
	"log/slog"

	"github.com/fukaraca/skypiea/internal/api/gemini"
	"github.com/fukaraca/skypiea/internal/service"
	"github.com/fukaraca/skypiea/internal/storage"
	"github.com/gin-gonic/gin"

	"github.com/fukaraca/skypiea/internal/config"
)

const V1 = "/v1"

type Server struct {
	Config  *config.Config
	engine  *gin.Engine
	Logger  *slog.Logger
	Repo    *storage.Registry
	Service *service.Service
}

func NewServer(cfg *config.Config, engine *gin.Engine, db *storage.DB, logger *slog.Logger) (*Server, error) {
	repo := storage.NewRegistry(db.Dialect, db.Pool)
	apiClient, err := gemini.NewClient(cfg.Gemini)
	if err != nil {
		return nil, err
	}
	srv := service.New(repo, apiClient)
	return &Server{
		Config:  cfg,
		Logger:  logger,
		engine:  engine,
		Repo:    repo,
		Service: srv,
	}, nil
}
