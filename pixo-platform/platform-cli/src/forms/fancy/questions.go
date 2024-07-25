package fancy

import (
	"errors"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/charmbracelet/huh"
	"strings"
)

func (f *Handler) AskQuestions(questions []forms.Question) (map[string]interface{}, error) {
	answers := make(map[string]interface{})
	var groupItems []huh.Field

	for _, question := range questions {
		var prompt string
		if question.Prompt != "" {
			prompt = question.Prompt
		} else {
			prompt = strings.ToUpper(question.Key)
		}

		switch question.Type {
		case forms.Input:
			var response string
			answers[question.Key] = &response
			groupItems = append(groupItems, f.InputField(prompt, &response))
		case forms.SensitiveInput:
			var response string
			answers[question.Key] = &response
			groupItems = append(groupItems, f.SensitiveInputField(prompt, &response))
		case forms.Confirm:
			var response bool
			answers[question.Key] = &response
			groupItems = append(groupItems, f.ConfirmField(prompt, &response))
		case forms.MultiSelectIDs:
			var response []int
			answers[question.Key] = &response
			groupItems = append(groupItems, f.MultiSelectIDsInput(prompt, question.Options, response))
		case forms.MultiSelect:
			var response []string
			answers[question.Key] = &response
			groupItems = append(groupItems, f.MultiSelectInput(prompt, question.Options, response))
		default:
			return nil, errors.New("unknown question type")
		}
	}

	if err := huh.NewForm(huh.NewGroup(groupItems...)).Run(); err != nil {
		return nil, err
	}

	return answers, nil
}
