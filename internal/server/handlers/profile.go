package handlers

import (
	"net/http"

	"github.com/fukaraca/skypiea/internal/model"
	"github.com/fukaraca/skypiea/internal/storage"
	"github.com/fukaraca/skypiea/pkg/gwt"
	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/gin-gonic/gin"
)

func (h *View) Profile(c *gin.Context) {
	loggedIn := c.GetBool(session.CtxLoggedIn)
	if !loggedIn {
		c.Redirect(http.StatusFound, model.PathMain)
		c.Abort()
		return
	}

	userID := session.Cache.GetUserUUIDByToken(c.GetString(gwt.CtxToken))
	if userID == nil {
		c.HTML(http.StatusInternalServerError, "failure", gin.H{
			"Title":         "Internal Error",
			"LoggedIn":      loggedIn,
			"StatusCode":    http.StatusInternalServerError,
			"StatusHead":    "Request not succeeded",
			"StatusMessage": "Token not found",
		})
		return
	}
	u, err := h.UserSvc.GetUser(c.Request.Context(), *userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "failure", gin.H{
			"Title":         "Internal Error",
			"LoggedIn":      loggedIn,
			"StatusCode":    http.StatusInternalServerError,
			"StatusHead":    "Request not succeeded",
			"StatusMessage": "User could not fetched",
		})
		return
	}
	c.HTML(http.StatusOK, "profile", gin.H{
		"Title":       "My Profile",
		"LoggedIn":    loggedIn,
		"Name":        u.Firstname,
		"Surname":     u.Lastname,
		"PhoneNumber": u.PhoneNumber.String,
		"AboutMe":     u.AboutMe.String,
	})
}

type UpdateProfileReq struct {
	Name        string `form:"name" binding:"required"`
	Surname     string `form:"surname" binding:"required"`
	PhoneNumber string `form:"phone"`
	AboutMe     string `form:"about_me"`
}

func (s *Strict) ProfileUpdate(c *gin.Context) {
	userID := session.Cache.GetUserUUIDByToken(c.GetString(gwt.CtxToken))
	if userID == nil {
		s.AlertUI(c, model.ErrSessionNotFound, ALError)
		return
	}
	var in UpdateProfileReq
	err := c.ShouldBind(&in)
	if err != nil {
		s.AlertUI(c, err, ALError)
		return
	}
	u := &storage.User{
		UserUUID:  userID.String(),
		Firstname: in.Name,
		Lastname:  in.Surname,
	}
	u.PhoneNumber.Scan(in.PhoneNumber)
	u.AboutMe.Scan(in.AboutMe)

	err = s.UserSvc.UpdateUserProfile(c.Request.Context(), u)
	if err != nil {
		s.AlertUI(c, err, ALError)
		return
	}
	s.AlertUI(c, "Profile updated successfully", ALInfo)
}
