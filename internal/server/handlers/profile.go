package handlers

import (
	"net/http"

	"github.com/fukaraca/skypiea/internal/model"
	"github.com/fukaraca/skypiea/pkg/gwt"
	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/gin-gonic/gin"
)

func (h *View) Profile(c *gin.Context) {
	loggedIn := c.GetBool(session.CtxLoggedIn)
	if !loggedIn {
		c.Header(model.HxRedirect, "/")
		c.Status(http.StatusFound)
		c.Abort()
		return
	}
	ctx, cancel := h.CtxWithTimout(c)
	defer cancel()
	userID := session.Cache.GetUserUUIDByToken(c.GetString(gwt.CtxToken))
	if userID == nil {
		c.HTML(http.StatusInternalServerError, "failure", gin.H{
			"Title":         "Internal Error",
			"LoggedIn":      c.GetBool(session.CtxLoggedIn),
			"StatusCode":    http.StatusInternalServerError,
			"StatusHead":    "Request not succeeded",
			"StatusMessage": "Token not found",
		})
		return
	}
	u, err := h.Repo.Users.GetUserByUUID(ctx, *userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "failure", gin.H{
			"Title":         "Internal Error",
			"LoggedIn":      c.GetBool(session.CtxLoggedIn),
			"StatusCode":    http.StatusInternalServerError,
			"StatusHead":    "Request not succeeded",
			"StatusMessage": "User could not fetched",
		})
		return
	}
	c.HTML(http.StatusOK, "profile", gin.H{
		"Title":    "My Profile",
		"LoggedIn": c.GetBool(session.CtxLoggedIn),
		"Name":     u.Firstname,
		"Surname":  u.Lastname,
	})
}
