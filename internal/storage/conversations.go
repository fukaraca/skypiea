package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	getConversationsByUserIDPG = `SELECT * FROM conversations WHERE user_uuid = $1 ORDER BY created_at DESC;`
	getConversationByIDPG      = `SELECT * FROM messages WHERE conv_id = $1 ORDER BY created_at DESC;`
	addNewConversationPG       = `INSERT INTO conversations(user_uuid,title,metadata) VALUES ($1,$2,$3) RETURNING id;`
	addNewMessagePG            = `INSERT INTO messages(conv_id,model_id,by_user,message,metadata) VALUES ($1,$2,$3,$4,$5) RETURNING id;`
	updateMessagePG            = `UPDATE messages SET message = '$1' WHERE id= $2;`
	deleteConversationByIDPG   = `DELETE FROM conversations WHERE id = $1;`
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
	GetConversationByID(ctx context.Context, conversationID int) ([]*Message, error)
	NewConversation(ctx context.Context, userID uuid.UUID, title, metadata string) (int, error) // returning conversation id if succeeds
	DeleteConversation(ctx context.Context, conversationID int) error
	AppendNewMessage(ctx context.Context, msg *Message) (int, error)
	UpdateMessage(ctx context.Context, msgID int) error
}

type conversationsRepoPgx struct {
	dbConn
}

func NewConversationsRepo(dia Dialect, conn dbConn) ConversationsRepo {
	switch dia {
	case DialectPostgres, DialectPgx:
		return &conversationsRepoPgx{conn}
	}
	return nil
}

func (c *conversationsRepoPgx) GetConversationsByUserUUID(ctx context.Context, userID uuid.UUID) ([]*Conversation, error) {
	out := make([]*Conversation, 0)
	rows, err := c.Query(ctx, getConversationsByUserIDPG, userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() { // TODO consider using generic solutions to map directly to slice
		var t Conversation
		err = rows.Scan(&t.ID, &t.UserUUID, t.Title, t.Metadata, t.CreatedAt)
		if err != nil {
			return nil, err
		}
		out = append(out, &t)
	}
	return out, nil
}

func (c *conversationsRepoPgx) GetConversationByID(ctx context.Context, i int) ([]*Message, error) {
	out := make([]*Message, 0)
	rows, err := c.Query(ctx, getConversationByIDPG, i)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var m Message
		err = rows.Scan(&m.ID, &m.ConvID, &m.ModelID, &m.ByUser, &m.MessageText, &m.Metadata, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		out = append(out, &m)
	}
	return out, nil
}

func (c *conversationsRepoPgx) NewConversation(ctx context.Context, userID uuid.UUID, title, metadata string) (int, error) {
	var id int
	if err := c.QueryRow(ctx, addNewConversationPG, userID.String(), title, metadata).Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

func (c *conversationsRepoPgx) DeleteConversation(ctx context.Context, conversationID int) error {
	_, err := c.Exec(ctx, deleteConversationByIDPG, conversationID)
	return err
}

func (c *conversationsRepoPgx) AppendNewMessage(ctx context.Context, msg *Message) (int, error) {
	var id int
	if err := c.QueryRow(ctx, addNewMessagePG, msg.ConvID, msg.ModelID, msg.ByUser, msg.MessageText, msg.Metadata).Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

// TODO: just emptying for now idk
func (c *conversationsRepoPgx) UpdateMessage(ctx context.Context, msgID int) error {
	_, err := c.Exec(ctx, updateMessagePG, "", msgID)
	return err
}
