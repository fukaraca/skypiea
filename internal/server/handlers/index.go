package handlers

import (
	"fmt"
	"net/http"

	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/gin-gonic/gin"
)

func (h *View) Index(c *gin.Context) {
	A := "ASDASD"
	fmt.Println(A)

	c.HTML(http.StatusOK, "index", gin.H{
		"Title":    "Home",
		"CSSFile":  "index.css",
		"LoggedIn": c.GetBool(session.CtxLoggedIn),
	})
}
