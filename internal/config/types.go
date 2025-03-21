package config

import (
	"github.com/fukaraca/skypiea/pkg/jwt"
	logg "github.com/fukaraca/skypiea/pkg/log"
)

const (
	ModeHttpServer       = "server"
	ModeBackgroundWorker = "worker"
)

type Config struct {
	ServiceMode string
	Server      *Server
	Log         logg.Config
	JWT         *jwt.Config
}

type Server struct {
	Address       string `yaml:"server.address"`
	Port          string `yaml:"server.port"`
	MaxBodySizeMB int    `yaml:"server.maxBodySizeMB"`
	GinMode       string `yaml:"server.ginMode"`
}
