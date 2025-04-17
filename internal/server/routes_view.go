package server

import (
	"net/http"

	"github.com/fukaraca/skypiea/internal/server/handlers"
	"github.com/fukaraca/skypiea/internal/server/middlewares"
	"github.com/gin-gonic/gin"
)

func viewRoutes(s *Server, common *handlers.Common) RouteMap {
	routes := NewRouteMap()
	h := handlers.NewViewHandler(common, s.Repo)
	routes[RouteKey{http.MethodGet, "/"}] = []gin.HandlerFunc{h.Index}
	routes[RouteKey{http.MethodGet, "/contact"}] = []gin.HandlerFunc{h.Contact}
	routes[RouteKey{http.MethodGet, "/pricing"}] = []gin.HandlerFunc{h.Pricing}
	routes[RouteKey{http.MethodGet, "/features"}] = []gin.HandlerFunc{h.Features}
	routes[RouteKey{http.MethodGet, "/login"}] = []gin.HandlerFunc{middlewares.NonAuthMw(), h.Login}
	routes[RouteKey{http.MethodGet, "/faq"}] = []gin.HandlerFunc{h.FAQ}
	routes[RouteKey{http.MethodGet, "/about"}] = []gin.HandlerFunc{h.About}
	routes[RouteKey{http.MethodGet, "/profile"}] = []gin.HandlerFunc{h.Profile}
	routes[RouteKey{http.MethodGet, "/signup"}] = []gin.HandlerFunc{middlewares.NonAuthMw(), h.Signup}
	routes[RouteKey{http.MethodGet, "/forgot-password"}] = []gin.HandlerFunc{middlewares.NonAuthMw(), h.ForgotPassword}
	return routes
}
