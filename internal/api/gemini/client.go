package gemini

import (
	"context"

	"google.golang.org/genai"
)

type Config struct {
	*genai.ClientConfig
	SupportedModels []string
}

type Client struct {
	config *Config
	*genai.Client
}

// NewClient returns reusable gemini client
func NewClient(cfg *Config) (*Client, error) {
	client, err := genai.NewClient(context.Background(), cfg.ClientConfig)
	if err != nil {
		return nil, err
	}
	return &Client{config: cfg, Client: client}, nil
}

func (c *Client) AskToGemini(ctx context.Context, msg, model string) (string, error) {
	result, err := c.Models.GenerateContent(
		ctx,
		model,
		genai.Text(msg),
		nil,
	)
	if err != nil {
		return "", err
	}
	return result.Text(), nil
}

func (c *Client) GetAllSupportedModels() []string {
	return c.config.SupportedModels
}
