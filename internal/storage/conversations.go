package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	getConversationsByUserIDPG  = `SELECT * FROM conversations WHERE user_uuid = $1 ORDER BY updated_at DESC LIMIT $2;`
	getConversationByIDPG       = `SELECT * FROM messages WHERE conv_id = $1 ORDER BY created_at ASC;`
	getMessageByIDPG            = `SELECT * FROM messages WHERE id = $1;`
	getResponseByIdPG           = `SELECT * FROM messages WHERE response_to = $1;`
	addNewConversationPG        = `INSERT INTO conversations(user_uuid,title,metadata) VALUES ($1,$2,$3) RETURNING id;`
	addNewMessagePG             = `INSERT INTO messages(conv_id,model_id,by_user,message,metadata,response_to) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id;`
	updateMessagePG             = `UPDATE messages SET message = $1 WHERE id= $2;`
	deleteConversationByIDPG    = `DELETE FROM conversations WHERE id = $1;`
	bumpConversationUpdatedAtPG = `UPDATE conversations SET updated_at = CURRENT_TIMESTAMP WHERE id = $1;`
	verifyUserForConversation   = `SELECT EXISTS (SELECT 1 FROM conversations WHERE id = $1 AND user_uuid = $2);`
	verifyUserForMessage        = `SELECT EXISTS (SELECT 1 FROM messages m JOIN conversations c ON c.id = m.conv_id WHERE m.id = $1 AND c.user_uuid = $2);`
)

type Conversation struct {
	ID        int
	UserUUID  string
	Title     string
	Metadata  *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Message struct {
	ID          int
	ConvID      int
	ModelID     string
	ByUser      bool
	MessageText *string
	Metadata    *string
	CreatedAt   time.Time
	ResponseTo  *int
}

// TODO: this is conventional and not needed right now. It is better created and implemented on caller side
type ConversationsRepo interface {
	GetConversationsByUserUUID(context.Context, uuid.UUID) ([]*Conversation, error)
	GetConversationByID(ctx context.Context, conversationID int) ([]*Message, error)
	GetMessageByID(ctx context.Context, i int) (*Message, error)
	// GetResponseByID returns message of response by the ID of query-message
	GetResponseByID(ctx context.Context, i int) (*Message, error)
	NewConversation(ctx context.Context, userID uuid.UUID, title, metadata string) (int, error) // returning conversation id if succeeds
	DeleteConversation(ctx context.Context, conversationID int) error
	AppendNewMessage(ctx context.Context, msg *Message) (int, error)
	UpdateMessage(ctx context.Context, msgID int) error
	BumpConversationUpdatedAtV(ctx context.Context, convID int) error
	VerifyUserForConversation(ctx context.Context, userUUID uuid.UUID, convID int) (bool, error)
	VerifyUserForMessage(ctx context.Context, userUUID uuid.UUID, msgID int) (bool, error)
}

type conversationsRepoPgx struct {
	dbConn
}

func NewConversationsRepo(dia Dialect, conn dbConn) *conversationsRepoPgx {
	switch dia {
	case DialectPostgres, DialectPgx:
		return &conversationsRepoPgx{conn}
	}
	return nil
}

func (c *conversationsRepoPgx) GetConversationsByUserUUID(ctx context.Context, userID uuid.UUID) ([]*Conversation, error) {
	out := make([]*Conversation, 0)
	limit := 20
	rows, err := c.Query(ctx, getConversationsByUserIDPG, userID.String(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() { // TODO consider using generic solutions to map directly to slice
		var t Conversation
		err = rows.Scan(&t.ID, &t.UserUUID, &t.Title, &t.Metadata, &t.CreatedAt, &t.UpdatedAt)
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
		err = rows.Scan(&m.ID, &m.ConvID, &m.ModelID, &m.ByUser, &m.MessageText, &m.Metadata, &m.CreatedAt, &m.ResponseTo)
		if err != nil {
			return nil, err
		}
		out = append(out, &m)
	}
	return out, nil
}

func (c *conversationsRepoPgx) GetMessageByID(ctx context.Context, i int) (*Message, error) {
	var out Message
	rows, err := c.Query(ctx, getMessageByIDPG, i)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&out.ID, &out.ConvID, &out.ModelID, &out.ByUser, &out.MessageText, &out.Metadata, &out.CreatedAt, &out.ResponseTo)
		if err != nil {
			return nil, err
		}
	}
	return &out, nil
}

// GetResponseByID returns message of response by the ID of query-message
func (c *conversationsRepoPgx) GetResponseByID(ctx context.Context, i int) (*Message, error) {
	var out Message
	err := c.QueryRow(ctx, getResponseByIdPG, i).Scan(&out.ID, &out.ConvID, &out.ModelID, &out.ByUser, &out.MessageText, &out.Metadata, &out.CreatedAt, &out.ResponseTo)
	if err != nil {
		return nil, err
	}
	return &out, nil
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
	if err := c.QueryRow(ctx, addNewMessagePG, msg.ConvID, msg.ModelID, msg.ByUser, msg.MessageText, msg.Metadata, msg.ResponseTo).Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

// TODO: just emptying for now idk
func (c *conversationsRepoPgx) UpdateMessage(ctx context.Context, msgID int) error {
	_, err := c.Exec(ctx, updateMessagePG, "", msgID)
	return err
}

func (c *conversationsRepoPgx) BumpConversationUpdatedAtV(ctx context.Context, convID int) error {
	_, err := c.Exec(ctx, bumpConversationUpdatedAtPG, convID)
	return err
}

func (c *conversationsRepoPgx) VerifyUserForConversation(ctx context.Context, userUUID uuid.UUID, convID int) (bool, error) {
	var ok bool
	err := c.QueryRow(ctx, verifyUserForConversation, convID, userUUID.String()).Scan(&ok)
	if err != nil {
		return ok, err
	}
	return ok, nil
}

func (c *conversationsRepoPgx) VerifyUserForMessage(ctx context.Context, userUUID uuid.UUID, msgID int) (bool, error) {
	var ok bool
	err := c.QueryRow(ctx, verifyUserForConversation, msgID, userUUID.String()).Scan(&ok)
	if err != nil {
		return ok, err
	}
	return ok, nil
}
