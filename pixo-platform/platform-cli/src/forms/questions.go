package forms

type QuestionType int

const (
	Input QuestionType = iota
	SensitiveInput
	Confirm
	Select
	SelectID
	MultiSelect
	MultiSelectIDs
)

type Question struct {
	Key     string
	Prompt  string
	Type    QuestionType
	Options []Option
	Answer  interface{}
}
