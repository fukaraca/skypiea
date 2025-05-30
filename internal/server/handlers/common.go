package handlers

import (
	"time"

	"github.com/fukaraca/skypiea/internal/config"
)

const (
	HX_REDIRECT = "HX-REDIRECT"

	ALInfo    AlertLevel = "alert-info"
	ALError   AlertLevel = "alert-danger"
	ALWarning AlertLevel = "alert-warning"

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
	MessageSvc MessageService
	UserSvc    UserService
}

func NewViewHandler(c *Common, msgSvc MessageService, userSvc UserService) *View {
	return &View{Common: c, MessageSvc: msgSvc, UserSvc: userSvc}
}

type Open struct {
	*Common
	UserSvc UserService
}

func NewOpenHandler(c *Common, userSvc UserService) *Open {
	return &Open{Common: c, UserSvc: userSvc}
}

type Strict struct {
	*Common
	MessageSvc MessageService
	UserSvc    UserService
}

func NewStrictHandler(c *Common, msgSvc MessageService, userSvc UserService) *Strict {
	return &Strict{Common: c, MessageSvc: msgSvc, UserSvc: userSvc}
}
