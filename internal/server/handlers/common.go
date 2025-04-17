package handlers

import (
	"context"
	"time"

	"github.com/fukaraca/skypiea/internal/config"
	"github.com/fukaraca/skypiea/internal/storage"
	"github.com/gin-gonic/gin"
)

const (
	HX_REDIRECT = "HX-REDIRECT"

	AlertLevelInfo    AlertLevel = "alert-info"
	AlertLevelError   AlertLevel = "alert-danger"
	AlertLevelWarning AlertLevel = "alert-warning"

	defaultReqTimeout = time.Second * 30
)

type AlertLevel string

type Common struct {
	reqTimeout time.Duration
}

func NewCommonHandler(s *config.Config) *Common {
	timeout := defaultReqTimeout
	if s.Server.DefaultRequestTimeout > time.Second {
		timeout = s.Server.DefaultRequestTimeout
	}
	return &Common{reqTimeout: timeout}
}

type View struct {
	*Common
	Repo *storage.Repositories
}

func NewViewHandler(c *Common, repo *storage.Repositories) *View {
	return &View{Common: c, Repo: repo}
}

type Open struct {
	*Common
	Repo *storage.Repositories
}

func NewOpenHandler(c *Common, repo *storage.Repositories) *Open {
	return &Open{Common: c, Repo: repo}
}

type Strict struct {
	*Common
	Repo *storage.Repositories
}

func NewStrictHandler(c *Common, repo *storage.Repositories) *Strict {
	return &Strict{Common: c, Repo: repo}
}

func (h *Common) CtxWithTimout(c *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c, h.reqTimeout)
}
