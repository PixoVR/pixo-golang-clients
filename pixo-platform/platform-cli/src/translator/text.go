package translator

import (
	"context"
	"github.com/ollama/ollama/api"
	"strings"
)

// TranslateText will translate text from one language to another
func (t *Translator) TranslateText(ctx context.Context, req Request, respFunc func(string) error) error {
	chatReq := api.ChatRequest{
		Model:    t.model,
		Messages: NewTextTranslationMessages(req),
	}

	chatRespFunc := func(resp api.ChatResponse) error {
		if strings.Trim(resp.Message.Content, " ") == "" {
			return nil
		}
		return respFunc(resp.Message.Content)
	}

	return t.client.Chat(ctx, &chatReq, chatRespFunc)
}
