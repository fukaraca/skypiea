package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *View) Profile(c *gin.Context) {
	c.HTML(http.StatusOK, "profile", gin.H{
		"Title": "My Profile",
	})
}
