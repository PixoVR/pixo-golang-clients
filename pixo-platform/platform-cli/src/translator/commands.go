package translator

import (
	"fmt"
	"github.com/ollama/ollama/api"
)

var (
	baseTranslationMessages = []api.Message{
		{
			Role:    "system",
			Content: "Your goal is to take the user's input, translate it to a given language and return it to them",
		},
		{
			Role:    "system",
			Content: "The output should be in the EXACT same format as the input and contain all the same information",
		},
		{
			Role:    "system",
			Content: "The output will be used as drop-in replacement for that language, so it's very important that the translation is accurate and as close as possible to the original",
		},
	}
)

func NewTextTranslationMessages(req Request) []api.Message {
	translationCommand := fmt.Sprintf("Translate the following text from %s to %s", req.OriginalLanguage, req.TranslatedLanguage)
	translationMessage := api.Message{
		Role:    "system",
		Content: translationCommand,
	}
	inputMessage := api.Message{
		Role:    "user",
		Content: req.Text,
	}
	return append(baseTranslationMessages, translationMessage, inputMessage)
}
