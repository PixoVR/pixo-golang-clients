package fancy

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/charmbracelet/huh"
)

func (f *Handler) MultiSelectIDsInput(prompt string, options []forms.Option, response []int) huh.Field {
	return huh.NewMultiSelect[int]().
		Options(transformIntOptions(options)...).
		Title(prompt).
		Value(&response)
}

func (f *Handler) MultiSelectIDs(prompt string, options []forms.Option) ([]int, error) {
	var response []int
	err := f.MultiSelectIDsInput(prompt, options, response).Run()
	return response, err
}

func (f *Handler) MultiSelectInput(prompt string, options []forms.Option, response []string) huh.Field {
	return huh.NewMultiSelect[string]().
		Options(transformStringOptions(options)...).
		Title(prompt).
		Value(&response)
}

func (f *Handler) MultiSelect(prompt string, options []forms.Option) ([]string, error) {
	var response []string
	err := f.MultiSelectInput(prompt, options, response).Run()
	return response, err
}
