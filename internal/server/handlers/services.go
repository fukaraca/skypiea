package handlers

import (
	"context"

	"github.com/fukaraca/skypiea/internal/storage"
	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/google/uuid"
)

type UserService interface {
	SignIn(ctx context.Context, email, password string) (*session.Cookie, error)
	RegisterNewUser(ctx context.Context, user *storage.User) error
	ChangePassword(ctx context.Context, email, newPass string) error
	GetUser(ctx context.Context, userID uuid.UUID) (*storage.User, error)
}

type MessageService interface {
	ProcessNewMessage(ctx context.Context, userID uuid.UUID, message *storage.Message) (int, error)
	GetResponseByMessageID(ctx context.Context, userID uuid.UUID, msgID, convID int) (*storage.Message, error)
	GetAllMessages(ctx context.Context, convID int) ([]*storage.Message, error)
	GetMessage(ctx context.Context, msgID int) (*storage.Message, error)
	GetAllConversations(ctx context.Context, userID uuid.UUID) ([]*storage.Conversation, error)
	DeleteConversation(ctx context.Context, convID int) error
}
