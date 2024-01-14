package graphql_api

import (
	"context"
	"encoding/json"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"time"
)

type SessionsClient interface {
	GetSession(ctx context.Context, id int) (*Session, error)
	CreateSession(ctx context.Context, moduleID int, ipAddress, deviceId string) (*Session, error)
	CreateEvent(ctx context.Context, sessionID int, uuid string, eventType string, data string) (*platform.Event, error)
}

type Session struct {
	ID       int `json:"id"`
	UserID   int `json:"userId"`
	OrgID    int `json:"orgId"`
	ModuleID int `json:"moduleId"`

	UUID      string `json:"uuid"`
	IPAddress string `json:"ipAddress"`
	DeviceID  string `json:"deviceId"`

	Module platform.Module  `json:"module"`
	Events []platform.Event `json:"events"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateSessionResponse struct {
	Session Session `json:"createSession"`
}

type SessionResponse struct {
	Session Session `json:"session"`
}

type CreateEventResponse struct {
	Event platform.Event `json:"createEvent"`
}

func (g *GraphQLAPIClient) GetSession(ctx context.Context, id int) (*Session, error) {
	query := `query session($id: ID!) { session(id: $id) { id userId } }`

	variables := map[string]interface{}{
		"id": id,
	}

	res, err := g.gqlClient.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var sessionResponse SessionResponse
	if err = json.Unmarshal(res, &sessionResponse); err != nil {
		return nil, err
	}

	return &sessionResponse.Session, nil
}

func (g *GraphQLAPIClient) CreateSession(ctx context.Context, moduleID int, ipAddress, deviceId string) (*Session, error) {
	query := `mutation createSession($input: SessionInput!) { createSession(input: $input) { id } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"moduleId":  moduleID,
			"ipAddress": ipAddress,
			"deviceId":  deviceId,
		},
	}

	res, err := g.gqlClient.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var sessionResponse CreateSessionResponse
	if err = json.Unmarshal(res, &sessionResponse); err != nil {
		return nil, err
	}

	return &sessionResponse.Session, nil
}

func (g *GraphQLAPIClient) CreateEvent(ctx context.Context, sessionID int, uuid string, eventType string, data string) (*platform.Event, error) {
	query := `mutation createEvent($input: EventInput!) { createEvent(input: $input) { id } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"sessionId": sessionID,
			"uuid":      uuid,
			"eventType": eventType,
			"jsonData":  data,
		},
	}

	res, err := g.gqlClient.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var eventResponse CreateEventResponse
	if err = json.Unmarshal(res, &eventResponse); err != nil {
		return nil, err
	}

	return &eventResponse.Event, nil
}
