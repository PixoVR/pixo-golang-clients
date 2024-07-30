package fancy

import (
	"errors"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/charmbracelet/huh"
)

func (f *Handler) SelectIDInput(question *forms.Question) huh.Field {
	return huh.NewSelect[int]().
		Options(transformIntOptions(question.Options)...).
		Title(question.Prompt).
		Value(question.Answer.(*int))
}

func (f *Handler) SelectID(question *forms.Question) (err error) {
	if question == nil {
		return errors.New("question not provided")
	}
	return f.SelectIDInput(question).Run()
}

func (f *Handler) SelectInput(question *forms.Question) huh.Field {
	return huh.NewSelect[string]().
		Options(transformStringOptions(question.Options)...).
		Title(question.Prompt).
		Value(question.Answer.(*string))
}

func (f *Handler) Select(question *forms.Question) (err error) {
	if question == nil {
		return errors.New("question not provided")
	}
	return f.SelectInput(question).Run()
}

func (f *Handler) MultiSelectIDsInput(question *forms.Question) huh.Field {
	return huh.NewMultiSelect[int]().
		Options(transformIntOptions(question.Options)...).
		Title(question.Prompt).
		Value(question.Answer.(*[]int))
}

func (f *Handler) MultiSelectIDs(question *forms.Question) (err error) {
	if question == nil {
		return errors.New("question not provided")
	}
	return f.MultiSelectIDsInput(question).Run()
}

func (f *Handler) MultiSelectInput(question *forms.Question) huh.Field {
	return huh.NewMultiSelect[string]().
		Options(transformStringOptions(question.Options)...).
		Title(question.Prompt).
		Value(question.Answer.(*[]string))
}

func (f *Handler) MultiSelect(question *forms.Question) (err error) {
	if question == nil {
		return errors.New("question not provided")
	}
	return f.MultiSelectInput(question).Run()
}
