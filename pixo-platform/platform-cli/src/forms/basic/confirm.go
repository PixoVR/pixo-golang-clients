package basic

import "strings"

const (
	yes = "yes"
)

func (f *Handler) Confirm(title string, response *bool) error {
	var answer string
	if err := f.GetResponseFromUser(title, &answer); err != nil {
		return err
	}

	var (
		isLowercase bool
		isFirstChar bool
	)
	isLowercase = strings.ToLower(answer) == strings.ToLower(yes)
	if len(answer) > 0 {
		isFirstChar = strings.ToLower(answer[:1]) == strings.ToLower(yes[:1])
	}
	confirmed := isLowercase || isFirstChar

	if response != nil {
		*response = confirmed
	}
	return nil
}
