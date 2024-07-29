package platform

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

type Session struct {
	ID int `json:"id,omitempty"`

	UUID        string  `json:"uuid,omitempty"`
	IPAddress   string  `json:"ipAddress,omitempty"`
	DeviceID    string  `json:"deviceId,omitempty"`
	RawScore    float64 `json:"rawScore,omitempty"`
	MaxScore    float64 `json:"maxScore,omitempty"`
	ScaledScore float64 `json:"scaledScore,omitempty"`
	Status      string  `json:"status,omitempty"`
	Completed   bool    `json:"completed,omitempty"`
	CompletedAt string  `json:"completedAt,omitempty"`
	Duration    string  `json:"duration,omitempty"`

	UserID   int     `json:"userId,omitempty"`
	User     User    `json:"user,omitempty"`
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
	query := `query session($id: ID!) { session(id: $id) { id userId user { orgId } moduleId } }`

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

	query := `mutation createSession($input: SessionInput!) { createSession(input: $input) { id uuid userId user { orgId } moduleId module { id abbreviation } } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"uuid":     session.UUID,
			"moduleId": session.ModuleID,
			"deviceId": session.DeviceID,
		},
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
	query := `mutation updateSession($input: SessionInput!) { updateSession(input: $input) { id rawScore maxScore scaledScore completedAt duration moduleId userId user { orgId } } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"id":        session.ID,
			"status":    session.Status,
			"completed": session.Completed,
			"rawScore":  session.RawScore,
			"maxScore":  session.MaxScore,
		},
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
