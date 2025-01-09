package service

import (
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"net/http"
)

func (s *Server) bindRoutes() {
	s.engine.NoRoute()
	s.engine.GET("/ping", func(c *gin.Context) {
		s.Logger.Info("ping ponged", "id", sloggin.GetRequestID(c))
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	s.engine.GET("home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", nil)
	})
}
