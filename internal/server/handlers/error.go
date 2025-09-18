package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Common) AlertUI(c *gin.Context, message any, level AlertLevel, optStatus ...int) {
	if level == ALError {
		v, ok := message.(error)
		if ok {
			c.Error(v)
			message = v.Error()
		}
	}
	status := http.StatusOK
	if len(optStatus) > 0 {
		status = optStatus[0]
	}
	c.HTML(status, "alerts", gin.H{
		"AlertMessage": message,
		"AlertStyle":   level,
		"ShowAlert":    true,
	})
}
