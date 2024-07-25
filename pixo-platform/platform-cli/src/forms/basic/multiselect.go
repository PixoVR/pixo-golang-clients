package basic

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"strconv"
	"strings"
)

func (f *Handler) MultiSelect(prompt string, options []forms.Option) ([]string, error) {
	if _, err := f.writer.Write([]byte(prompt)); err != nil {
		return nil, err
	}
	_, _ = f.writer.Write([]byte("\n"))

	for _, option := range options {
		if _, err := f.writer.Write([]byte(option.Label)); err != nil {
			return nil, err
		}
		_, _ = f.writer.Write([]byte("\n"))
	}

	response, err := f.ReadLine()
	if err != nil {
		return nil, err
	}

	return strings.Split(strings.Trim(response, "\n"), ","), nil
}

func (f *Handler) MultiSelectIDs(prompt string, options []forms.Option) ([]int, error) {
	answers, err := f.MultiSelect(prompt, options)
	if err != nil {
		return nil, err
	}

	ids := make([]int, len(answers))
	for i, answer := range answers {
		for _, option := range options {
			if answer == option.Label {
				id, err := strconv.Atoi(option.Value)
				if err != nil {
					return nil, err
				}
				ids[i] = id
			}
		}
	}

	return ids, nil
}
