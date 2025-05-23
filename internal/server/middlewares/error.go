package middlewares

import (
	"github.com/gin-gonic/gin"
)

// ErrorHandlerMw used - at least was planned - for processing errors on exit, logging maybe saving
func ErrorHandlerMw() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			GetLoggerFromContext(c).Error("request has failed", "errors", c.Errors)
		}
	}
}
