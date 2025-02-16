package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Features(c *gin.Context) {
	c.HTML(http.StatusOK, "features", gin.H{
		"Title": "Features",
	})
}