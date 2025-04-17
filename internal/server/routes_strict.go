package server

import (
	"net/http"

	"github.com/fukaraca/skypiea/internal/server/handlers"
	"github.com/gin-gonic/gin"
)

func strictRoutes(s *Server, common *handlers.Common) RouteMap {
	routes := NewRouteMap()
	h := handlers.NewStrictHandler(common, s.Repo)
	routes[RouteKey{http.MethodDelete, "/logout"}] = []gin.HandlerFunc{h.Logout}
	routes[RouteKey{http.MethodPost, "/password"}] = []gin.HandlerFunc{h.ChangePassword}
	return routes
}
