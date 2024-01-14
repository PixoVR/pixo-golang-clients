package graphql_api

import (
	"encoding/json"
	"errors"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	. "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/commonerrors"
	"time"
)

type MockSessionsClient struct {
	CalledGetSession    bool
	CalledCreateSession bool
	CalledCreateEvent   bool
}

func (m *MockSessionsClient) GetSession(id int) (*Session, error) {

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

func (m *MockSessionsClient) CreateSession(moduleID int, ipAddress, deviceId string) (*Session, error) {

	m.CalledCreateSession = true

	if moduleID <= 0 {
		return nil, errors.New("invalid module id")
	}

	if ipAddress == "" {
		return nil, ErrorRequired("ip address")
	}

	return &Session{
		ID:        1,
		UserID:    1,
		ModuleID:  moduleID,
		IPAddress: ipAddress,
		DeviceID:  deviceId,
	}, nil
}

func (m *MockSessionsClient) CreateEvent(sessionID int, uuid string, eventType string, data string) (*platform.Event, error) {

	m.CalledCreateEvent = true

	if sessionID <= 0 {
		return nil, InvalidSessionError
	}

	if eventType == "" {
		return nil, ErrorRequired("event type")
	}

	if data == "" {
		return nil, ErrorRequired("data")
	}

	var jsonData platform.EventResult
	if err := json.Unmarshal([]byte(data), &jsonData); err != nil {
		return nil, err
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
