package fancy

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/charmbracelet/huh"
)

func (f *Handler) MultiSelectIDsInput(prompt string, options []forms.Option, response *[]int) huh.Field {
	return huh.NewMultiSelect[int]().
		Options(transformIntOptions(options)...).
		Title(prompt).
		Value(response)
}

func (f *Handler) MultiSelectIDs(prompt string, options []forms.Option, response *[]int) error {
	return f.MultiSelectIDsInput(prompt, options, response).Run()
}

func (f *Handler) MultiSelectInput(prompt string, options []forms.Option, response *[]string) huh.Field {
	return huh.NewMultiSelect[string]().
		Options(transformStringOptions(options)...).
		Title(prompt).
		Value(response)
}

func (f *Handler) MultiSelect(prompt string, options []forms.Option, response *[]string) error {
	return f.MultiSelectInput(prompt, options, response).Run()
}
