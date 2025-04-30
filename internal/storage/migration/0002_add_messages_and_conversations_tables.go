package migration

import "github.com/rubenv/sql-migrate"

var Mig0002AddMessagesAndConversationsTables = &migrate.Migration{
	Id:   "Migration0002AddMessagesAndConversationsTables",
	Up:   migration0002AddMessagesAndConversationsTablesUp,
	Down: migration0002AddMessagesAndConversationsTablesDown,
}

var migration0002AddMessagesAndConversationsTablesUp = []string{
	`CREATE TABLE conversations (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_uuid uuid NOT NULL, 
    title TEXT NOT NULL,
    metadata TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
        FOREIGN KEY (user_uuid) 
        REFERENCES users (user_uuid)
        ON UPDATE CASCADE 
        ON DELETE CASCADE                 
	);`,

	`CREATE TABLE messages (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY ,
    conv_id BIGINT NOT NULL,
    model_id TEXT,
    by_user BOOLEAN NOT NULL DEFAULT FALSE,
    message TEXT,
    metadata TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    FOREIGN KEY (conv_id) 
        REFERENCES conversations (id)
        ON UPDATE cascade 
        ON DELETE cascade                    
	);`,
}

var migration0002AddMessagesAndConversationsTablesDown = []string{
	`DROP TABLE IF EXISTS messages;`,

	`DROP TABLE IF EXISTS conversations;`,
}
