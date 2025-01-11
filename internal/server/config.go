package service

import (
	"github.com/fukaraca/skypiea/internal/config"
	"github.com/gin-gonic/gin"
	"log/slog"
)

const V1 = "/v1"

type Server struct {
	Config *config.Config
	engine *gin.Engine
	Logger *slog.Logger
}

func NewServer(cfg *config.Config, engine *gin.Engine, logger *slog.Logger) *Server {
	return &Server{
		Config: cfg,
		Logger: logger,
		engine: engine,
	}
}
