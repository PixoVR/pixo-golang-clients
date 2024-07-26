package forms

import "io"

type FormHandler interface {
	GetResponseFromUser(prompt string, response *string) error
	GetSensitiveResponseFromUser(prompt string, response *string) error
	Select(prompt string, options []Option, response *string) error
	SelectID(prompt string, options []Option, response *int) error
	MultiSelect(prompt string, options []Option, response *[]string) error
	MultiSelectIDs(prompt string, options []Option, response *[]int) error
	Confirm(prompt string, response *bool) error
	AskQuestions(questions []Question) (map[string]interface{}, error)
	SetReader(reader io.Reader)
	SetWriter(writer io.Writer)
}
