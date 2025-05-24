package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Common) AlertUI(c *gin.Context, message any, level AlertLevel) {
	if level == ALError {
		v, ok := message.(error)
		if ok {
			c.Error(v)
			message = v.Error()
		}
	}
	c.HTML(http.StatusOK, "alerts", gin.H{
		"AlertMessage": message,
		"AlertStyle":   level,
		"ShowAlert":    true,
	})
}
