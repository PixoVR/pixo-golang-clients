package basic

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/rs/zerolog/log"
	"strings"
)

func (f *Handler) AskQuestions(questions []forms.Question) (map[string]interface{}, error) {
	answers := make(map[string]interface{})
	for _, question := range questions {
		var (
			prompt string
			answer interface{}
			err    error
		)

		if question.Prompt != "" {
			prompt = question.Prompt
		} else {
			prompt = strings.ToUpper(question.Key)
		}

		switch question.Type {
		case forms.Input:
			answer, err = f.GetResponseFromUser(prompt)
		case forms.SensitiveInput:
			answer, err = f.GetSensitiveResponseFromUser(prompt)
		case forms.Confirm:
			answer = f.Confirm(prompt)
		case forms.MultiSelectIDs:
			answer, err = f.MultiSelectIDs(prompt, question.Options)
		case forms.MultiSelect:
			answer, err = f.MultiSelect(prompt, question.Options)
		default:
			log.Panic().Msg("unknown question type")
		}

		if err != nil {
			return nil, fmt.Errorf("%s not provided", prompt)
		}
		answers[question.Key] = answer
	}
	return answers, nil
}
