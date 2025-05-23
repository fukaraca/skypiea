package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/fukaraca/skypiea/internal/model"
	"github.com/fukaraca/skypiea/internal/storage"
	"github.com/fukaraca/skypiea/pkg/gwt"
	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/gin-gonic/gin"
)

type SendMessageReq struct {
	MessageText    string `form:"messageText" binding:"required"`
	ConversationID int    `form:"conv_id"`
}

func (s *Strict) PostMessage(c *gin.Context) {
	var in SendMessageReq
	if err := c.ShouldBind(&in); err != nil {
		s.AlertUI(c, err.Error(), ALError)
		return
	}

	userID := session.Cache.GetUserUUIDByToken(c.GetString(gwt.CtxToken))
	if userID == nil {
		s.AlertUI(c, model.ErrSessionNotFound.Message, ALError)
		return
	}

	_, err := s.MessageSvc.ProcessNewMessage(c.Request.Context(), *userID, &storage.Message{
		ConvID:      in.ConversationID,
		ModelID:     "",
		ByUser:      true,
		MessageText: in.MessageText,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		s.AlertUI(c, err.Error(), ALError)
	}

	c.Status(http.StatusCreated)
}

func (s *Strict) GetMessagesByConversationID(c *gin.Context) {
	conversationID, err := strconv.Atoi(c.Param("conv_id"))
	if err != nil {
		s.AlertUI(c, err.Error(), ALError)
		return
	}

	messages, err := s.MessageSvc.GetAllMessages(c.Request.Context(), conversationID)
	if err != nil {
		s.AlertUI(c, err.Error(), ALError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Messages": messages,
	})
}

func (s *Strict) DeleteConversationByID(c *gin.Context) {
	conversationID, err := strconv.Atoi(c.Param("conv_id"))
	if err != nil {
		s.AlertUI(c, err.Error(), ALError)
		return
	}

	if err = s.MessageSvc.DeleteConversation(c.Request.Context(), conversationID); err != nil {
		s.AlertUI(c, err.Error(), ALError)
		return
	}
	c.Status(http.StatusNoContent)
}
