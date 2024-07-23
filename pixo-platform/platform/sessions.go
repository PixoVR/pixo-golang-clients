package platform

import (
	"context"
	"encoding/json"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/legacy"
	"time"
)

type Session struct {
	ID int `json:"id"`

	UUID        string  `json:"uuid"`
	IPAddress   string  `json:"ipAddress"`
	DeviceID    string  `json:"deviceId"`
	RawScore    float64 `json:"rawScore"`
	MaxScore    float64 `json:"maxScore"`
	ScaledScore float64 `json:"scaledScore"`
	Status      string  `json:"status"`
	Completed   bool    `json:"completed"`
	CompletedAt string  `json:"completedAt"`
	Duration    string  `json:"duration"`

	UserID   int              `json:"userId"`
	User     platform.User    `json:"user"`
	OrgID    int              `json:"orgId"`
	Org      Org              `json:"org"`
	ModuleID int              `json:"moduleId"`
	Module   platform.Module  `json:"module"`
	Events   []platform.Event `json:"events"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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

type CreateEventResponse struct {
	Event platform.Event `json:"createEvent"`
}

func (g *PlatformAPIClient) GetSession(ctx context.Context, id int) (*Session, error) {
	query := `query session($id: ID!) { session(id: $id) { id userId user { orgId } moduleId } }`

	variables := map[string]interface{}{
		"id": id,
	}

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var sessionResponse SessionResponse
	if err = json.Unmarshal(res, &sessionResponse); err != nil {
		return nil, err
	}

	return &sessionResponse.Session, nil
}

func (g *PlatformAPIClient) CreateSession(ctx context.Context, moduleID int, ipAddress, deviceId string) (*Session, error) {
	query := `mutation createSession($input: SessionInput!) { createSession(input: $input) { id userId user { orgId } moduleId } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"moduleId":  moduleID,
			"ipAddress": ipAddress,
			"deviceId":  deviceId,
		},
	}

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var sessionResponse CreateSessionResponse
	if err = json.Unmarshal(res, &sessionResponse); err != nil {
		return nil, err
	}

	return &sessionResponse.Session, nil
}

func (g *PlatformAPIClient) UpdateSession(ctx context.Context, session Session) (*Session, error) {
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

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var sessionResponse UpdateSessionResponse
	if err = json.Unmarshal(res, &sessionResponse); err != nil {
		return nil, err
	}

	return &sessionResponse.Session, nil
}

func (g *PlatformAPIClient) CreateEvent(ctx context.Context, sessionID int, uuid string, eventType string, data string) (*platform.Event, error) {
	query := `mutation createEvent($input: EventInput!) { createEvent(input: $input) { id } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"sessionId": sessionID,
			"uuid":      uuid,
			"eventType": eventType,
			"jsonData":  data,
		},
	}

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var eventResponse CreateEventResponse
	if err = json.Unmarshal(res, &eventResponse); err != nil {
		return nil, err
	}

	return &eventResponse.Event, nil
}
