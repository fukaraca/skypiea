package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) FAQ(c *gin.Context) {
	c.HTML(http.StatusOK, "faq", gin.H{
		"Title": "FAQ",
	})
}
