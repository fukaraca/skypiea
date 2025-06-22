package migration

import "github.com/rubenv/sql-migrate"

var Mig0003AddAboutMeColumnsUsersTable = &migrate.Migration{
	Id:   "Migration0003AddAboutMeColumnsUsersTable",
	Up:   migration0003AddAboutMeColumnsUsersTableUp,
	Down: migration0003AddAboutMeColumnsUsersTableDown,
}

var migration0003AddAboutMeColumnsUsersTableUp = []string{
	`ALTER TABLE users
	ADD COLUMN IF NOT EXISTS "phone_number" VARCHAR(36),
	ADD COLUMN IF NOT EXISTS "about_me" text,
	ADD COLUMN IF NOT EXISTS "summary" text;`,
}

var migration0003AddAboutMeColumnsUsersTableDown = []string{
	`ALTER TABLE users
    DROP COLUMN IF EXISTS "summary",
    DROP COLUMN IF EXISTS "about_me",
    DROP COLUMN IF EXISTS "phone_number";`,
}
