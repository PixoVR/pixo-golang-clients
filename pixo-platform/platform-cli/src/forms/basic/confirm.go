package basic

import "strings"

const (
	yes = "yes"
)

func (f *Handler) Confirm(title string) bool {
	response, _ := f.GetResponseFromUser(title)

	if response == "" {
		return false
	}

	isLowercase := strings.ToLower(response) == strings.ToLower(yes)
	isFirstChar := strings.ToLower(response[:1]) == strings.ToLower(yes[:1])
	return isLowercase || isFirstChar
}
