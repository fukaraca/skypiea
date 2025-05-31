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
	"github.com/google/uuid"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func getDBTestContainer(ctx context.Context, cfg storage.Database) (*postgres.PostgresContainer, error) {
	return postgres.Run(
		ctx,
		"postgres:16.8-alpine",
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver(string(cfg.Dialect)),
		postgres.WithDatabase(cfg.PostgreSQL.Database),
		postgres.WithUsername(cfg.PostgreSQL.Username),
		postgres.WithPassword(cfg.PostgreSQL.Password),
	)
}

func prepareDB(ctx context.Context, cfg *storage.Database,
	ctr *postgres.PostgresContainer, t *testing.T) (*storage.Registry, error) {
	t.Helper()
	connStr, err := ctr.ConnectionString(ctx)
	require.NoError(t, err)
	uRL, err := url.Parse(connStr)
	require.NoError(t, err)
	cfg.PostgreSQL.Host, cfg.PostgreSQL.Port = uRL.Hostname(), uRL.Port()
	t.Log(connStr, cfg.PostgreSQL)
	t.Run("migration", func(t *testing.T) {
		err = migration.RunMigration(cfg, migrate.Up)
		require.NoError(t, err)
	})
	db, err := cfg.GetDBConn()
	require.NoError(t, err)
	out, err := os.ReadFile("seeding.sql")
	require.NoError(t, err)
	_, err = db.Exec(ctx, string(out))
	require.NoError(t, err)
	return storage.NewRegistry(cfg.Dialect, db), nil
}

func TestStorage(t *testing.T) {
	ctx := context.Background()
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

	ctr, err := getDBTestContainer(ctx, cfg)
	defer testcontainers.CleanupContainer(t, ctr)
	require.NoError(t, err)
	repo, err := prepareDB(ctx, &cfg, ctr, t)
	defer repo.Close()
	require.NoError(t, err)
	require.NotNil(t, repo)

	t.Run("Test User Operations", func(t *testing.T) {
		var u *storage.User
		u, err = repo.Users.GetUserByEmail(ctx, "johndoe@example.com")
		require.NoError(t, err)
		require.NotNil(t, u)
		u, err = repo.Users.GetUserByUUID(ctx, uuid.MustParse(u.UserUUID))
		require.NoError(t, err)
		require.NotNil(t, u)
		_, err = repo.Users.GetHPassword(ctx, u.Email)
		require.NoError(t, err)
		err = repo.DoInTx(ctx, slog.Default(), func(reg *storage.Registry) error {
			return reg.Users.DeleteUsersByUUID(ctx, uuid.MustParse(u.UserUUID))
		})
		require.NoError(t, err)
	})

	t.Run("Test Conversation Operations", func(t *testing.T) {
		var u *storage.User
		u, err = repo.Users.GetUserByEmail(ctx, "janedoe@example.com")
		require.NoError(t, err)
		require.NotNil(t, u)
		convs, err := repo.Conversations.GetConversationsByUserUUID(ctx, uuid.MustParse(u.UserUUID))
		require.NoError(t, err)
		require.NotNil(t, convs)
		msgs, err := repo.Conversations.GetConversationByID(ctx, convs[0].ID)
		require.NoError(t, err)
		require.NotNil(t, msgs)
	})

}
