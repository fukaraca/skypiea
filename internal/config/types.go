package config

import (
	"time"

	"github.com/fukaraca/skypiea/internal/api/gemini"
	"github.com/fukaraca/skypiea/internal/storage"
	"github.com/fukaraca/skypiea/pkg/gwt"
	logg "github.com/fukaraca/skypiea/pkg/log"
	"golang.org/x/oauth2"
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
	Gemini      *gemini.Config
	Oauth2      *Oauth2Config
}

type Server struct {
	Address               string        `mapstructure:"address"`
	Port                  string        `mapstructure:"port"`
	MaxBodySizeMB         int           `mapstructure:"maxBodySizeMB"`
	GinMode               string        `mapstructure:"ginMode"`
	SessionTimeout        time.Duration `mapstructure:"sessionTimeout"`
	DefaultRequestTimeout time.Duration `mapstructure:"defaultRequestTimeout"`
	Version               string
}

type Worker struct {
	IntervalTicker time.Duration `mapstructure:"intervalTicker"`
	Version        string
}

type Oauth2Config struct {
	Google oauth2.Config
}
