package charm

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/forms"
	"github.com/charmbracelet/huh"
	"io"
	"strconv"
)

type FormHandler struct {
	input  io.Reader
	output io.Writer
}

func NewFormHandler() *FormHandler {
	return &FormHandler{}
}

func (f *FormHandler) SetReader(reader io.Reader) {
	f.input = reader
}

func (f *FormHandler) SetWriter(writer io.Writer) {
	f.output = writer
}

func (f *FormHandler) GetResponseFromUser(prompt string) (string, error) {
	var response string
	if err := huh.NewInput().
		Title(prompt).
		Value(&response).
		Run(); err != nil {
		return "", err
	}

	return response, nil
}

func (f *FormHandler) MultiSelectIDs(prompt string, options []forms.Option) ([]int, error) {
	var response []int

	err := huh.NewMultiSelect[int]().
		Options(transformIntOptions(options)...).
		Title(prompt).
		Value(&response).
		Run()

	return response, err
}

func (f *FormHandler) MultiSelect(prompt string, options []forms.Option) ([]string, error) {
	var response []string

	err := huh.NewMultiSelect[string]().
		Options(transformStringOptions(options)...).
		Title(prompt).
		Value(&response).
		Run()

	return response, err
}

func transformStringOptions(options []forms.Option) []huh.Option[string] {
	var transformedOptions []huh.Option[string]
	for _, o := range options {
		transformedOptions = append(transformedOptions, huh.NewOption[string](o.Label, o.Value))
	}

	return transformedOptions
}

func transformIntOptions(options []forms.Option) []huh.Option[int] {
	var transformedOptions []huh.Option[int]
	for _, o := range options {
		val, err := strconv.Atoi(o.Value)
		if err != nil {
			continue
		}
		transformedOptions = append(transformedOptions, huh.NewOption[int](o.Label, val))
	}

	return transformedOptions
}
