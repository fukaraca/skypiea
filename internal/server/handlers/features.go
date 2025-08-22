package handlers

import (
	"net/http"

	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/gin-gonic/gin"
)

func (h *View) Features(c *gin.Context) {
	c.HTML(http.StatusOK, "features", gin.H{
		"Title":    "Features",
		"CSSFile":  "index.css",
		"LoggedIn": c.GetBool(session.CtxLoggedIn),
	})
}
