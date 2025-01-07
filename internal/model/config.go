package model

import (
	logg "github.com/fukaraca/skypiea/pkg/log"
)

const (
	MODEHTTPSERVER       = "http"
	MODEBACKGROUNDWORKER = "background"
)

type Config struct {
	Server ServerConfig
	Log    logg.Config
}

type ServerConfig struct {
	Address       string `yaml:"server.address"`
	Port          string `yaml:"server.port"`
	MaxBodySizeMB int    `yaml:"server.maxBodySizeMB"`
}
