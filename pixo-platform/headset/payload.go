package headset

// Payload is the struct that represents common items in the event payload
type Payload struct {
	LessonStatus    *string `json:"lessonStatus,omitempty"`
	LessonStatusOld *string `json:"Lesson_status,omitempty"`

	Score    *float64 `json:"score,omitempty"`
	ScoreOld *float64 `json:"Score,omitempty"`

	ScoreMax    *float64 `json:"scoreMax,omitempty"`
	ScoreMaxOld *float64 `json:"Score_max,omitempty"`
	Success     *bool    `json:"success,omitempty"`

	TimeJoined *string `json:"TimeJoined,omitempty"`

	SessionDuration    *float64 `json:"sessionDuration,omitempty"` // this needs to be a float so that we can handle 1 and 1.0
	OldSessionDuration *float64 `json:"sessionduration,omitempty"`

	Result *struct {
		Completion bool   `json:"completion,omitempty"`
		Success    bool   `json:"success,omitempty"`
		Duration   string `json:"duration,omitempty"`
		Score      *struct {
			Raw float64 `json:"raw,omitempty"`
			Max float64 `json:"max,omitempty"`
		} `json:"score,omitempty"`
		Extensions map[string]interface{} `json:"extensions,omitempty"`
	} `json:"result,omitempty"`
	Object *struct {
		ID *string `json:"id,omitempty"`
	} `json:"object,omitempty"`
	Context *struct {
		Registration string                 `json:"registration,omitempty"`
		Revision     string                 `json:"revision,omitempty"`
		Extensions   map[string]interface{} `json:"extensions,omitempty"`
	} `json:"context,omitempty"`
}
