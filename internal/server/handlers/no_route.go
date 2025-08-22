package handlers

import (
	"net/http"

	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/gin-gonic/gin"
)

func NoRoute404(c *gin.Context) {
	c.HTML(http.StatusNotFound, "failure", gin.H{
		"Title":         "Page not exist",
		"CSSFile":       "index.css",
		"LoggedIn":      c.GetBool(session.CtxLoggedIn),
		"StatusCode":    http.StatusNotFound,
		"StatusHead":    "Page not found.",
		"StatusMessage": "The page you’re looking for doesn’t exist or an error occurred.",
	})
}
