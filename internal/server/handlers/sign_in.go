package handlers

import (
	"fmt"
	"net/http"

	"github.com/fukaraca/skypiea/internal/model"
	"github.com/gin-gonic/gin"
)

func (h *View) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{
		"Title":    "Login",
		"LoggedIn": false,
	})
}

type SignInReq struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

func (h *Open) SignIn(c *gin.Context) {
	var in SignInReq
	if err := c.ShouldBind(&in); err != nil {
		h.AlertUI(c, err.Error(), ALError)
		return
	}

	sess, err := h.UserSvc.SignIn(c.Request.Context(), in.Email, in.Password)
	if err != nil {
		h.AlertUI(c, fmt.Sprintf("couldn't sign in due to: %v", err), ALError)
		return
	}
	c.SetCookie(sess.Name, sess.Value, sess.MaxAge, sess.Path, sess.Domain, sess.Secure, sess.HTTPOnly)
	c.Header(model.HxRedirect, "/")
	c.Status(http.StatusFound)
}
