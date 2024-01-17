package graphql_api

import (
	"context"
	"encoding/json"
	"errors"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	commonerrors "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/commonerrors"
	"github.com/rs/zerolog/log"
	"time"
)

type MockGraphQLClient struct {
	CalledCreateUser    bool
	CalledGetSession    bool
	CalledCreateSession bool
	CalledUpdateSession bool
	CalledCreateEvent   bool
}

func (m *MockGraphQLClient) CreateUser(ctx context.Context, user platform.User) (*platform.User, error) {

	m.CalledCreateUser = true

	if user.Username == "" {
		return nil, commonerrors.ErrorRequired("username")
	}

	if user.Password == "" {
		return nil, commonerrors.ErrorRequired("password")
	}

	if user.OrgID <= 0 {
		return nil, errors.New("invalid org id")
	}

	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()
	return &user, nil
}

func (m *MockGraphQLClient) GetSession(ctx context.Context, id int) (*Session, error) {

	m.CalledGetSession = true

	if id <= 0 {
		return nil, errors.New("invalid session id")
	}

	return &Session{
		ID:        id,
		UserID:    1,
		ModuleID:  1,
		IPAddress: "127.0.0.1",
		DeviceID:  "1234567890",
	}, nil
}

func (m *MockGraphQLClient) CreateSession(ctx context.Context, moduleID int, ipAddress, deviceId string) (*Session, error) {

	m.CalledCreateSession = true

	if moduleID <= 0 {
		return nil, errors.New("invalid module id")
	}

	if ipAddress == "" {
		return nil, commonerrors.ErrorRequired("ip address")
	}

	return &Session{
		ID:        1,
		UserID:    1,
		ModuleID:  moduleID,
		IPAddress: ipAddress,
		DeviceID:  deviceId,
	}, nil
}

func (m *MockGraphQLClient) UpdateSession(ctx context.Context, id int, status string, completed bool) (*Session, error) {

	m.CalledUpdateSession = true

	if id <= 0 {
		return nil, errors.New("invalid session id")
	}

	return &Session{
		ID:        id,
		UserID:    1,
		ModuleID:  1,
		IPAddress: "127.0.0.1",
	}, nil
}

func (m *MockGraphQLClient) CreateEvent(ctx context.Context, sessionID int, uuid string, eventType string, data string) (*platform.Event, error) {

	m.CalledCreateEvent = true

	if sessionID <= 0 {
		return nil, InvalidSessionError
	}

	if eventType == "" {
		return nil, commonerrors.ErrorRequired("event type")
	}

	var jsonData platform.EventResult
	if err := json.Unmarshal([]byte(data), &jsonData); err != nil {
		log.Error().Err(err).Msg("error unmarshalling event data")
	}

	return &platform.Event{
		ID:        1,
		SessionID: sessionID,
		UUID:      uuid,
		EventType: eventType,
		Data:      jsonData,
		CreatedAt: time.Now().UTC(),
	}, nil
}
