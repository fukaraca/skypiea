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

func (s *Server) RegisterRoutes(rGroup *gin.RouterGroup, routes RouteMap, optionalMWs ...gin.HandlerFunc) {
	for key, chain := range routes {
		mws := make(gin.HandlersChain, 0)
		mws = append(mws, optionalMWs...)
		mws = append(mws, chain...)

		if rGroup != nil {
			rGroup.Handle(key.Method, key.Path, mws...)
		} else {
			s.engine.Handle(key.Method, key.Path, mws...)
		}
	}
}

func commonRoutes(s *Server) RouteMap {
	routes := NewRouteMap()
	h := handlers.Common{Repo: s.Repo}
	routes[RouteKey{http.MethodPost, "/login"}] = []gin.HandlerFunc{h.SignIn}
	routes[RouteKey{http.MethodPost, "/signup"}] = []gin.HandlerFunc{h.SignUp}

	return routes
}
