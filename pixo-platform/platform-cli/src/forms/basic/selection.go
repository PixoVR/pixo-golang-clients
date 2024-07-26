package basic

import (
	"errors"
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

	if question.Options == nil && question.GetOptionsFunc != nil {
		options, err := question.GetOptionsFunc()
		if err != nil {
			return err
		}
		question.Options = options
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
			question.Answer = option.Label
			return nil
		}
	}

	return errors.New("invalid option")
}

func (f *Handler) SelectID(question *forms.Question) error {
	if question == nil {
		return fmt.Errorf("question not provided")
	}

	options, err := question.GetOptions()
	if err != nil {
		return err
	}

	f.printOptions(question.Prompt, options)

	line, err := f.ReadLine()
	if err != nil {
		return err
	}

	if line == "" {
		return fmt.Errorf("%s not provided", question.Prompt)
	}

	for _, option := range options {
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
	if question == nil {
		return fmt.Errorf("question not provided")
	}

	options, err := question.GetOptions()
	if err != nil {
		return err
	}

	f.printOptions(question.Prompt, options)

	line, err := f.ReadLine()
	if err != nil {
		return err
	}

	if line == "" {
		return fmt.Errorf("%s not provided", question.Prompt)
	}

	selectedOptions := strings.Split(strings.Trim(line, "\n"), ",")

	for _, selectedOption := range selectedOptions {
		found := false
		for _, option := range options {
			if selectedOption == option.Label {
				found = true
				break
			}
		}

		if !found {
			return errors.New("invalid option")
		}
	}

	question.Answer = forms.StringSlice(selectedOptions)
	return nil
}

func (f *Handler) MultiSelectIDs(question *forms.Question) error {
	if err := f.MultiSelect(question); err != nil {
		return err
	}

	answers := forms.StringSlice(question.Answer)
	if len(answers) == 0 {
		return fmt.Errorf("%s not provided", question.Prompt)
	}

	options, err := question.GetOptions()
	if err != nil {
		return err
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

	question.Answer = forms.IntSlice(ids)
	return nil
}
