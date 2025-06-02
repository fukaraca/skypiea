package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/fukaraca/skypiea/internal/model"
	"github.com/fukaraca/skypiea/internal/server/middlewares"
	"github.com/fukaraca/skypiea/internal/storage"
	"github.com/fukaraca/skypiea/pkg/gwt"
	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Service) ProcessNewMessage(ctx context.Context, userID uuid.UUID, message *storage.Message) (int, error) {
	logger := middlewares.GetLoggerFromContext(ctx)
	var err error
	if message.ConvID == 0 {
		// New conversation >> find some title
		err = s.Repositories.DoInTx(ctx, logger, func(reg *storage.Registry) error {
			message.ConvID, err = reg.Conversations.NewConversation(ctx, userID, TitleFromString(*message.MessageText, 2, 16), "")
			return err
		})
		if err != nil {
			logger.Error("new conversation couldn't be processed", slog.Any("error", err))
			return 0, model.ErrNewConversationCouldNotBeAdded.WithError(err)
		}
	}
	ok, err := s.Repositories.Conversations.VerifyUserForConversation(ctx, userID, message.ConvID) // Maybe this should be done somewhere else, doesnt scale
	switch {
	case err != nil:
		return 0, err
	case !ok:
		return 0, model.ErrConversationCouldNotGet.WithError(errors.New("conversation is not for this user"))
	}

	err = s.Repositories.DoInTx(ctx, logger, func(reg *storage.Registry) error {
		message.ID, err = reg.Conversations.AppendNewMessage(ctx, message)
		if err == nil {
			err = reg.Conversations.BumpConversationUpdatedAtV(ctx, message.ConvID)
		}
		return err
	})
	if err != nil {
		logger.Error("new message couldn't be appended", slog.Any("error", err))
		return 0, model.ErrNewMessageCouldNotBeAdded.WithError(err)
	}
	return message.ID, nil
}

func (s *Service) GetAllMessages(ctx context.Context, convID int) ([]*storage.Message, error) {
	logger := middlewares.GetLoggerFromContext(ctx)
	userID := session.Cache.GetUserUUIDByToken(middlewares.GetGinCtxFromContext(ctx).GetString(gwt.CtxToken))

	ok, err := s.Repositories.Conversations.VerifyUserForConversation(ctx, *userID, convID)
	switch {
	case err != nil:
		return nil, err
	case !ok:
		return nil, model.ErrConversationCouldNotGet.WithError(errors.New("conversation is not for this user"))
	}

	messages, err := s.Repositories.Conversations.GetConversationByID(ctx, convID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		logger.Error("GetAllMessages failed", slog.Any("error", err))
		return nil, model.ErrMessagesCouldNotBeReloaded.WithError(err)
	}
	return messages, err
}

func (s *Service) GetMessage(ctx context.Context, msgID int) (*storage.Message, error) {
	logger := middlewares.GetLoggerFromContext(ctx)
	userID := session.Cache.GetUserUUIDByToken(middlewares.GetGinCtxFromContext(ctx).GetString(gwt.CtxToken))
	ok, err := s.Repositories.Conversations.VerifyUserForMessage(ctx, *userID, msgID)
	switch {
	case err != nil:
		return nil, err
	case !ok:
		return nil, model.ErrMessageCouldNotGet.WithError(errors.New("message is not for this user"))
	}

	message, err := s.Repositories.Conversations.GetMessageByID(ctx, msgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return message, nil
		}
		logger.Error("GetAllMessages failed", slog.Any("error", err))
		return nil, model.ErrMessagesCouldNotBeReloaded.WithError(err)
	}
	return message, err
}

// GetResponseByMessageID is a temporary BS(#!$%) method instead of maintaining an API to GPT
func (s *Service) GetResponseByMessageID(ctx context.Context, userID uuid.UUID, msgID, convID int) (*storage.Message, error) {
	ok, err := s.Repositories.Conversations.VerifyUserForConversation(ctx, userID, convID)
	switch {
	case err != nil:
		return nil, err
	case !ok:
		return nil, model.ErrMessageCouldNotGet.WithError(errors.New("conversation is not for this user"))
	}

	message, err := s.Repositories.Conversations.GetResponseByID(ctx, msgID)
	switch {
	case err == nil:
		return message, nil
	case !errors.Is(err, pgx.ErrNoRows):
		return nil, err
	}

	time.Sleep(2 * time.Second)
	// take it like i'm thinking hard
	msgTxt := fmt.Sprintf("This is auto generated response for message %d.\nLets see maybe we link it to Gemini or something cheaper in the future just for fun.. ", msgID) //nolint: lll
	message = &storage.Message{
		ConvID:      convID,
		ModelID:     model.GPTFurkan,
		ByUser:      false,
		MessageText: &msgTxt,
		Metadata:    nil,
		CreatedAt:   time.Now(),
		ResponseTo:  &msgID,
	}

	message.ID, err = s.ProcessNewMessage(ctx, userID, message)
	return message, err
}

func (s *Service) GetAllConversations(ctx context.Context, userID uuid.UUID) ([]*storage.Conversation, error) {
	logger := middlewares.GetLoggerFromContext(ctx)
	conversations, err := s.Repositories.Conversations.GetConversationsByUserUUID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		logger.Error("GetAllConversations failed", slog.Any("error", err))
		return nil, model.ErrConversationCouldNotGet.WithError(err)
	}
	return conversations, err
}

func (s *Service) DeleteConversation(ctx context.Context, convID int) error {
	return s.Repositories.DoInTx(ctx, middlewares.GetLoggerFromContext(ctx), func(reg *storage.Registry) error {
		return reg.Conversations.DeleteConversation(ctx, convID)
	})
}

func TitleFromString(input string, maxWords, maxLen int) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return fmt.Sprintf("untitled-%s", uuid.New().String()[:8])
	}

	words := strings.Fields(input)
	var selected []string
	currLen := 0

	for i, w := range words {
		if i >= maxWords {
			break
		}

		addLen := len(w)
		if len(selected) > 0 {
			addLen++
		}
		if currLen+addLen > maxLen {
			break
		}
		selected = append(selected, w)
		currLen += addLen
	}

	if len(selected) > 0 {
		return strings.Join(selected, " ")
	}

	// fallback for giant first word: take a prefix of the raw input
	truncated := input
	if len(truncated) > maxLen {
		truncated = truncated[:maxLen]
		truncated = strings.TrimSpace(truncated)
	}
	return fmt.Sprintf("%s-%s", truncated, uuid.New().String()[:8])
}
