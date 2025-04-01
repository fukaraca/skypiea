package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/fukaraca/skypiea/internal/storage"
	"github.com/gin-gonic/gin"
)

func (h *View) Signup(c *gin.Context) {
	c.HTML(http.StatusOK, "signup", gin.H{
		"Title": "Sign Up",
	})
}

type SignUpReq struct {
	Password  string `form:"password" binding:"required"`
	Firstname string `form:"firstName" binding:"required"`
	Lastname  string `form:"lastName" binding:"required"`
	Email     string `form:"email" binding:"required,email"`
}

func (h *Common) SignUp(c *gin.Context) {
	var in SignUpReq
	if err := c.ShouldBind(&in); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	user := &storage.User{
		Firstname: in.Firstname,
		Lastname:  in.Lastname,
		Email:     in.Email,
		Password:  in.Password,
		Role:      "admin",
		Status:    "New",
	}
	_, err := h.Repo.Users.AddUser(ctx, user)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}
