package service

import (
	"context"
	"errors"

	"github.com/fukaraca/skypiea/internal/model"
	"github.com/fukaraca/skypiea/internal/server/middlewares"
	"github.com/fukaraca/skypiea/internal/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Service) ProcessNewMessage(ctx context.Context, userID uuid.UUID, message *storage.Message) (int, error) {
	logger := middlewares.GetLoggerFromContext(ctx)
	var err error
	if message.ConvID == 0 {
		// New conversation >> find some title
		err = s.Repositories.DoInTx(ctx, logger, func(reg *storage.Registry) error {
			message.ConvID, err = reg.Conversations.NewConversation(ctx, userID, uuid.New().String(), "")
			return err
		})
		if err != nil {
			logger.Error("new conversation couldn't be processed", err)
			return 0, model.ErrNewConversationCouldNotBeAdded
		}
	}

	err = s.Repositories.DoInTx(ctx, logger, func(reg *storage.Registry) error {
		message.ID, err = reg.Conversations.AppendNewMessage(ctx, message)
		return err
	})
	if err != nil {
		logger.Error("new message couldn't be appended", err)
	}
	return message.ID, model.ErrNewMessageCouldNotBeAdded
}

func (s *Service) GetAllMessages(ctx context.Context, convID int) ([]*storage.Message, error) {
	logger := middlewares.GetLoggerFromContext(ctx)
	messages := make([]*storage.Message, 0)
	messages, err := s.Repositories.Conversations.GetConversationByID(ctx, convID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		logger.Error("GetAllMessages failed", err)
		return nil, model.ErrMessagesCouldNotBeReloaded
	}
	return messages, err
}

func (s *Service) GetAllConversations(ctx context.Context, userID uuid.UUID) ([]*storage.Conversation, error) {
	logger := middlewares.GetLoggerFromContext(ctx)
	conversations := make([]*storage.Conversation, 0)
	conversations, err := s.Repositories.Conversations.GetConversationsByUserUUID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		logger.Error("GetAllConversations failed", err)
		return nil, model.ErrConversationCouldNotGet
	}
	return conversations, err
}

func (s *Service) DeleteConversation(ctx context.Context, convID int) error {
	return s.Repositories.DoInTx(ctx, middlewares.GetLoggerFromContext(ctx), func(reg *storage.Registry) error {
		return reg.Conversations.DeleteConversation(ctx, convID)
	})

}
