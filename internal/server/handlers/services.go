package handlers

import (
	"context"
	"html/template"

	"github.com/fukaraca/skypiea/internal/storage"
	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/google/uuid"
)

type UserService interface {
	SignIn(ctx context.Context, email, password string) (*session.Cookie, error)
	RegisterNewUser(ctx context.Context, user *storage.User) error
	ChangePassword(ctx context.Context, email, newPass string) error
	UpdateRole(ctx context.Context, userUUID, role string) error
	GetUser(ctx context.Context, userID uuid.UUID) (*storage.User, error)
	SupportedModels(ctx context.Context, userID uuid.UUID) []string
	UpdateUserProfile(ctx context.Context, userNew *storage.User) error
	GetAllUsers(ctx context.Context) ([]*storage.User, error)
	GetAdoptionStatistics(ctx context.Context) ([]*storage.AdoptionStat, error)
}

type MessageService interface {
	ProcessNewMessage(ctx context.Context, userID uuid.UUID, message *storage.Message) (int, error)
	GetResponseByMessageID(ctx context.Context, userID uuid.UUID, msgID, convID int) (*storage.Message, error)
	GetAllMessages(ctx context.Context, convID int) ([]*storage.Message, error)
	GetMessage(ctx context.Context, msgID int) (*storage.Message, error)
	GetAllConversations(ctx context.Context, userID uuid.UUID) ([]*storage.Conversation, error)
	DeleteConversation(ctx context.Context, convID int) error
	Sanitize(txt string, safe bool) *template.HTML
}

type AuthService interface {
	Start(ctx context.Context, isSignUp bool) string
	Callback(ctx context.Context, code, state string) (*session.Cookie, error)
}
