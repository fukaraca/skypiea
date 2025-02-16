package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ForgotPassword(c *gin.Context) {
	c.HTML(http.StatusOK, "forgot-password", gin.H{
		"Title": "Recover your password",
	})
}
