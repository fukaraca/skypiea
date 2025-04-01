package handlers

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/fukaraca/skypiea/internal/server/middlewares"
)

func (h *View) Features(c *gin.Context) {
	middlewares.GetLoggerFromContext(c).Debug("features")
	c.AbortWithError(400, os.ErrClosed).SetMeta("features")
}
