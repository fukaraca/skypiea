package handlers

import (
	"fmt"
	"net/http"

	"github.com/fukaraca/skypiea/internal/model"
	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/gin-gonic/gin"
)

func (s *Strict) Logout(c *gin.Context) {
	ck, err := c.Cookie(session.DefaultCookieName)
	if err != nil {
		s.AlertUI(c, fmt.Sprintf("Couldn't log out: %v", err), AlertLevelError)
		return
	}
	session.Cache.Delete(ck)
	c.Header(model.HxRedirect, "/")
	c.Status(http.StatusFound)
}
