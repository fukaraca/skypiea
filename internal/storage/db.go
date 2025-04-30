package storage

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Dialect string

const (
	DialectPgx      = "pgx"
	DialectPostgres = "postgres"
	DialectMySQL    = "mysql"
)

type Repositories struct {
	Users         UsersRepo
	Conversations ConversationsRepo
}

func NewRepositories(db *DB) *Repositories {
	return &Repositories{
		Users:         NewUsersRepo(db),
		Conversations: NewConversationsRepo(db),
	}
}

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
	*rw
	Dialect Dialect
	conf    *Database
}

type rw struct {
	*ro
}

func newRW(readOnly *ro) *rw {
	return &rw{readOnly}
}

type ro struct {
	*pgxpool.Pool
}

func newRO(pool *pgxpool.Pool) *ro {
	return &ro{pool}
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
		return &DB{newRW(newRO(conn)), d.Dialect, d}, nil
	}
	return nil, errors.New("no dialect was selected to get db connection")
}
