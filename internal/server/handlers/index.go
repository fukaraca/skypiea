package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Index(c *gin.Context) {
	A := "ASDASD"
	fmt.Println(A)
	c.HTML(http.StatusOK, "index", gin.H{
		"Title":    "Home",
		"CSSFile":  "index.css",
		"LoggedIn": false,
	})
}
