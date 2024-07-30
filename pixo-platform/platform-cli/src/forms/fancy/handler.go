package fancy

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/charmbracelet/huh"
	"io"
	"strconv"
)

var _ forms.FormHandler = &Handler{}

type Handler struct {
	input  io.Reader
	output io.Writer
}

func NewFormHandler() *Handler {
	return &Handler{}
}

func (f *Handler) SetReader(reader io.Reader) {
	f.input = reader
}

func (f *Handler) SetWriter(writer io.Writer) {
	f.output = writer
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
