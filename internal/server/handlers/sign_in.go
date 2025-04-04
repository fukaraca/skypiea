package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

func (h *Common) SignIn(c *gin.Context) {
	var in SignInReq
	if err := c.ShouldBind(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"loginResponse": err.Error(),
		})
		return
	}
	ctx, cancel := context.WithTimeout(c, time.Second*100)
	defer cancel()
	user, err := h.Repo.Users.GetUserByEmail(ctx, in.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"loginResponse": fmt.Sprintf("user not found for %q", in.Email),
		})
		return
	}
	if !encryption.CheckPasswordHash(in.Password, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"loginResponse": model.ErrIncorrectCred.Message,
		})
		return
	}

	sess := session.Cache.NewSession(ctx, uuid.MustParse(user.UserUUID))
	session.Cache.Set(sess)
	c.SetCookie(session.DefaultCookieName, sess.ID, 100, "/", "localhost", false, true)
	c.Header(model.HxRedirect, "/")
	c.Status(http.StatusFound)
}
