package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Contact(c *gin.Context) {
	c.HTML(http.StatusOK, "contact", gin.H{
		"Title": "Contact",
	})
}
