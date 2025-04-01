package handlers

import (
	"net/http"

	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/gin-gonic/gin"
)

func NoRoute404(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404", gin.H{
		"Title":    "Page not exist",
		"LoggedIn": c.GetBool(session.CtxLoggedIn),
	})
}
