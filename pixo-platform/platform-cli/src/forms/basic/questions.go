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
			err error
		)

		if question.Prompt == "" {
			question.Prompt = forms.CleanPrompt(question.Key)
		}

		switch question.Type {
		case forms.Input:
			question.Answer = new(string)
			err = f.GetResponseFromUser(&question)
		case forms.SensitiveInput:
			question.Answer = new(string)
			err = f.GetSensitiveResponseFromUser(&question)
		case forms.Confirm:
			question.Answer = new(bool)
			err = f.Confirm(&question)
		case forms.Select:
			question.Answer = new(string)
			err = f.Select(&question)
		case forms.SelectID:
			question.Answer = new(int)
			err = f.SelectID(&question)
		case forms.MultiSelectIDs:
			question.Answer = new([]int)
			err = f.MultiSelectIDs(&question)
		case forms.MultiSelect:
			question.Answer = new([]string)
			err = f.MultiSelect(&question)
		default:
			log.Panic().Msg("unknown question type")
		}

		if !question.Optional {
			if err != nil || question.Answer == nil {
				return nil, fmt.Errorf("%s not provided", forms.CleanPrompt(question.Key))
			}
			answers[question.Key] = question.Answer
		}
	}

	return answers, nil
}
