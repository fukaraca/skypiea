package handlers

import "github.com/fukaraca/skypiea/internal/storage"

const (
	HX_REDIRECT = "HX-REDIRECT"
)

type View struct {
	Repo *storage.Repositories
}

type Common struct {
	Repo *storage.Repositories
}

type Strict struct {
	Repo *storage.Repositories
}
