package fancy

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/charmbracelet/huh"
)

func (f *Handler) ConfirmField(question *forms.Question) huh.Field {
	return huh.NewConfirm().
		Title(question.Prompt).
		Affirmative("Yes").
		Negative("No").
		Value(question.Answer.(*bool))
}

func (f *Handler) Confirm(question *forms.Question) error {
	return f.ConfirmField(question).Run()
}
