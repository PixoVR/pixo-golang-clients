package forms

import (
	"io"
)

type FormHandler interface {
	GetResponseFromUser(question *Question) error
	GetSensitiveResponseFromUser(question *Question) error
	Confirm(question *Question) error
	Select(question *Question) error
	SelectID(question *Question) error
	MultiSelect(question *Question) error
	MultiSelectIDs(question *Question) error
	AskQuestions(questions []Question) (map[string]interface{}, error)
	SetReader(reader io.Reader)
	SetWriter(writer io.Writer)
}
