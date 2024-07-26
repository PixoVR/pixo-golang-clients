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
	Type        QuestionType
	Options     []Option
	OptionsFunc func() ([]Option, error)
	Answer      interface{}
	Key         string
	Prompt      string
	Optional    bool
}
