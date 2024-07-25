package basic

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"strconv"
	"strings"
)

func (f *Handler) MultiSelect(prompt string, options []forms.Option, response *[]string) error {
	if _, err := f.writer.Write([]byte(prompt)); err != nil {
		return err
	}

	_, _ = f.writer.Write([]byte("\n"))

	for _, option := range options {
		if _, err := f.writer.Write([]byte(option.Label)); err != nil {
			return err
		}
		_, _ = f.writer.Write([]byte("\n"))
	}

	line, err := f.ReadLine()
	if err != nil {
		return err
	}

	if line == "" {
		return fmt.Errorf("%s not provided", prompt)
	}

	splitTrimmedLine := strings.Split(strings.Trim(line, "\n"), ",")
	if response != nil {
		*response = splitTrimmedLine
	}
	return nil
}

func (f *Handler) MultiSelectIDs(prompt string, options []forms.Option, response *[]int) error {
	var answers []string
	if err := f.MultiSelect(prompt, options, &answers); err != nil {
		return err
	}
	if len(answers) == 0 {
		return fmt.Errorf("%s not provided", prompt)
	}

	ids := make([]int, len(answers))
	for i, answer := range answers {
		for _, option := range options {
			if answer == option.Label {
				id, err := strconv.Atoi(option.Value)
				if err != nil {
					return err
				}
				ids[i] = id
			}
		}
	}

	if response != nil {
		*response = ids
	}
	return nil
}
