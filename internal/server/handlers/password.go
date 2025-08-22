package handlers

import (
	"net/http"

	"github.com/fukaraca/skypiea/internal/model"
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
		"Title":    "Recover your password",
		"CSSFile":  "index.css",
		"LoggedIn": c.GetBool(session.CtxLoggedIn),
	})
}

func (s *Strict) ChangePassword(c *gin.Context) {
	userID := session.Cache.GetUserUUIDByToken(c.GetString(gwt.CtxToken))
	if userID == nil {
		s.AlertUI(c, model.ErrSessionNotFound, ALError)
		return
	}
	var in ChangePasswordReq
	err := c.ShouldBind(&in)
	if err != nil {
		s.AlertUI(c, err, ALError)
		return
	}
	if in.NewPassword == "" || in.NewPassword != in.ConfirmPassword {
		s.AlertUI(c, model.ErrIncorrectCred, ALError)
		return
	}
	ctx := c.Request.Context()
	user, err := s.UserSvc.GetUser(ctx, *userID)
	if err != nil {
		s.AlertUI(c, err, ALError)
		return
	}
	if err = s.UserSvc.ChangePassword(c.Request.Context(), user.Email, in.NewPassword); err != nil {
		s.AlertUI(c, err, ALError)
		return
	}
	s.Logout(c)
}

type ForgotPasswordReq struct {
	Email string `form:"email" binding:"required"`
}

// TODO a proper password change logic
func (h *Open) ForgotPassword(c *gin.Context) {
	var in ForgotPasswordReq
	if err := c.ShouldBind(&in); err != nil {
		h.AlertUI(c, err, ALError)
		return
	}
	if err := h.UserSvc.ChangePassword(c.Request.Context(), in.Email, "iForgotMyPassword"); err != nil {
		h.AlertUI(c, err, ALError)
		return
	}
	h.AlertUI(c, "New Password sent to your email... (kidding it's 'iForgotMyPassword' now, go and try now).", ALInfo)
}
