package middlewares

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/fukaraca/skypiea/internal/model"
)

func ErrorHandlerMw() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				var err *model.Error
				if !errors.As(e.Err, &err) {
					c.JSON(-1, e)
				} else {
					c.JSON(-1, e.Err)
				}
			}
			if l := GetLoggerFromContext(c); l != nil {
				l.Error("error handled", "errors", c.Errors)
			}
		}
	}
}
