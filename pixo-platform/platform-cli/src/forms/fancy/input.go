package fancy

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/charmbracelet/huh"
)

func (f *Handler) InputField(question *forms.Question) huh.Field {
	if question.Answer == nil {
		question.Answer = new(string)
	}
	return huh.NewInput().
		Title(question.Prompt).
		Value(question.Answer.(*string))
}

func (f *Handler) GetResponseFromUser(question *forms.Question) error {
	return f.InputField(question).Run()
}

func (f *Handler) SensitiveInputField(question *forms.Question) huh.Field {
	if question.Answer == nil {
		question.Answer = new(string)
	}
	return huh.NewInput().
		Title(question.Prompt).
		EchoMode(huh.EchoModePassword).
		Value(question.Answer.(*string))
}

func (f *Handler) GetSensitiveResponseFromUser(question *forms.Question) error {
	return f.SensitiveInputField(question).Run()
}
