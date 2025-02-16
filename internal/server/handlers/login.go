package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{
		"Title": "Login",
	})
}