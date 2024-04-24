package graphql_api

import (
	"context"
	"encoding/json"
	"errors"
	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	commonerrors "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/commonerrors"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/agones"
	"github.com/go-faker/faker/v4"
	"github.com/rs/zerolog/log"
	"time"
)

var _ PlatformClient = (*MockGraphQLClient)(nil)

type MockGraphQLClient struct {
	abstract_client.MockAbstractClient
	isAuthenticated bool

	CalledGetUser bool
	GetUserError  bool

	CalledCreateUser bool
	CreateUserError  bool

	CalledUpdateUser bool
	UpdateUserError  bool

	CalledDeleteUser bool
	DeleteUserError  bool

	CalledGetAPIKeys bool
	GetAPIKeysEmpty  bool
	GetAPIKeysError  bool

	CalledCreateAPIKey bool
	CreateAPIKeyError  bool

	CalledDeleteAPIKey bool
	DeleteAPIKeyError  bool

	CalledGetSession bool
	GetSessionError  bool

	CalledCreateSession bool
	CreateSessionError  bool

	CalledUpdateSession bool
	UpdateSessionError  bool

	CalledCreateEvent bool
	CreateEventError  bool

	CalledGetPlatforms bool
	GetPlatformsError  bool

	CalledGetControlTypes bool
	GetControlTypesError  bool

	CalledCreateModuleVersion bool
	CreateModuleVersionError  bool

	CalledGetMultiplayerServerConfigs        bool
	GetMultiplayerServerConfigsError         bool
	GetMultiplayerServerConfigsEmpty         bool
	GetMultiplayerServerConfigsEmptyVersions bool

	CalledGetMultiplayerServerVersions bool
	GetMultiplayerServerVersionsError  bool
	GetMultiplayerServerVersionsEmpty  bool

	CalledGetMultiplayerServerVersion bool
	GetMultiplayerServerVersionError  bool
	GetMultiplayerServerVersionEmpty  bool

	CalledCreateMultiplayerServerVersion bool
	CreateMultiplayerServerVersionError  bool
}

func (m *MockGraphQLClient) GetURL() string {
	return faker.URL()
}

func (m *MockGraphQLClient) Login(username, password string) error {
	m.isAuthenticated = true
	return nil
}

func (m *MockGraphQLClient) ActiveUserID() int {
	return 1
}

func (m *MockGraphQLClient) ActiveOrgID() int {
	return 1
}

func (m *MockGraphQLClient) GetToken() string {
	return faker.UUIDHyphenated()
}

func (m *MockGraphQLClient) SetToken(token string) {
	m.isAuthenticated = true
	return
}

func (m *MockGraphQLClient) SetAPIKey(apiKey string) {
	m.isAuthenticated = true
	return
}

func (m *MockGraphQLClient) IsAuthenticated() bool {
	return m.isAuthenticated
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

	user.ID = 1
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

func (m *MockGraphQLClient) CreateAPIKey(ctx context.Context, input platform.APIKey) (*platform.APIKey, error) {

	m.CalledCreateAPIKey = true

	if m.CreateAPIKeyError {
		return nil, errors.New("error creating api key")
	}

	input.ID = 1
	input.Key = faker.UUIDHyphenated()
	input.CreatedAt = time.Now().UTC()
	input.UpdatedAt = time.Now().UTC()

	return &input, nil
}

func (m *MockGraphQLClient) GetAPIKeys(ctx context.Context, params *APIKeyQueryParams) ([]*platform.APIKey, error) {

	m.CalledGetAPIKeys = true

	if m.GetAPIKeysEmpty {
		return []*platform.APIKey{}, nil
	}

	if m.GetAPIKeysError {
		return nil, errors.New("error getting user")
	}

	if params.UserID == nil {
		params.UserID = &[]int{1}[0]
	}

	return []*platform.APIKey{
		{
			ID:        1,
			UserID:    *params.UserID,
			Key:       faker.UUIDHyphenated(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
	}, nil
}

func (m *MockGraphQLClient) DeleteAPIKey(ctx context.Context, id int) error {

	m.CalledDeleteAPIKey = true

	if m.DeleteAPIKeyError {
		return errors.New("error deleting api key")
	}

	if id <= 0 {
		return errors.New("invalid api key id")
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

func (m *MockGraphQLClient) UpdateSession(ctx context.Context, session Session) (*Session, error) {

	m.CalledUpdateSession = true

	if m.UpdateSessionError {
		return nil, errors.New("error updating session")
	}

	if session.ID <= 0 {
		return nil, errors.New("invalid session id")
	}

	if session.MaxScore > 0 {
		session.ScaledScore = session.RawScore / session.MaxScore
	}

	session.Duration = "1s"

	return &session, nil
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

func (m *MockGraphQLClient) GetPlatforms(ctx context.Context) ([]*Platform, error) {
	m.CalledGetPlatforms = true

	if m.GetPlatformsError {
		return nil, errors.New("error getting platforms")
	}

	return []*Platform{{
		ID:   1,
		Name: "android",
	}}, nil
}

func (m *MockGraphQLClient) GetControlTypes(ctx context.Context) ([]*ControlType, error) {
	m.CalledGetControlTypes = true

	if m.GetControlTypesError {
		return nil, errors.New("error getting control types")
	}

	return []*ControlType{{
		ID:   1,
		Name: "keyboard/mouse",
	}}, nil
}

func (m *MockGraphQLClient) CreateModuleVersion(ctx context.Context, input ModuleVersion) (*ModuleVersion, error) {

	m.CalledCreateModuleVersion = true

	if m.CreateModuleVersionError {
		return nil, errors.New("error creating module version")
	}

	if input.ModuleID == 0 {
		return nil, errors.New("invalid module id")
	}

	if input.LocalFilePath == "" {
		return nil, errors.New("local file path required")
	}

	if input.SemanticVersion == "" {
		return nil, errors.New("invalid version")
	}

	return &ModuleVersion{
		ModuleID:        input.ModuleID,
		SemanticVersion: input.SemanticVersion,
		Package:         input.Package,
	}, nil
}

func (m *MockGraphQLClient) GetMultiplayerServerConfigs(ctx context.Context, params *MultiplayerServerConfigParams) ([]*MultiplayerServerConfigQueryParams, error) {
	m.CalledGetMultiplayerServerConfigs = true

	if m.GetMultiplayerServerConfigsError {
		return nil, errors.New("error getting multiplayer server configs")
	}

	if m.GetMultiplayerServerConfigsEmpty {
		return []*MultiplayerServerConfigQueryParams{}, nil
	}

	if m.GetMultiplayerServerConfigsEmptyVersions {
		return []*MultiplayerServerConfigQueryParams{
			{
				ModuleID: 1,
				Capacity: 5,
			},
		}, nil
	}

	return []*MultiplayerServerConfigQueryParams{
		{
			ModuleID: 1,
			Capacity: 5,
			ServerVersions: []*MultiplayerServerVersion{
				{
					Engine:          "unreal",
					ImageRegistry:   agones.SimpleGameServerImage,
					Status:          "enabled",
					SemanticVersion: "1.0.0",
				},
			},
		},
	}, nil
}

func (m *MockGraphQLClient) GetMultiplayerServerVersions(ctx context.Context, params *MultiplayerServerVersionQueryParams) ([]*MultiplayerServerVersion, error) {

	m.CalledGetMultiplayerServerVersions = true

	if m.GetMultiplayerServerVersionsEmpty {
		return []*MultiplayerServerVersion{}, nil
	}

	if m.GetMultiplayerServerVersionsError {
		return nil, errors.New("error getting multiplayer server versions")
	}

	if params.ModuleID == 0 {
		return nil, errors.New("invalid module id")
	}

	if params.SemanticVersion == "" {
		return nil, errors.New("invalid semantic version")
	}

	return []*MultiplayerServerVersion{
		{
			ModuleID:        1,
			SemanticVersion: "1.0.0",
			Status:          "enabled",
			Engine:          "unreal",
		},
	}, nil
}

func (m *MockGraphQLClient) GetMultiplayerServerVersion(ctx context.Context, versionID int) (*MultiplayerServerVersion, error) {

	m.CalledGetMultiplayerServerVersion = true

	if m.GetMultiplayerServerVersionEmpty {
		return nil, nil
	}

	if m.GetMultiplayerServerVersionError {
		return nil, errors.New("error getting multiplayer server version")
	}

	return &MultiplayerServerVersion{
		ID:              versionID,
		ModuleID:        1,
		SemanticVersion: "1.0.0",
		Status:          "enabled",
		Engine:          "unreal",
	}, nil
}

func (m *MockGraphQLClient) CreateMultiplayerServerVersion(ctx context.Context, input MultiplayerServerVersion) (*MultiplayerServerVersion, error) {
	m.CalledCreateMultiplayerServerVersion = true

	if input.ModuleID == 0 {
		return nil, errors.New("invalid module id")
	}

	if input.ImageRegistry == "" && input.LocalFilePath == "" {
		return nil, errors.New("image or file path required")
	}

	if input.SemanticVersion == "" {
		return nil, errors.New("invalid semantic version")
	}

	if m.CreateMultiplayerServerVersionError {
		return nil, errors.New("error creating multiplayer server version")
	}

	return &MultiplayerServerVersion{
		ModuleID:        input.ModuleID,
		SemanticVersion: input.SemanticVersion,
		Status:          input.Status,
		Engine:          input.Engine,
	}, nil
}
