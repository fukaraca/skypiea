package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *View) About(c *gin.Context) {
	c.HTML(http.StatusOK, "about", gin.H{
		"Title": "About",
	})
}
