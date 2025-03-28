package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Database struct {
	Dialect    string           `yaml:"dialect"`
	PostgreSQL PostgreSQLConfig `yaml:"postgresql"`
}

type PostgreSQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
}

type Conn struct {
	*sql.DB
	Dialect string
}

func (d *Database) GetDBConn(dialect string) (*Conn, error) {
	switch dialect {
	case "postgres":
		dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
			d.PostgreSQL.Username,
			d.PostgreSQL.Password,
			d.PostgreSQL.Host,
			d.PostgreSQL.Port,
			d.PostgreSQL.Database)
		conn, err := pgxpool.New(context.Background(), dbURL)
		if err != nil {
			return nil, err
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		if err = conn.Ping(ctx); err != nil {
			return nil, err
		}
	}
	return nil, fmt.Errorf("no dialect was selected to get db connection")
}
