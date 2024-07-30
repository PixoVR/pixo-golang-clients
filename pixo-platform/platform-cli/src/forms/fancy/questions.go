package fancy

import (
	"context"
	"errors"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/charmbracelet/huh"
)

func (f *Handler) AskQuestions(questions []forms.Question) (map[string]interface{}, error) {
	answers := make(map[string]interface{})
	var groupItems []huh.Field

	for _, question := range questions {
		if question.Prompt == "" {
			question.Prompt = forms.CleanPrompt(question.Key)
		}

		if err := question.GetOptions(context.TODO()); err != nil {
			return nil, err
		}

		switch question.Type {
		case forms.Input:
			question.Answer = new(string)
			answers[question.Key] = question.Answer
			groupItems = append(groupItems, f.InputField(&question))
		case forms.SensitiveInput:
			question.Answer = new(string)
			answers[question.Key] = question.Answer
			groupItems = append(groupItems, f.SensitiveInputField(&question))
		case forms.Confirm:
			question.Answer = new(bool)
			answers[question.Key] = question.Answer
			groupItems = append(groupItems, f.ConfirmField(&question))
		case forms.Select:
			question.Answer = new(string)
			answers[question.Key] = question.Answer
			groupItems = append(groupItems, f.SelectInput(&question))
		case forms.SelectID:
			question.Answer = new(int)
			answers[question.Key] = question.Answer
			groupItems = append(groupItems, f.SelectIDInput(&question))
		case forms.MultiSelect:
			question.Answer = new([]string)
			answers[question.Key] = question.Answer
			groupItems = append(groupItems, f.MultiSelectInput(&question))
		case forms.MultiSelectIDs:
			question.Answer = new([]int)
			answers[question.Key] = question.Answer
			groupItems = append(groupItems, f.MultiSelectIDsInput(&question))
		default:
			return nil, errors.New("unknown question type")
		}
	}

	if len(groupItems) > 0 {
		if err := huh.NewForm(huh.NewGroup(groupItems...)).Run(); err != nil {
			return nil, err
		}
	}

	return answers, nil
}
