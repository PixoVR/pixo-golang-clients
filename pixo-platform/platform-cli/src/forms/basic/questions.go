package basic

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/rs/zerolog/log"
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
			prompt = forms.CleanPrompt(question.Key)
		}

		switch question.Type {
		case forms.Input:
			answer = new(string)
			err = f.GetResponseFromUser(prompt, answer.(*string))
		case forms.SensitiveInput:
			answer = new(string)
			err = f.GetSensitiveResponseFromUser(prompt, answer.(*string))
		case forms.Confirm:
			answer = new(bool)
			err = f.Confirm(prompt, answer.(*bool))
		case forms.MultiSelectIDs:
			answer = new([]int)
			err = f.MultiSelectIDs(prompt, question.Options, answer.(*[]int))
		case forms.MultiSelect:
			answer = new([]string)
			err = f.MultiSelect(prompt, question.Options, answer.(*[]string))
		default:
			log.Panic().Msg("unknown question type")
		}

		if err != nil {
			return nil, fmt.Errorf("%s not provided", forms.CleanPrompt(question.Key))
		}
		answers[question.Key] = answer
		answer = nil
	}
	return answers, nil
}
