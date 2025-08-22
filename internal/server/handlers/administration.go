package handlers

import (
	"net/http"

	"github.com/fukaraca/skypiea/internal/model"
	gb "github.com/fukaraca/skypiea/pkg/guest_book"
	"github.com/gin-gonic/gin"
)

func (s *Strict) GetAdminPage(c *gin.Context) {
	stats, err := s.UserSvc.GetAdoptionStatistics(c.Request.Context())
	if err != nil {
		s.AlertUI(c, model.ErrSessionNotFound, ALError)
		return
	}
	visitors := gb.GuestBook.DumpVisitorMetric()
	c.HTML(http.StatusOK, "adminship", gin.H{
		"Title":       "Admin Only",
		"LoggedIn":    true,
		"CSSFile":     "index.css",
		"Users":       stats,
		"Visitors":    visitors,
		"RoleOptions": []string{model.RoleAdmin, model.RoleUserStd, model.RoleUserVip},
	})
}

func (s *Strict) UpdateRole(c *gin.Context) {
	userID := c.PostForm("user_uuid")
	role := c.PostForm("role")
	err := s.UserSvc.UpdateRole(c.Request.Context(), userID, role)
	if err != nil {
		s.AlertUI(c, err.Error(), ALError)
	}
	s.AlertUI(c, "Role has been updated", ALInfo)
}
