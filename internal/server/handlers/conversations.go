package handlers

import (
	"fmt"
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
	Model          string `form:"llmodel"`
	ConversationID int    `form:"conv_id"`
}

func (s *Strict) PostMessage(c *gin.Context) {
	var in SendMessageReq
	if err := c.ShouldBind(&in); err != nil {
		s.AlertUI(c, err, ALError)
		return
	}

	userID := session.Cache.GetUserUUIDByToken(c.GetString(gwt.CtxToken))
	if userID == nil {
		s.AlertUI(c, model.ErrSessionNotFound, ALError)
		return
	}
	msg := &storage.Message{
		ConvID:      in.ConversationID,
		ModelID:     in.Model,
		ByUser:      true,
		MessageText: &in.MessageText,
		CreatedAt:   time.Now(),
	}

	msgID, err := s.MessageSvc.ProcessNewMessage(c.Request.Context(), *userID, msg)
	if err != nil {
		s.AlertUI(c, err, ALError)
		return
	}
	msg.ID = msgID

	c.Header("HX-Trigger", fmt.Sprintf(`{"chat:new":{"conv_id":%d}}`, msg.ConvID))
	c.HTML(http.StatusOK, "chat-post-post", gin.H{
		"Messages": []*storage.Message{msg},
		"MsgID":    msg.ID,
		"ConvID":   msg.ConvID,
	})
}

// ResponseOfMessage is complementary of post message
func (s *Strict) ResponseOfMessage(c *gin.Context) {
	msgID, err := strconv.Atoi(c.Param("msg_id"))
	if err != nil {
		s.AlertUI(c, err, ALError)
		return
	}
	convID, err := strconv.Atoi(c.Param("conv_id"))
	if err != nil {
		s.AlertUI(c, err, ALError)
		return
	}
	userID := session.Cache.GetUserUUIDByToken(c.GetString(gwt.CtxToken))
	if userID == nil {
		s.AlertUI(c, model.ErrSessionNotFound, ALError)
		return
	}
	resp, err := s.MessageSvc.GetResponseByMessageID(c.Request.Context(), *userID, msgID, convID)
	if err != nil {
		s.AlertUI(c, err, ALError)
		return
	}
	c.HTML(http.StatusOK, "chat-panel", gin.H{
		"Messages": []*storage.Message{resp},
	})
}

func (s *Strict) GetMessagesByConversationID(c *gin.Context) {
	conversationID, err := strconv.Atoi(c.Param("conv_id"))
	if err != nil {
		s.AlertUI(c, err, ALError)
		return
	}

	messages, err := s.MessageSvc.GetAllMessages(c.Request.Context(), conversationID)
	if err != nil {
		s.AlertUI(c, err, ALError)
		return
	}
	c.HTML(http.StatusOK, "chat-panel", gin.H{
		"Messages": messages,
	})
}

func (s *Strict) DeleteConversationByID(c *gin.Context) {
	conversationID, err := strconv.Atoi(c.Param("conv_id"))
	if err != nil {
		s.AlertUI(c, err, ALError)
		return
	}

	if err = s.MessageSvc.DeleteConversation(c.Request.Context(), conversationID); err != nil {
		s.AlertUI(c, err, ALError)
		return
	}
	c.Status(http.StatusNoContent)
	c.Header(HX_REDIRECT, model.PathMain)
}
