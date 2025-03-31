package server

import (
	"net/http"

	"github.com/fukaraca/skypiea/internal/server/handlers"
	"github.com/gin-gonic/gin"
)

type RouteKey struct {
	Method string
	Path   string
}

type RouteMap map[RouteKey]gin.HandlersChain

func NewRouteMap() RouteMap {
	return make(RouteMap)
}

func (s *Server) RegisterRoutes(rGroup *gin.RouterGroup, routes RouteMap) {
	if rGroup != nil {
		for k, v := range routes {
			rGroup.Handle(k.Method, k.Path, v...)
		}
	} else {
		for k, v := range routes {
			s.engine.Handle(k.Method, k.Path, v...)
		}
	}
}

func commonRoutes(s *Server) RouteMap {
	routes := NewRouteMap()
	h := handlers.Handler2{Repo: s.Repo}
	routes[RouteKey{http.MethodPost, "/login"}] = []gin.HandlerFunc{h.SignIn}
	routes[RouteKey{http.MethodPost, "/signup"}] = []gin.HandlerFunc{h.SignUp}

	return routes
}
