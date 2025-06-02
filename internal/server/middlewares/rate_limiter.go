package middlewares

import "github.com/gin-gonic/gin"

func RateLimiterMw() gin.HandlerFunc {
	return func(c *gin.Context) {
		// skip here for now and use AWS api gateway
		c.Next()
	}
}
