package basic

import (
	"errors"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"strings"
)

const (
	yes = "yes"
)

func (f *Handler) Confirm(question *forms.Question) error {
	if question == nil {
		return errors.New("question not provided")
	}

	if err := f.GetResponseFromUser(question); err != nil {
		return err
	}

	var (
		isLowercase bool
		isFirstChar bool
		answer      = question.Answer.(string)
	)

	isLowercase = strings.EqualFold(answer, yes)
	if len(answer) > 0 {
		isFirstChar = strings.EqualFold(answer[:1], yes[:1])
	}
	confirmed := isLowercase || isFirstChar

	question.Answer = confirmed
	return nil
}
