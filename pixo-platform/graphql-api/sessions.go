package graphql_api

import (
	"encoding/json"
	primary_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"time"
)

type Session struct {
	ID       int `json:"id"`
	UserID   int `json:"userId"`
	OrgID    int `json:"orgId"`
	ModuleID int `json:"moduleId"`

	UUID      string `json:"uuid"`
	IPAddress string `json:"ipAddress"`
	DeviceID  string `json:"deviceId"`

	Module primary_api.Module  `json:"module"`
	Events []primary_api.Event `json:"events"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateSessionResponse struct {
	Session Session `json:"createSession"`
}

type SessionResponse struct {
	Session Session `json:"session"`
}

func (g *GraphQLAPIClient) GetSession(id int) (*Session, error) {
	query := `query session($id: ID!) { session(id: $id) { id userId } }`

	variables := map[string]interface{}{
		"id": id,
	}

	res, err := g.ExecRaw(query, variables)
	if err != nil {
		return nil, err
	}

	var sessionResponse SessionResponse
	if err = json.Unmarshal(res, &sessionResponse); err != nil {
		return nil, err
	}

	return &sessionResponse.Session, nil
}

func (g *GraphQLAPIClient) CreateSession(moduleID int, ipAddress, deviceId string) (*Session, error) {
	query := `mutation createSession($input: SessionInput!) { createSession(input: $input) { id } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"moduleId":  moduleID,
			"ipAddress": ipAddress,
			"deviceId":  deviceId,
		},
	}

	res, err := g.ExecRaw(query, variables)
	if err != nil {
		return nil, err
	}

	var sessionResponse CreateSessionResponse
	if err = json.Unmarshal(res, &sessionResponse); err != nil {
		return nil, err
	}

	return &sessionResponse.Session, nil
}
