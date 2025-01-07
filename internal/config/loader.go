package config

import (
	"github.com/fukaraca/skypiea/internal/model"
	"github.com/spf13/viper"
)

func LoadConfig(filename, path string) (*model.Config, error) {
	v := viper.New()
	v.SetConfigName(filename)
	v.SetConfigType("yml")
	v.AddConfigPath(path)
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var config model.Config
	err = v.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
