package handlers

import (
	"fmt"
	"net/http"

	"github.com/fukaraca/skypiea/internal/model"
	"github.com/fukaraca/skypiea/pkg/encryption"
	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		h.AlertUI(c, err.Error(), AlertLevelError)
		return
	}
	ctx, cancel := h.CtxWithTimout(c)
	defer cancel()
	user, err := h.Repo.Users.GetUserByEmail(ctx, in.Email)
	if err != nil {
		h.AlertUI(c, fmt.Sprintf("user not found for %q", in.Email), AlertLevelError)
		return
	}
	if !encryption.CheckPasswordHash(in.Password, user.Password) {
		h.AlertUI(c, model.ErrIncorrectCred.Message, AlertLevelError)
		return
	}

	sess := session.Cache.NewSession(ctx, uuid.MustParse(user.UserUUID))
	session.Cache.Set(sess)
	c.SetCookie(session.DefaultCookieName, sess.ID, session.DefaultCookieMaxAge, "/", session.DefaultCookieDomain, false, true)
	c.Header(model.HxRedirect, "/")
	c.Status(http.StatusFound)
}
