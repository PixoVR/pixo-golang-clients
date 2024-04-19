package primary_api

import "time"

type JSONEvent struct {
	ID              string  `json:"id,omitempty"`
	SessionDuration float64 `json:"sessionDuration,omitempty"`
	LessonStatus    *string `json:"lessonStatus,omitempty"`
	ModuleName      string  `json:"moduleName,omitempty"`
	Actor           struct {
		Name *string `json:"name,omitempty"`
		Mbox *string `json:"mBox,omitempty"`
	} `json:"actor,omitempty"`
	Verb struct {
		ID      *string `json:"id,omitempty"`
		Display struct {
			EN *string `json:"en,omitempty"`
		} `json:"display,omitempty"`
	} `json:"verb,omitempty"`
	Object *struct {
		ID *string `json:"id,omitempty"`
	} `json:"object,omitempty"`
	Result *struct {
		Completion bool   `json:"completion,omitempty"`
		Success    bool   `json:"success,omitempty"`
		Duration   string `json:"duration,omitempty"`
		Score      *struct {
			Scaled float32 `json:"scaled,omitempty"`
			Raw    float64 `json:"raw,omitempty"`
			Min    float32 `json:"min,omitempty"`
			Max    float32 `json:"max,omitempty"`
		} `json:"score,omitempty"`
	} `json:"result,omitempty"`
	Context *struct {
		Registration string                 `json:"registration,omitempty"`
		Revision     string                 `json:"revision,omitempty"`
		Extensions   map[string]interface{} `json:"extensions,omitempty"`
	} `json:"context,omitempty"`

	Score       *float64 `json:"score,omitempty"`
	ScoreMin    *float64 `json:"scoreMin,omitempty"`
	ScoreMax    *float64 `json:"scoreMax,omitempty"`
	ScoreScaled *float64 `json:"scoreScaled,omitempty"`
	Success     *bool    `json:"success,omitempty"`
}

type OrgModule struct {
	ID             int       `json:"id"`
	ModuleID       int       `json:"moduleId"`
	OrgID          int       `json:"orgId"`
	Lifetime       bool      `json:"lifetime"`
	ExpirationDate time.Time `json:"expirDate"`
	LearningType   string    `json:"learningType"`
	CompleteStatus string    `json:"completeStatus"`
	ExternalID     string    `json:"externalId"`
}

type Module struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	//Description   string `json:"description,omitempty"`
	//ImageName     string `json:"imageName,omitempty"`
	//ImagePath     string `json:"imagePath,omitempty"`
	//ThumbnailPath string `json:"thumbnailPath,omitempty"`
	//ShortDesc     string `json:"shortDesc,omitempty"`

	GitConfigID int       `json:"gitConfigId,omitempty"`
	GitConfig   GitConfig `json:"gitConfig,omitempty"`

	//OrgModules []OrgModule

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type Event struct {
	ID        int         `json:"id"`
	UUID      string      `json:"uuid"`
	SessionID int         `json:"sessionId"`
	UserID    int         `json:"userId"`
	EventType string      `json:"eventType"`
	Data      EventResult `json:"jsonData"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type EventResult struct {
	Score           float32 `json:"score"`
	SessionDuration int     `json:"sessionDuration"`
}

type Org struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Status    string    `json:"enabled"`
	CreatedAt time.Time `json:"createdAt"`
}

type OrgResponse struct {
	Error      bool                   `json:"error"`
	Message    string                 `json:"message"`
	StatusCode *int                   `json:"statusCode"`
	Orgs       *[]Org                 `json:"orgs"`
	Errors     *FieldValidationErrors `json:"errors"`
}
