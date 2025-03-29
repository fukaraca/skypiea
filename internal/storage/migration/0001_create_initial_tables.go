package migration

import "github.com/rubenv/sql-migrate"

var Mig0001CreateInitialTables = &migrate.Migration{
	Id:   "Migration0001CreateInitialTables",
	Up:   migration0001CreateInitialTablesUp,
	Down: migration0001CreateInitialTablesDown,
}

var migration0001CreateInitialTablesUp = []string{
	`CREATE TABLE roles (
    id TEXT PRIMARY KEY ,
    name TEXT NOT NULL
	);`,

	`CREATE TABLE users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY ,
    user_uuid uuid NOT NULL UNIQUE ,
    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE ,
    password TEXT,
    role TEXT,
    status TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    FOREIGN KEY (role) 
        REFERENCES roles (id)
        ON UPDATE restrict 
        ON DELETE restrict                    
	);`,

	`CREATE TABLE subscriptions (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INTEGER NOT NULL,
    plan TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    started_at TIMESTAMP WITHOUT TIME ZONE,
    valid_thru TIMESTAMP WITHOUT TIME ZONE,
    FOREIGN KEY (user_id) 
        REFERENCES users (id)
        ON UPDATE CASCADE 
        ON DELETE CASCADE    
	);`,
}

var migration0001CreateInitialTablesDown = []string{
	`DROP TABLE IF EXISTS subscriptions;`,

	`DROP TABLE IF EXISTS users;`,

	`DROP TABLE IF EXISTS roles;`,
}
