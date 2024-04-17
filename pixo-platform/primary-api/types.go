package primary_api

import "time"

type JSONEvent struct {
	ID              string  `json:"id"`
	SessionDuration float64 `json:"sessionDuration"` // this needs to be a float so that we can handle 1 and 1.0
	LessonStatus    *string `json:"lessonStatus"`
	ModuleName      string  `json:"moduleName"`
	Actor           struct {
		Name *string `json:"name"`
		Mbox *string `json:"mBox"`
	} `json:"actor"`
	Verb struct {
		ID      *string `json:"id"`
		Display struct {
			EN *string `json:"en"`
		} `json:"display"`
	} `json:"verb"`
	Object *struct {
		ID *string `json:"id"`
	} `json:"object"`
	Result *struct {
		Completion bool   `json:"completion"`
		Success    bool   `json:"success"`
		Duration   string `json:"duration"`
		Score      *struct {
			Scaled float32 `json:"scaled"`
			Raw    float64 `json:"raw"`
			Min    float32 `json:"min"`
			Max    float32 `json:"max"`
		} `json:"score"`
	} `json:"result"`
	Context *struct {
		Registration string                 `json:"registration"`
		Revision     string                 `json:"revision"`
		Extensions   map[string]interface{} `json:"extensions"`
	} `json:"context"`

	Score       *float64 `json:"score"`
	ScoreMin    *float64 `json:"scoreMin"`
	ScoreMax    *float64 `json:"scoreMax"`
	ScoreScaled *float64 `json:"scoreScaled"`
	Success     *bool    `json:"success"`
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
