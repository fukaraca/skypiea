package handlers

import (
	"net/http"

	"github.com/fukaraca/skypiea/internal/model"
	"github.com/fukaraca/skypiea/pkg/gwt"
	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/gin-gonic/gin"
)

func (h *View) Index(c *gin.Context) {
	loggedIn := c.GetBool(session.CtxLoggedIn)
	if !loggedIn {
		c.HTML(http.StatusOK, "index", gin.H{
			"Title":    "Home",
			"CSSFile":  "index.css",
			"LoggedIn": false,
		})
		return
	}
	userID := session.Cache.GetUserUUIDByToken(c.GetString(gwt.CtxToken))
	if userID == nil {
		h.AlertUI(c, model.ErrSessionNotFound, ALError)
		return
	}
	convs, err := h.MessageSvc.GetAllConversations(c.Request.Context(), *userID)
	if err != nil {
		h.AlertUI(c, model.ErrConversationCouldNotGet, ALError)
		return
	}

	c.HTML(http.StatusOK, "index", gin.H{
		"Title":         "Home",
		"CSSFile":       "index.css",
		"LoggedIn":      true,
		"Conversations": convs,
		"CurrentConv":   nil,
		"Messages":      nil,
	})
}
