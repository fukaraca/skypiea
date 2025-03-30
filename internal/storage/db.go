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

type Database struct {
	Dialect    Dialect          `yaml:"dialect"`
	PostgreSQL PostgreSQLConfig `yaml:"postgresql"`
}

type PostgreSQLConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
}

type Conn struct {
	*pgxpool.Pool
	Dialect Dialect
}

func (d *Database) GetDBConn() (*Conn, error) {
	switch d.Dialect {
	case DialectPostgres, DialectPgx:
		dbURL := fmt.Sprintf("postgres://%s:%s@%s/%s",
			d.PostgreSQL.Username,
			d.PostgreSQL.Password,
			net.JoinHostPort(d.PostgreSQL.Host, d.PostgreSQL.Port),
			d.PostgreSQL.Database)
		conn, err := pgxpool.New(context.Background(), dbURL)
		if err != nil {
			return nil, err
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) //nolint: govet
		if err = conn.Ping(ctx); err != nil {
			return nil, err
		}
		return &Conn{conn, d.Dialect}, nil
	}
	return nil, errors.New("no dialect was selected to get db connection")
}
