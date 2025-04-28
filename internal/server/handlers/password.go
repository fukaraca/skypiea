package handlers

import (
	"fmt"
	"net/http"

	"github.com/fukaraca/skypiea/internal/model"
	"github.com/fukaraca/skypiea/pkg/encryption"
	"github.com/fukaraca/skypiea/pkg/gwt"
	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		s.AlertUI(c, model.ErrSessionNotFound.Message, AlertLevelError)
		return
	}
	var in ChangePasswordReq
	err := c.ShouldBind(&in)
	if err != nil {
		s.AlertUI(c, err.Error(), AlertLevelError)
		return
	}
	if in.NewPassword == "" || in.NewPassword != in.ConfirmPassword {
		s.AlertUI(c, model.ErrIncorrectCred.Message, AlertLevelError)
		return
	}
	ctx, cancel := s.CtxWithTimout(c)
	defer cancel()
	pass, err := encryption.HashPassword(in.ConfirmPassword)
	if err != nil {
		s.AlertUI(c, fmt.Sprintf("password change not succeeded: %v", err), AlertLevelError)
		return
	}
	err = s.Repo.Users.ChangePassword(ctx, *userID, pass)
	if err != nil {
		s.AlertUI(c, fmt.Sprintf("password change not succeeded: %v", err), AlertLevelError)
		return
	}
	s.Logout(c)
}

type ForgotPasswordReq struct {
	Email string `form:"email" binding:"required"`
}

func (h *Open) ForgotPassword(c *gin.Context) {
	var in ForgotPasswordReq
	if err := c.ShouldBind(&in); err != nil {
		h.AlertUI(c, err.Error(), AlertLevelError)
		return
	}
	ctx, cancel := h.CtxWithTimout(c)
	defer cancel()
	user, err := h.Repo.Users.GetUserByEmail(ctx, in.Email)
	if err != nil {
		h.AlertUI(c, err.Error(), AlertLevelError)
		return
	}
	hPass, err := encryption.HashPassword("iForgotMyPassword")
	if err != nil {
		h.AlertUI(c, err.Error(), AlertLevelError)
		return
	}
	err = h.Repo.Users.ChangePassword(ctx, uuid.MustParse(user.UserUUID), hPass)
	if err != nil {
		h.AlertUI(c, err.Error(), AlertLevelError)
		return
	}
	h.AlertUI(c, "New Password sent to your email... (kidding its iForgotMyEmail now, go and try now).", AlertLevelInfo)
}
