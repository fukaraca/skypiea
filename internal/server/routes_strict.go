package server

import (
	"net/http"

	"github.com/fukaraca/skypiea/internal/server/handlers"
	"github.com/fukaraca/skypiea/internal/server/middlewares"
	"github.com/gin-gonic/gin"
)

func strictRoutes(s *Server, common *handlers.Common) RouteMap {
	routes := NewRouteMap()
	h := handlers.NewStrictHandler(common, s.Service, s.Service)
	routes[RouteKey{http.MethodDelete, "/logout"}] = []gin.HandlerFunc{h.Logout}
	routes[RouteKey{http.MethodPut, "/profile/update"}] = []gin.HandlerFunc{h.ProfileUpdate}
	routes[RouteKey{http.MethodPost, "/password"}] = []gin.HandlerFunc{h.ChangePassword}
	routes[RouteKey{http.MethodPost, "/message"}] = []gin.HandlerFunc{h.PostMessage}
	routes[RouteKey{http.MethodGet, "/message/:conv_id/:msg_id/response"}] = []gin.HandlerFunc{h.ResponseOfMessage}
	routes[RouteKey{http.MethodGet, "/conversations/:conv_id/messages"}] = []gin.HandlerFunc{h.GetMessagesByConversationID}
	routes[RouteKey{http.MethodDelete, "/conversations/:conv_id"}] = []gin.HandlerFunc{h.DeleteConversationByID}
	routes[RouteKey{http.MethodGet, "/adminship"}] = []gin.HandlerFunc{middlewares.AdminAuthMw(), h.GetAdminPage}
	routes[RouteKey{http.MethodPut, "/adminship/role"}] = []gin.HandlerFunc{middlewares.AdminAuthMw(), h.UpdateRole}
	return routes
}
