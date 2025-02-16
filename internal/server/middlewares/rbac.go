package middlewares

import "github.com/gin-gonic/gin"

func RBACMw() gin.HandlerFunc {
	return func(c *gin.Context) {
		// rbac logic
		c.Next()
	}
}
