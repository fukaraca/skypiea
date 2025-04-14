package handlers

import "github.com/fukaraca/skypiea/internal/storage"

const (
	HX_REDIRECT = "HX-REDIRECT"

	AlertLevelInfo    AlertLevel = "alert-info"
	AlertLevelError   AlertLevel = "alert-danger"
	AlertLevelWarning AlertLevel = "alert-warning"
)

type AlertLevel string

type Common struct{}

type View struct {
	Common
	Repo *storage.Repositories
}

type Open struct {
	Common
	Repo *storage.Repositories
}

type Strict struct {
	Common
	Repo *storage.Repositories
}
