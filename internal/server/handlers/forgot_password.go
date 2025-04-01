package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *View) ForgotPassword(c *gin.Context) {
	c.HTML(http.StatusOK, "forgot-password", gin.H{
		"Title": "Recover your password",
	})
}
