package platform

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

type Session struct {
	ID int `json:"id,omitempty"`

	UUID           *string `json:"uuid,omitempty"`
	IPAddress      string  `json:"ipAddress,omitempty"`
	DeviceID       string  `json:"deviceId,omitempty"`
	RawScore       float64 `json:"rawScore,omitempty"`
	MaxScore       float64 `json:"maxScore,omitempty"`
	ScaledScore    float64 `json:"scaledScore,omitempty"`
	Status         string  `json:"status,omitempty"`
	LessonStatus   string  `json:"lessonStatus,omitempty"`
	Scenario       string  `json:"scenario,omitempty"`
	PlayMode       string  `json:"playMode,omitempty"`
	Focus          string  `json:"focus,omitempty"`
	Specialization string  `json:"specialization,omitempty"`
	Completed      bool    `json:"completed,omitempty"`
	CompletedAt    string  `json:"completedAt,omitempty"`
	Duration       string  `json:"duration,omitempty"`

	UserID   int     `json:"userId,omitempty"`
	User     User    `json:"user,omitempty"`
	OrgID    int     `json:"orgId,omitempty"`
	Org      Org     `json:"org,omitempty"`
	ModuleID int     `json:"moduleId,omitempty"`
	Module   Module  `json:"module,omitempty"`
	Events   []Event `json:"events,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type CreateSessionResponse struct {
	Session Session `json:"createSession"`
}

type UpdateSessionResponse struct {
	Session Session `json:"updateSession"`
}

type SessionResponse struct {
	Session Session `json:"session"`
}

func (p *PlatformClient) GetSession(ctx context.Context, id int) (*Session, error) {
	query := `query session($id: ID!) { session(id: $id) { id uuid deviceId status lessonStatus scenario playMode focus specialization rawScore maxScore scaledScore completedAt orgId org { id name } userId user { id orgId firstName lastName } moduleId module { id abbreviation description externalId } } }`

	variables := map[string]interface{}{
		"id": id,
	}

	res, err := p.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var sessionResponse SessionResponse
	if err = json.Unmarshal(res, &sessionResponse); err != nil {
		return nil, err
	}

	return &sessionResponse.Session, nil
}

func (p *PlatformClient) CreateSession(ctx context.Context, session *Session) error {
	if session == nil {
		return errors.New("session is nil")
	}

	query := `mutation createSession($input: SessionInput!) { createSession(input: $input) { id uuid status lessonStatus scenario playMode focus specialization maxScore deviceId userId user { orgId } moduleId module { id abbreviation } } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"uuid":     session.UUID,
			"moduleId": session.ModuleID,
		},
	}

	if session.ModuleID <= 0 {
		return errors.New("module id is required")
	}

	if session.DeviceID != "" {
		variables["input"].(map[string]interface{})["deviceId"] = session.DeviceID
	}

	if session.Scenario != "" {
		variables["input"].(map[string]interface{})["scenario"] = session.Scenario
	}

	if session.PlayMode != "" {
		variables["input"].(map[string]interface{})["playMode"] = session.PlayMode
	}

	if session.Focus != "" {
		variables["input"].(map[string]interface{})["focus"] = session.Focus
	}

	if session.Specialization != "" {
		variables["input"].(map[string]interface{})["specialization"] = session.Specialization
	}

	res, err := p.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return err
	}

	var sessionResponse CreateSessionResponse
	if err = json.Unmarshal(res, &sessionResponse); err != nil {
		return err
	}

	*session = sessionResponse.Session

	return nil
}

func (p *PlatformClient) UpdateSession(ctx context.Context, session Session) (*Session, error) {
	query := `mutation updateSession($input: SessionInput!) { updateSession(input: $input) { id status lessonStatus scenario playMode focus specialization rawScore maxScore scaledScore completedAt duration moduleId userId user { orgId } } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"status":       session.Status,
			"lessonStatus": session.LessonStatus,
			"completed":    session.Completed,
		},
	}

	if session.ID != 0 {
		variables["input"].(map[string]interface{})["id"] = session.ID
	} else if session.UUID != nil {
		variables["input"].(map[string]interface{})["uuid"] = session.UUID
	} else {
		return nil, errors.New("id or uuid is required")
	}

	if session.RawScore != 0 {
		variables["input"].(map[string]interface{})["rawScore"] = session.RawScore
	}

	if session.MaxScore != 0 {
		variables["input"].(map[string]interface{})["maxScore"] = session.MaxScore
	}

	if session.Scenario != "" {
		variables["input"].(map[string]interface{})["scenario"] = session.Scenario
	}

	res, err := p.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var sessionResponse UpdateSessionResponse
	if err = json.Unmarshal(res, &sessionResponse); err != nil {
		return nil, err
	}

	return &sessionResponse.Session, nil
}
