package storage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dialect string

const (
	defaultTxTimeout = 10 * time.Second

	DialectPgx      = "pgx"
	DialectPostgres = "postgres"
	DialectMySQL    = "mysql"
)

type Database struct {
	Dialect    Dialect          `yaml:"dialect"`
	PostgreSQL PostgreSQLConfig `yaml:"postgresql"`
}

type PostgreSQLConfig struct {
	Host     string        `yaml:"host"`
	Port     string        `yaml:"port"`
	Username string        `yaml:"username"`
	Password string        `yaml:"password"`
	Database string        `yaml:"database"`
	SSLMode  string        `yaml:"sslmode"`
	Timeout  time.Duration `yaml:"timeout"`
}

type DB struct {
	*pgxpool.Pool
	Dialect Dialect
	conf    *Database
}

func (d *Database) GetDBConn() (*DB, error) {
	switch d.Dialect {
	case DialectPostgres, DialectPgx:
		dbURL := fmt.Sprintf("postgres://%s:%s@%s/%s",
			d.PostgreSQL.Username,
			d.PostgreSQL.Password,
			net.JoinHostPort(d.PostgreSQL.Host, d.PostgreSQL.Port),
			d.PostgreSQL.Database)
		conf, err := pgxpool.ParseConfig(dbURL)
		if err != nil {
			return nil, err
		}
		conn, err := pgxpool.NewWithConfig(context.Background(), conf)
		if err != nil {
			return nil, err
		}
		ctx, _ := context.WithTimeout(context.Background(), d.PostgreSQL.Timeout) //nolint: govet
		if err = conn.Ping(ctx); err != nil {
			return nil, err
		}
		return &DB{conn, d.Dialect, d}, nil
	}
	return nil, errors.New("no dialect was selected to get db connection")
}

// dbConn supports TX, pool, conn etc..
type dbConn interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row
}

type Registry struct {
	dia           Dialect
	conn          dbConn
	Users         UsersRepo
	Conversations ConversationsRepo
}

// NewRegistry returns new registry, this can be regular pool of sessions or a transaction registry
func NewRegistry(dia Dialect, conn dbConn) *Registry {
	return &Registry{
		dia:           dia,
		conn:          conn,
		Users:         NewUsersRepo(dia, conn),
		Conversations: NewConversationsRepo(dia, conn),
	}
}

func (r *Registry) DoInTx(ctx context.Context, logger *slog.Logger, fn func(reg *Registry) error) error {
	if logger == nil {
		logger = slog.Default()
	}
	ctx, cancel := context.WithTimeout(ctx, defaultTxTimeout)
	defer cancel()

	tx, err := r.conn.Begin(ctx) // defaults read committed
	if err != nil {
		logger.Error("DoInTx.Begin failed", err)
		return errors.Join(err, errors.New("DoInTx.Begin failed"))
	}
	reg := NewRegistry(r.dia, tx)
	if err = fn(reg); err != nil {
		logger.Error("DoInTx.fn() failed", err)
		return errors.Join(err, tx.Rollback(ctx))
	}
	return tx.Commit(ctx)
}
