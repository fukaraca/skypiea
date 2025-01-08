package config

import (
	"github.com/spf13/viper"
)

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load(filename, path string) error {
	v := viper.New()
	v.SetConfigName(filename)
	v.SetConfigType("yml")
	v.AddConfigPath(path)
	err := v.ReadInConfig()
	if err != nil {
		return err
	}
	err = v.Unmarshal(c)
	if err != nil {
		return err
	}

	return nil
}
