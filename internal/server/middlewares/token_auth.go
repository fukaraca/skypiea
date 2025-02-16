package middlewares

import "github.com/gin-gonic/gin"

func TokenAuthMw() gin.HandlerFunc {
	return func(c *gin.Context) {
		// token auth logic
		c.Next()
	}
}
