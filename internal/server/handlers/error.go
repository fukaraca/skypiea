package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Common) AlertUI(c *gin.Context, message any, level AlertLevel) {
	if level == ALError {
		c.Error(message.(error))
	}
	c.HTML(http.StatusOK, "alerts", gin.H{
		"AlertMessage": message,
		"AlertStyle":   level,
		"ShowAlert":    true,
	})
}
