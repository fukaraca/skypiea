package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/fukaraca/skypiea/internal/model"
	"github.com/fukaraca/skypiea/pkg/encryption"
	"github.com/fukaraca/skypiea/pkg/gwt"
	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/gin-gonic/gin"
)

type ChangePasswordReq struct {
	NewPassword     string `form:"new_password" binding:"required"`
	ConfirmPassword string `form:"confirm_password" binding:"required"`
}

func (h *View) ForgotPassword(c *gin.Context) {
	c.HTML(http.StatusOK, "forgot-password", gin.H{
		"Title": "Recover your password",
	})
}

func (s *Strict) ChangePassword(c *gin.Context) {
	userID := session.Cache.GetUserUUIDByToken(c.GetString(gwt.CtxToken))
	if userID == nil {
		c.Error(model.ErrSessionNotFound)
		return
	}
	var in ChangePasswordReq
	err := c.ShouldBind(&in)
	if err != nil {
		c.Error(err)
		return
	}
	if in.NewPassword == "" || in.NewPassword != in.ConfirmPassword {
		c.Error(model.ErrIncorrectCred)
		return
	}
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	pass, err := encryption.HashPassword(in.ConfirmPassword)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "failure", gin.H{
			"Title":         "Internal Error",
			"LoggedIn":      c.GetBool(session.CtxLoggedIn),
			"StatusCode":    http.StatusInternalServerError,
			"StatusHead":    "Password change not succeeded",
			"StatusMessage": err.Error(),
		})
		return
	}
	err = s.Repo.Users.ChangePassword(ctx, *userID, pass)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "failure", gin.H{
			"Title":         "Internal Error",
			"LoggedIn":      c.GetBool(session.CtxLoggedIn),
			"StatusCode":    http.StatusInternalServerError,
			"StatusHead":    "Password change not succeeded",
			"StatusMessage": err.Error(),
		})
		return
	}
	s.Logout(c)
}
