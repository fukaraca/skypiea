package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Conversation struct {
	ID        int
	UserUUID  string
	Title     string
	Metadata  string
	CreatedAt time.Time
}

type Message struct {
	ID          int
	ConvID      int
	ModelID     string
	ByUser      bool
	MessageText string
	Metadata    string
	CreatedAt   time.Time
}

type ConversationsRepo interface {
	GetConversationsByUserUUID(context.Context, uuid.UUID) ([]*Conversation, error)
	GetConversationByID(context.Context, int) ([]*Message, error)
	NewConversation(ctx context.Context, userID uuid.UUID, title, metadata string) (int, error) // returning conversation id if succeeds
	AppendNewMessage(ctx context.Context, msg *Message) (int, error)
	UpdateMessage(ctx context.Context, msgID int) error
}

type conversationsRepoPgx struct {
	*pgxpool.Pool
}

func NewConversationsRepo(db *DB) ConversationsRepo {
	switch db.Dialect {
	case DialectPostgres, DialectPgx:
		return &conversationsRepoPgx{db.Pool}
	}
	return nil
}

func (c *conversationsRepoPgx) GetConversationsByUserUUID(ctx context.Context, userID uuid.UUID) ([]*Conversation, error) {
	//TODO implement me
	panic("implement me")
}

func (c *conversationsRepoPgx) GetConversationByID(ctx context.Context, i int) ([]*Message, error) {
	//TODO implement me
	panic("implement me")
}

func (c *conversationsRepoPgx) NewConversation(ctx context.Context, userID uuid.UUID, title, metadata string) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (c *conversationsRepoPgx) AppendNewMessage(ctx context.Context, msg *Message) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (c *conversationsRepoPgx) UpdateMessage(ctx context.Context, msgID int) error {
	//TODO implement me
	panic("implement me")
}
