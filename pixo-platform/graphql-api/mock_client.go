package graphql_api

import (
	"context"
	"encoding/json"
	"errors"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	commonerrors "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/commonerrors"
	"github.com/go-faker/faker/v4"
	"github.com/rs/zerolog/log"
	"time"
)

type MockGraphQLClient struct {
	CalledGetUser bool
	GetUserError  bool

	CalledCreateUser bool
	CreateUserError  bool

	CalledUpdateUser bool
	UpdateUserError  bool

	CalledDeleteUser bool
	DeleteUserError  bool

	CalledGetSession bool
	GetSessionError  bool

	CalledCreateSession bool
	CreateSessionError  bool

	CalledUpdateSession bool
	UpdateSessionError  bool

	CalledCreateEvent bool
	CreateEventError  bool
}

func (m *MockGraphQLClient) GetUserByUsername(ctx context.Context, username string) (*platform.User, error) {

	m.CalledGetUser = true

	if m.GetUserError {
		return nil, errors.New("error getting user")
	}

	if username == "" {
		return nil, commonerrors.ErrorRequired("username")
	}

	return &platform.User{
		ID:        1,
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		Username:  username,
		OrgID:     1,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}

func (m *MockGraphQLClient) CreateUser(ctx context.Context, user platform.User) (*platform.User, error) {

	m.CalledCreateUser = true

	if m.CreateUserError {
		return nil, errors.New("error creating user")
	}

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

func (m *MockGraphQLClient) UpdateUser(ctx context.Context, user platform.User) (*platform.User, error) {

	m.CalledUpdateUser = true

	if m.UpdateUserError {
		return nil, errors.New("error updating user")
	}

	if user.ID <= 0 {
		return nil, errors.New("invalid user id")
	}

	if user.Username == "" {
		return nil, commonerrors.ErrorRequired("username")
	}

	if user.OrgID <= 0 {
		return nil, errors.New("invalid org id")
	}

	user.UpdatedAt = time.Now().UTC()
	return &user, nil
}

func (m *MockGraphQLClient) DeleteUser(ctx context.Context, id int) error {

	m.CalledDeleteUser = true

	if m.DeleteUserError {
		return errors.New("error deleting user")
	}

	if id <= 0 {
		return errors.New("invalid user id")
	}

	return nil
}

func (m *MockGraphQLClient) GetSession(ctx context.Context, id int) (*Session, error) {

	m.CalledGetSession = true

	if m.GetSessionError {
		return nil, errors.New("error getting session")
	}

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

	if m.CreateSessionError {
		return nil, errors.New("error creating session")
	}

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

	if m.UpdateSessionError {
		return nil, errors.New("error updating session")
	}

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

	if m.CreateEventError {
		return nil, errors.New("error creating event")
	}

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
