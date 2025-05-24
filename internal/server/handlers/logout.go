package handlers

import (
	"net/http"

	"github.com/fukaraca/skypiea/internal/model"
	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/gin-gonic/gin"
)

func (s *Strict) Logout(c *gin.Context) {
	ck, err := c.Cookie(session.DefaultCookieName)
	if err != nil {
		s.AlertUI(c, err, ALError)
		return
	}
	session.Cache.Delete(ck)
	c.Header(model.HxRedirect, model.PathMain)
	c.Status(http.StatusFound)
}
