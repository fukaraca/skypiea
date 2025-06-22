package service

import (
	"github.com/fukaraca/skypiea/internal/storage"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

// Service is main implementor of service layer. Bad naming tho
type Service struct {
	Repositories     *storage.Registry
	GeminiClient     GeminiAPI
	policy           *bluemonday.Policy
	MDSafe, MDUnsafe goldmark.Markdown
}

func New(reg *storage.Registry, geminiClient GeminiAPI) *Service {
	return &Service{
		Repositories: reg,
		GeminiClient: geminiClient,
		policy:       bluemonday.UGCPolicy(),
		MDUnsafe: goldmark.New(
			goldmark.WithRendererOptions(
				html.WithUnsafe(), // <-- Allow raw HTML but this is for trusted sources like gemini
			),
		),
		MDSafe: goldmark.New(),
	}
}
