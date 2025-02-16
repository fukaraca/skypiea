package middlewares

import "github.com/gin-gonic/gin"

func RateLimiterMw() gin.HandlerFunc {
	return func(c *gin.Context) {
		// rate limiter logic
		c.Next()
	}
}
