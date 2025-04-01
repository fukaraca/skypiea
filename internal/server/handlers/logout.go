package handlers

import (
	"net/http"

	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/gin-gonic/gin"
)

func (s *Strict) Logout(c *gin.Context) {
	ck, err := c.Cookie(session.DefaultCookieName)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	session.Cache.Delete(ck)
	c.Redirect(http.StatusPermanentRedirect, "/")
}
