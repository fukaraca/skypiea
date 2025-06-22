package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/fukaraca/skypiea/internal/model"
	"github.com/fukaraca/skypiea/internal/server/middlewares"
	"github.com/fukaraca/skypiea/internal/storage"
	"github.com/fukaraca/skypiea/pkg/gwt"
	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const summarizer = `You are a summarization engine. Read the following text and return only a single sentence that serves as a conversation title.
Each word must start with a capital letter. Do not include any punctuation marks at all. Only return the title, nothing else.
Text:
%s`

func (s *Service) ProcessNewMessage(ctx context.Context, userID uuid.UUID, message *storage.Message) (int, error) {
	logger := middlewares.GetLoggerFromContext(ctx)
	var err error
	if message.ConvID == 0 {
		// New conversation >> find some title
		title, err := s.GenerateTitle(ctx, string(*message.MessageText)) // nolint:govet
		if err != nil {
			title = TitleFromString(string(*message.MessageText), 2, 16)
		}
		err = s.Repositories.DoInTx(ctx, logger, func(reg *storage.Registry) error {
			message.ConvID, err = reg.Conversations.NewConversation(ctx, userID, title, "")
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

// GetResponseByMessageID returns the response from either DB or asistant.
func (s *Service) GetResponseByMessageID(ctx context.Context, userID uuid.UUID, msgID, convID int) (*storage.Message, error) {
	ok, err := s.Repositories.Conversations.VerifyUserForConversation(ctx, userID, convID)
	switch {
	case err != nil:
		return nil, err
	case !ok:
		return nil, model.ErrMessageCouldNotGet.WithError(errors.New("conversation is not for this user"))
	}

	found, err := s.Repositories.Conversations.GetResponseByID(ctx, msgID)
	switch {
	case err == nil:
		return found, nil
	case !errors.Is(err, pgx.ErrNoRows):
		return nil, err
	}
	message, err := s.Repositories.Conversations.GetMessageByID(ctx, msgID)
	if err != nil {
		return nil, err
	}

	history, err := s.Repositories.Conversations.GetConversationByID(ctx, message.ConvID)
	if err != nil {
		return nil, err
	}

	u, err := s.Repositories.Users.GetUserByUUID(ctx, userID)
	if err != nil {
		return nil, err
	}

	question := BuildPrompt(&u.AboutMe.String, &u.Summary.String, history, message)

	response, err := s.GeminiClient.AskToGemini(ctx, question, message.ModelID)
	if err != nil {
		return nil, err
	}

	messageToSave := &storage.Message{
		ConvID:      convID,
		ModelID:     message.ModelID,
		ByUser:      false,
		MessageText: s.Sanitize(response, false),
		Metadata:    nil,
		CreatedAt:   time.Now(),
		ResponseTo:  &msgID,
	}

	messageToSave.ID, err = s.ProcessNewMessage(ctx, userID, messageToSave)
	return messageToSave, err
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

func (s *Service) GenerateTitle(ctx context.Context, msg string) (string, error) {
	msg = fmt.Sprintf(summarizer, msg)
	title, err := s.GeminiClient.AskToGemini(ctx, msg, "")
	return title, err
}
