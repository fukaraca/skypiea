package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Open) Healthz(c *gin.Context) {
	// TODO some internal check like ping to DB etc...
	c.JSON(http.StatusOK, gin.H{"ping": "OK"})
}
