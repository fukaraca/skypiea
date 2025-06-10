package server

import (
	"net/http"

	"github.com/fukaraca/skypiea/internal/server/handlers"
	"github.com/fukaraca/skypiea/internal/server/middlewares"
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

func openRoutes(s *Server, common *handlers.Common) RouteMap {
	routes := NewRouteMap()
	h := handlers.NewOpenHandler(common, s.Service)
	routes[RouteKey{http.MethodGet, "/healthz"}] = []gin.HandlerFunc{middlewares.NonAuthMw(), h.Healthz}
	routes[RouteKey{http.MethodPost, "/login"}] = []gin.HandlerFunc{middlewares.NonAuthMw(), h.SignIn}
	routes[RouteKey{http.MethodPost, "/signup"}] = []gin.HandlerFunc{middlewares.NonAuthMw(), h.SignUp}
	routes[RouteKey{http.MethodPost, "/forgot-password"}] = []gin.HandlerFunc{middlewares.NonAuthMw(), h.ForgotPassword}

	return routes
}
