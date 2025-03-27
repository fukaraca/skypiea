package handlers

import (
	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/google/uuid"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Login(c *gin.Context) {
	uid := "d653ea8b-7e0b-4a33-9910-85e393598c55"
	id := uuid.MustParse(uid)
	sess := session.Cache.New(id)
	session.Cache.Set(sess)
	c.SetCookie(session.DefaultCookieName, sess.ID, 100, "/", "localhost", false, true)
	c.HTML(http.StatusOK, "login", gin.H{
		"Title": "Login",
	})
}
