package migration

import migrate "github.com/rubenv/sql-migrate"

var Mig0004AddOauth2FieldsForUsers = &migrate.Migration{
	Id:   "Migration0004AddOauth2FieldsForUsers",
	Up:   migration0004AddOauth2FieldsForUsersUp,
	Down: migration0004AddOauth2FieldsForUsersDown,
}

var migration0004AddOauth2FieldsForUsersUp = []string{
	`ALTER TABLE users
	ADD COLUMN IF NOT EXISTS "picture" text,
	ADD COLUMN IF NOT EXISTS "auth_type" VARCHAR(32) NOT NULL DEFAULT 'password';`,
}

var migration0004AddOauth2FieldsForUsersDown = []string{
	`ALTER TABLE users
    DROP COLUMN IF EXISTS "auth_type",
    DROP COLUMN IF EXISTS "picture";`,
}
