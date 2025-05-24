package handlers

import (
	"net/http"

	"github.com/fukaraca/skypiea/internal/model"
	"github.com/fukaraca/skypiea/internal/storage"
	"github.com/gin-gonic/gin"
)

func (h *View) Signup(c *gin.Context) {
	c.HTML(http.StatusOK, "signup", gin.H{
		"Title":    "Sign Up",
		"LoggedIn": false,
	})
}

type SignUpReq struct {
	Password  string `form:"password" binding:"required"`
	Firstname string `form:"firstName" binding:"required"`
	Lastname  string `form:"lastName" binding:"required"`
	Email     string `form:"email" binding:"required,email"`
	Role      string `form:"role" binding:"-"`
}

// TODO use passed role
func (h *Open) SignUp(c *gin.Context) {
	var in SignUpReq
	if err := c.ShouldBind(&in); err != nil {
		h.AlertUI(c, err, ALError)
		return
	}
	user := &storage.User{
		Firstname: in.Firstname,
		Lastname:  in.Lastname,
		Email:     in.Email,
		Password:  in.Password,
		Role:      model.RoleAdmin,
		Status:    model.StatusNew,
	}

	err := h.UserSvc.RegisterNewUser(c.Request.Context(), user)
	if err != nil {
		h.AlertUI(c, err, ALError)
		return
	}
	c.Header(model.HxRedirect, model.PathLogin)
	c.Status(http.StatusCreated)
}
