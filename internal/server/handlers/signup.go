package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Signup(c *gin.Context) {
	c.HTML(http.StatusOK, "signup", gin.H{
		"Title": "Sign Up",
	})
}