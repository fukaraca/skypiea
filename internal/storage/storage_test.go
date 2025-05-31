package storage_test

import (
	"context"
	"log/slog"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/fukaraca/skypiea/internal/storage"
	"github.com/fukaraca/skypiea/internal/storage/migration"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func prepareDB(ctx context.Context, t *testing.T) (*storage.Registry, error) {
	cfg := storage.Database{
		Dialect: storage.DialectPostgres,
		PostgreSQL: storage.PostgreSQLConfig{
			Host:     "localhost",
			Port:     "5432",
			Username: "test_pg_user",
			Password: "test_pg_pass",
			Database: "test_skypiea_ai",
			Timeout:  100 * time.Second,
		},
	}

	ctr, err := postgres.Run(
		ctx,
		"postgres:16.8-alpine",
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver(string(cfg.Dialect)),
		postgres.WithDatabase(cfg.PostgreSQL.Database),
		postgres.WithUsername(cfg.PostgreSQL.Username),
		postgres.WithPassword(cfg.PostgreSQL.Password),
		//postgres.WithInitScripts(`CREATE DATABASE test_skypiea_ai;`),
	)
	defer testcontainers.CleanupContainer(t, ctr)
	require.NoError(t, err)

	connStr, err := ctr.ConnectionString(ctx)
	require.NoError(t, err)
	uRL, err := url.Parse(connStr)
	require.NoError(t, err)
	cfg.PostgreSQL.Host, cfg.PostgreSQL.Port = uRL.Hostname(), uRL.Port()
	t.Log(connStr, cfg.PostgreSQL)

	err = migration.RunMigration(&cfg, migrate.Up)
	require.NoError(t, err)
	db, err := cfg.GetDBConn()
	require.NoError(t, err)
	out, err := os.ReadFile("seeding.sql")
	require.NoError(t, err)
	_, err = db.Exec(ctx, string(out))
	require.NoError(t, err)
	return storage.NewRegistry(cfg.Dialect, db), nil
}

func TestConversationsRepoPgx(t *testing.T) {
	ctx := context.Background()
	repo, err := prepareDB(ctx, t)
	require.NoError(t, err)
	require.NotNil(t, repo)
	var u *storage.User
	err = repo.DoInTx(ctx, slog.Default(), func(reg *storage.Registry) error {
		u, err = reg.Users.GetUserByEmail(ctx, "johndoe@example.com")
		return err
	})
	require.NoError(t, err)
	require.Equal(t, "johndoe@example.com", u.Email)

}
