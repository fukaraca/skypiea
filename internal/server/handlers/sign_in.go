package handlers

import (
	"errors"
	"net/http"

	"github.com/fukaraca/skypiea/internal/model"
	"github.com/gin-gonic/gin"
)

func (h *View) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{
		"Title":    "Login",
		"LoggedIn": false,
		"CSSFile":  "index.css",
	})
}

type SignInReq struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

func (h *Open) SignIn(c *gin.Context) {
	var in SignInReq
	if err := c.ShouldBind(&in); err != nil {
		h.AlertUI(c, err, ALError)
		return
	}

	sess, err := h.UserSvc.SignIn(c.Request.Context(), in.Email, in.Password)
	if err != nil {
		h.AlertUI(c, err, ALError)
		return
	}
	c.SetCookie(sess.Name, sess.Value, sess.MaxAge, sess.Path, sess.Domain, sess.Secure, sess.HTTPOnly)
	c.Header(model.HxRedirect, model.PathMain)
	c.Status(http.StatusFound)
}

func (h *Open) OAuth2Start(c *gin.Context) {
	isSignup := c.Query("flow") == "signup"

	provider := c.Param("provider")
	if provider != "google" {
		h.AlertUI(c, errors.New("provider is not supported"), ALError)
		return
	}
	redirect := h.Auth.Start(c, isSignup)
	c.Header(model.HxRedirect, redirect)
	c.Header("Access-Control-Expose-Headers", model.HxRedirect)
	c.Header("Access-Control-Allow-Origin", h.origin)
	c.String(http.StatusOK, "Redirecting")
}

func (h *Open) OAuth2Callback(c *gin.Context) {
	errorResp := c.Query("error")
	code := c.Query("code")
	state := c.Query("state")
	if errorResp != "" || code == "" || state == "" {
		h.AlertUI(c, errors.New("authentication failed"), ALError, http.StatusBadRequest)
		return
	}

	provider := c.Param("provider")
	if provider != "google" {
		h.AlertUI(c, errors.New("provider is not supported"), ALError)
		return
	}
	sess, err := h.Auth.Callback(c, code, state)
	if err != nil {
		h.AlertUI(c, err, ALError, http.StatusUnauthorized)
		return
	}
	c.SetCookie(sess.Name, sess.Value, sess.MaxAge, sess.Path, sess.Domain, sess.Secure, sess.HTTPOnly)
	c.Redirect(http.StatusSeeOther, model.PathMain)
	h.AlertUI(c, "Successfully signed in", ALInfo)
}
