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
	Type           QuestionType
	Options        []Option
	GetOptionsFunc func() ([]Option, error)
	Answer         interface{}
	Key            string
	Prompt         string
	Optional       bool
}

func (q *Question) GetOptions() ([]Option, error) {
	if q.Options == nil && q.GetOptionsFunc != nil {
		options, err := q.GetOptionsFunc()
		if err != nil {
			return nil, err
		}
		return options, nil
	}
	return q.Options, nil
}
