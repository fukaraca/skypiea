package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func NewConfig() *Config {
	return &Config{}
}

var envs = []string{
	"database.dialect",
	"database.postgresql.host",
	"database.postgresql.username",
	"database.postgresql.password",
	"database.postgresql.database",
	"database.dynamodb.region",
	"database.dynamodb.endpoint",
	"database.dynamodb.accessKeyId",
	"database.dynamodb.secretAccessKey",
}

func (c *Config) Load(filename, path string) error {
	v := viper.New()
	v.SetConfigName(filename)
	v.SetConfigType("yml")
	v.SetDefault("database.postgresql.port", 5432)
	v.SetDefault("database.postgresql.sslmode", "disable")
	v.AddConfigPath(path)
	v.AllowEmptyEnv(true)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	for _, env := range envs {
		if err := v.BindEnv(env); err != nil {
			return err
		}
	}
	err := v.ReadInConfig()
	if err != nil {
		return err
	}
	fmt.Println(os.Environ())

	err = v.Unmarshal(c)
	if err != nil {
		return err
	}

	return nil
}
