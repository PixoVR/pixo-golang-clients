package primary_api

import "time"

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User struct {
		Token string `json:"authToken"`
		Role  string `json:"role"`
	}
}
type FieldValidationError struct {
	Location string `json:"location"`
	Msg      string `json:"msg"`
	Path     string `json:"path"`
	Type     string `json:"type"`
	Value    string `json:"value"`
}

type MultiplayerServerVersion struct {
	Status        string `json:"status"`
	ImageRegistry string `json:"imageRegistry"`
}

type FieldValidationErrors struct {
	OrgID     *FieldValidationError `json:"orgID"`
	StartDate *FieldValidationError `json:"startDate"`
	EndDate   *FieldValidationError `json:"endDate"`
}

type CourseDataResponse struct {
	Message    string                 `json:"message"`
	StatusCode int                    `json:"statusCode"`
	Data       []CourseDataResult     `json:"data"`
	Errors     *FieldValidationErrors `json:"errors"`
	Error      bool                   `json:"error"`
}

type LearningHistoryResponse struct {
	Message    string                  `json:"message"`
	StatusCode int                     `json:"statusCode"`
	Data       []LearningHistoryResult `json:"data"`
	Errors     *FieldValidationErrors  `json:"errors"`
	Error      bool                    `json:"error"`
}

type CourseDataResult struct {
	ID                  int         `json:"id"`
	Name                string      `json:"name"`
	Description         string      `json:"description"`
	ImageName           interface{} `json:"imageName"`
	ImagePath           interface{} `json:"imagePath"`
	ThumbnailPath       string      `json:"thumbnailPath"`
	PdfPath             interface{} `json:"pdfPath"`
	PdfName             interface{} `json:"pdfName"`
	ShortDesc           string      `json:"shortDesc"`
	LongDesc            string      `json:"longDesc"`
	Details             interface{} `json:"details"`
	Public              bool        `json:"public"`
	Demo                bool        `json:"demo"`
	Status              string      `json:"status"`
	Deleted             bool        `json:"deleted"`
	Industry            string      `json:"industry"`
	DistributorID       int         `json:"distributorId"`
	Developer           string      `json:"developer"`
	UpdatedBy           string      `json:"updatedBy"`
	IsMultiplayer       bool        `json:"isMultiplayer"`
	PassingScoreEnabled bool        `json:"passingScoreEnabled"`
	OrgModule           []OrgModule `json:"orgModule"`
	CreatedAt           time.Time   `json:"createdAt"`
	UpdatedAt           time.Time   `json:"updatedAt"`
}

type LearningHistoryResult struct {
	ID        int       `json:"id"`
	ModuleID  int       `json:"moduleId"`
	UserID    int       `json:"userId"`
	Module    Module    `json:"module"`
	User      User      `json:"user"`
	Events    []Event   `json:"events"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type User struct {
	ID         int    `json:"id"`
	First      string `json:"first"`
	Last       string `json:"last"`
	Email      string `json:"email"`
	ExternalID string `json:"externalId"`
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
	ID            int         `json:"id"`
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	ImageName     interface{} `json:"imageName"`
	ImagePath     interface{} `json:"imagePath"`
	ThumbnailPath string      `json:"thumbnailPath"`
	ShortDesc     string      `json:"shortDesc"`
	UpdatedAt     time.Time   `json:"updatedAt"`
	OrgModule     []OrgModule `json:"orgModule"`
}

type Event struct {
	ID          int         `json:"id"`
	SessionID   int         `json:"sessionId"`
	UserID      int         `json:"userId"`
	EventType   string      `json:"eventType"`
	CreatedAt   time.Time   `json:"createdAt"`
	EventResult EventResult `json:"jsonData"`
}

type EventResult struct {
	Score           float32 `json:"score"`
	SessionDuration int     `json:"sessionDuration"`
}

type Org struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type OrgResponse struct {
	Error      bool                   `json:"error"`
	Message    string                 `json:"message"`
	StatusCode *int                   `json:"statusCode"`
	Orgs       *[]Org                 `json:"orgs"`
	Errors     *FieldValidationErrors `json:"errors"`
}
