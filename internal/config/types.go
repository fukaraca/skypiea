package config

import (
	"time"

	"github.com/fukaraca/skypiea/internal/storage"
	"github.com/fukaraca/skypiea/pkg/gwt"
	logg "github.com/fukaraca/skypiea/pkg/log"
)

const (
	ModeHttpServer       = "server"
	ModeBackgroundWorker = "worker"
)

type Config struct {
	ServiceMode string
	Server      *Server
	Worker      *Worker
	Log         logg.Config
	Database    *storage.Database
	JWT         *gwt.Config
}

type Server struct {
	Address               string        `yaml:"server.address"`
	Port                  string        `yaml:"server.port"`
	MaxBodySizeMB         int           `yaml:"server.maxBodySizeMB"`
	GinMode               string        `yaml:"server.ginMode"`
	SessionTimeout        time.Duration `yaml:"sessionTimeout"`
	DefaultRequestTimeout time.Duration `yaml:"defaultRequestTimeout"`
	Version               string
}

type Worker struct {
	IntervalTicker time.Duration `yaml:"intervalTicker"`
	Version        string
}
