package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *View) Pricing(c *gin.Context) {
	c.HTML(http.StatusOK, "pricing", gin.H{
		"Title": "Pricing",
	})
}
