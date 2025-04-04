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
		c.HTML(http.StatusInternalServerError, "failure", gin.H{
			"Title":         "Internal Error",
			"LoggedIn":      c.GetBool(session.CtxLoggedIn),
			"StatusCode":    http.StatusInternalServerError,
			"StatusHead":    "Logout not succeeded",
			"StatusMessage": err.Error(),
		})
		return
	}
	session.Cache.Delete(ck)
	c.Header(model.HxRedirect, "/")
	c.Status(http.StatusFound)
}
