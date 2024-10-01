package translator

import (
	"context"
	"github.com/ollama/ollama/api"
	"os"
)

type Client interface {
	Chat(ctx context.Context, req *api.ChatRequest, responseFunc api.ChatResponseFunc) error
}

type Translator struct {
	model  string
	client Client
}

// NewTranslator creates a new Translator instance given a chat client
func NewTranslator(input ...Client) (*Translator, error) {
	var client Client
	if len(input) > 0 {
		client = input[0]
	} else {
		chatClient, err := api.ClientFromEnvironment()
		if err != nil {
			return nil, err
		}
		client = chatClient
	}

	model, ok := os.LookupEnv("LLM_MODEL")
	if !ok {
		model = "icky/translate"
	}

	return &Translator{
		model:  model,
		client: client,
	}, nil
}

func (t *Translator) Model() string {
	return t.model
}
