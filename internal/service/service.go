package service

import "github.com/fukaraca/skypiea/internal/storage"

// Service is main implementor of service layer. Bad naming tho
type Service struct {
	Repositories *storage.Registry
}

func New(reg *storage.Registry) *Service {
	return &Service{
		Repositories: reg,
	}
}
