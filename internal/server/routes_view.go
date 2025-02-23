package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fukaraca/skypiea/internal/server/handlers"
	"github.com/fukaraca/skypiea/internal/server/middlewares"
)

func viewRoutes() RouteMap {
	routes := NewRouteMap()
	h := handlers.Handler{}
	routes[RouteKey{http.MethodGet, "/"}] = []gin.HandlerFunc{h.Index, middlewares.TokenAuthMw()}
	routes[RouteKey{http.MethodGet, "/contact"}] = []gin.HandlerFunc{h.Contact}
	routes[RouteKey{http.MethodGet, "/features"}] = []gin.HandlerFunc{h.Features}
	routes[RouteKey{http.MethodGet, "/login"}] = []gin.HandlerFunc{h.Login}
	routes[RouteKey{http.MethodGet, "/faq"}] = []gin.HandlerFunc{h.FAQ}
	routes[RouteKey{http.MethodGet, "/about"}] = []gin.HandlerFunc{h.About}
	routes[RouteKey{http.MethodGet, "/signup"}] = []gin.HandlerFunc{h.Signup}
	routes[RouteKey{http.MethodGet, "/forgot-password"}] = []gin.HandlerFunc{h.ForgotPassword}
	return routes
}
