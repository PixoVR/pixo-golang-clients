package forms

import "io"

type FormHandler interface {
	GetResponseFromUser(prompt string) (string, error)
	GetSensitiveResponseFromUser(prompt string) (string, error)
	MultiSelect(prompt string, options []Option) ([]string, error)
	MultiSelectIDs(prompt string, options []Option) ([]int, error)
	Confirm(prompt string) bool
	AskQuestions(questions []Question) (map[string]interface{}, error)
	SetReader(reader io.Reader)
	SetWriter(writer io.Writer)
}
