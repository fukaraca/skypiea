package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Profile(c *gin.Context) {
	c.HTML(http.StatusOK, "profile", gin.H{
		"Title": "My Profile",
	})
}
