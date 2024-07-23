package legacy

import "time"

type FieldValidationError struct {
	Location string `json:"location"`
	Msg      string `json:"msg"`
	Path     string `json:"path"`
	Type     string `json:"type"`
	Value    string `json:"value"`
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
