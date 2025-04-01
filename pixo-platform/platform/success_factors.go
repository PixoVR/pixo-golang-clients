package platform

import (
	"context"
	"time"
)

type LearningHistory struct {
	ID                   int     `json:"id"`
	UserID               int     `json:"userId"`
	ModuleID             int     `json:"moduleId"`
	OrgID                int     `json:"orgId"`
	UserEmail            string  `json:"userEmail"`
	ModuleExternalID     string  `json:"moduleExternalId"`
	ModuleLearningType   string  `json:"moduleLearningType"`
	ModuleCompleteStatus string  `json:"moduleCompleteStatus"`
	Score                float64 `json:"score"`
	SessionDuration      int     `json:"sessionDuration"`
	CreatedAt            string  `json:"createdAt"`
}

type CourseData struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	ImageName           string `json:"imageName"`
	ImagePath           string `json:"imagePath"`
	ThumbnailPath       string `json:"thumbnailPath"`
	PDFPath             string `json:"pdfPath"`
	PDFName             string `json:"pdfName"`
	ShortDesc           string `json:"shortDesc"`
	LongDesc            string `json:"longDesc"`
	Details             string `json:"details"`
	Public              bool   `json:"public"`
	Demo                bool   `json:"demo"`
	Status              string `json:"status"`
	Deleted             bool   `json:"deleted"`
	Industry            string `json:"industry"`
	DistributorID       int    `json:"distributorId"`
	Developer           string `json:"developer"`
	UpdatedBy           string `json:"updatedBy"`
	IsMultiplayer       bool   `json:"isMultiplayer"`
	PassingScoreEnabled bool   `json:"passingScoreEnabled"`
	ExternalID          string `json:"externalId"`
	LearningType        string `json:"learningType"`
	CompleteStatus      string `json:"completeStatus"`
	CreatedAt           string `json:"createdAt"`
	UpdatedAt           string `json:"updatedAt"`
}

type OrgSuccessFactor struct {
	ID    int    `json:"id" gorm:"primaryKey"`
	Path  string `json:"path" gorm:"column:path"`
	OrgID int    `json:"orgId" gorm:"column:org_id"`
}

type LearningHistoryParams struct {
	OrgID     int       `json:"orgId"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

type GetLearningHistoryResponse struct {
	LearningHistoryRecords []LearningHistory `json:"learningHistory"`
}

func (p *clientImpl) GetLearningHistoryRecords(ctx context.Context, input LearningHistoryParams) ([]LearningHistory, error) {
	query := `
query learningHistory($params: SuccessFactorsParams!) {
    learningHistory(params: $params) {
        id
        moduleId
        userId
        orgId
        userEmail
        score
        sessionDuration
        moduleExternalId
        moduleLearningType
        moduleCompleteStatus
        createdAt
    }
}`

	params := map[string]interface{}{
		"params": map[string]interface{}{
			"orgId":     input.OrgID,
			"startDate": input.StartDate,
			"endDate":   input.EndDate,
		},
	}

	var res GetLearningHistoryResponse
	if err := p.Exec(ctx, query, &res, params); err != nil {
		return nil, err
	}

	return res.LearningHistoryRecords, nil
}

type GetCourseDataResponse struct {
	CourseDataRecords []CourseData `json:"courseData"`
}

func (p *clientImpl) GetCourseDataRecords(ctx context.Context, orgID int) ([]CourseData, error) {
	query := `
query courseData($orgId: ID!){
    courseData(orgId: $orgId) {
        id
        name
        description
        imageName
        imagePath
        thumbnailPath
        pdfPath
        pdfName
        shortDesc
        longDesc
        details
        public
        demo
        status
        deleted
        industry
        distributorId
        developer
        updatedBy
        isMultiplayer
        passingScoreEnabled
        externalId
        learningType
        completeStatus
        createdAt
        updatedAt
    }
}`

	params := map[string]interface{}{
		"orgId": orgID,
	}

	var res GetCourseDataResponse
	if err := p.Exec(ctx, query, &res, params); err != nil {
		return nil, err
	}

	return res.CourseDataRecords, nil
}

type OrgSuccessFactorsResponse struct {
	OrgSuccessFactors []OrgSuccessFactor `json:"orgSuccessFactors"`
}

func (p *clientImpl) GetOrgSuccessFactors(ctx context.Context) ([]OrgSuccessFactor, error) {
	query := `
query orgSuccessFactors {
    orgSuccessFactors {
        id
        orgId
        path
    }
}
`

	var res OrgSuccessFactorsResponse
	if err := p.Exec(ctx, query, &res, nil); err != nil {
		return nil, err
	}

	return res.OrgSuccessFactors, nil
}
