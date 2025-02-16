package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{
		"Title":    "Home",
		"CSSFile":  "index.css",
		"LoggedIn": false,
	})
}
