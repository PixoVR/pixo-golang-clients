package basic

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"strconv"
	"strings"
)

func (f *Handler) printOptions(prompt string, options []forms.Option) {
	if _, err := f.writer.Write([]byte(prompt)); err != nil {
		return
	}

	_, _ = f.writer.Write([]byte("\n"))

	for _, option := range options {
		if _, err := f.writer.Write([]byte(option.Label)); err != nil {
			return
		}
		_, _ = f.writer.Write([]byte("\n"))
	}
}

func (f *Handler) Select(question *forms.Question) error {
	if question == nil {
		return fmt.Errorf("question not provided")
	}

	f.printOptions(question.Prompt, question.Options)

	line, err := f.ReadLine()
	if err != nil {
		return err
	}

	if line == "" {
		return fmt.Errorf("%s not provided", question.Prompt)
	}

	question.Answer = strings.Trim(line, "\n")
	return nil
}

func (f *Handler) SelectID(question *forms.Question) error {
	if question == nil {
		return fmt.Errorf("question not provided")
	}

	f.printOptions(question.Prompt, question.Options)

	line, err := f.ReadLine()
	if err != nil {
		return err
	}

	if line == "" {
		return fmt.Errorf("%s not provided", question.Prompt)
	}

	for _, option := range question.Options {
		if line == option.Label {
			id, err := strconv.Atoi(option.Value)
			if err != nil {
				return err
			}
			question.Answer = id
			return nil
		}
	}

	return nil
}

func (f *Handler) MultiSelect(question *forms.Question) error {
	if _, err := f.writer.Write([]byte(question.Prompt)); err != nil {
		return err
	}

	_, _ = f.writer.Write([]byte("\n"))

	for _, option := range question.Options {
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
		return fmt.Errorf("%s not provided", question.Prompt)
	}

	question.Answer = strings.Split(strings.Trim(line, "\n"), ",")
	return nil
}

func (f *Handler) MultiSelectIDs(question *forms.Question) error {
	if err := f.MultiSelect(question); err != nil {
		return err
	}

	answers := question.Answer.([]string)
	if len(answers) == 0 {
		return fmt.Errorf("%s not provided", question.Prompt)
	}

	ids := make([]int, len(answers))
	for i, answer := range answers {
		for _, option := range question.Options {
			if answer == option.Label {
				id, err := strconv.Atoi(option.Value)
				if err != nil {
					return err
				}
				ids[i] = id
			}
		}
	}

	question.Answer = ids
	return nil
}
