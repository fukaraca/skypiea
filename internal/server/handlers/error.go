package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Common) AlertUI(c *gin.Context, message string, level AlertLevel) {
	c.HTML(http.StatusOK, "alerts", gin.H{
		"AlertMessage": message,
		"AlertStyle":   level,
		"ShowAlert":    true,
	})
}
