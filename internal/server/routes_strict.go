package server

import (
	"net/http"

	"github.com/fukaraca/skypiea/internal/server/handlers"
	"github.com/gin-gonic/gin"
)

func strictRoutes(s *Server) RouteMap {
	routes := NewRouteMap()
	h := handlers.Strict{Repo: s.Repo}
	routes[RouteKey{http.MethodDelete, "/logout"}] = []gin.HandlerFunc{h.Logout}
	return routes
}
