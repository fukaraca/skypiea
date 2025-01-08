package config

import logg "github.com/fukaraca/skypiea/pkg/log"

const (
	ModeHttpServer       = "server"
	ModeBackgroundWorker = "worker"
)

type Config struct {
	RunningMode string
	Server      Server
	Log         logg.Config
}

type Server struct {
	Address       string `yaml:"server.address"`
	Port          string `yaml:"server.port"`
	MaxBodySizeMB int    `yaml:"server.maxBodySizeMB"`
}
