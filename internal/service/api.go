package service

import "context"

// TODO just experimenting
type GeminiAPI interface {
	AskToGemini(ctx context.Context, msg, model string) (string, error)
	GetAllSupportedModels() []string
}
