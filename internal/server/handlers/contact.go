package handlers

import (
	"net/http"

	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/gin-gonic/gin"
)

func (h *View) Contact(c *gin.Context) {
	c.HTML(http.StatusOK, "contact", gin.H{
		"Title":    "Contact",
		"CSSFile":  "index.css",
		"LoggedIn": c.GetBool(session.CtxLoggedIn),
	})
}
