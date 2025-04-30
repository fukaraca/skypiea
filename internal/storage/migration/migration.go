package migration

import (
	"database/sql"
	"fmt"

	"github.com/fukaraca/skypiea/internal/storage"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rubenv/sql-migrate"
)

func RunMigration(conf *storage.Database, direction migrate.MigrationDirection) error {
	conn, err := getConn(conf)
	if err != nil {
		return err
	}
	source := migrate.MemoryMigrationSource{
		Migrations: bindMigrations(),
	}

	_, err = migrate.Exec(conn, string(conf.Dialect), source, direction)
	return err
}

func bindMigrations() []*migrate.Migration {
	return []*migrate.Migration{
		Mig0001CreateInitialTables,
		Mig0002AddMessagesAndConversationsTables,
	}
}

func getConn(conf *storage.Database) (*sql.DB, error) {
	switch conf.Dialect {
	case storage.DialectPgx, storage.DialectPostgres:
		conn, err := conf.GetDBConn()
		if err != nil {
			return nil, err
		}
		sqlDB := stdlib.OpenDBFromPool(conn.Pool)
		return sqlDB, nil
	}
	return nil, fmt.Errorf("unsupported dialect: %q", conf.Dialect)
}
