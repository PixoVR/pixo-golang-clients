package platform

import (
	"context"
	"encoding/json"
	"errors"
	abstract "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/legacy"
	commonerrors "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/commonerrors"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/agones"
	"github.com/go-faker/faker/v4"
	"github.com/rs/zerolog/log"
	"time"
)

var _ Client = (*MockClient)(nil)

type MockClient struct {
	abstract.MockAbstractClient

	NumCalledGetUser int
	GetUserError     error

	NumCalledCreateUser int
	CreateUserError     error

	NumCalledUpdateUser int
	UpdateUserError     error

	NumCalledDeleteUser int
	DeleteUserError     error

	NumCalledGetOrg int
	GetOrgError     error

	NumCalledCreateOrg int
	CreateOrgError     error

	NumCalledUpdateOrg int
	UpdateOrgError     error

	NumCalledDeleteOrg int
	DeleteOrgError     error

	NumCalledGetWebhooks int
	GetWebhooksError     error

	NumCalledGetWebhook int
	GetWebhookError     error

	NumCalledCreateWebhook int
	CreateWebhookError     error

	NumCalledUpdateWebhook int
	UpdateWebhookError     error

	NumCalledDeleteWebhook int
	DeleteWebhookError     error

	NumCalledGetAPIKeys int
	GetAPIKeysEmpty     bool
	GetAPIKeysError     error

	NumCalledCreateAPIKey int
	CreateAPIKeyError     error

	NumCalledDeleteAPIKey int
	DeleteAPIKeyError     error

	NumCalledGetSession int
	GetSessionError     error

	NumCalledCreateSession int
	CreateSessionError     error

	NumCalledUpdateSession int
	UpdateSessionError     error

	NumCalledCreateEvent int
	CreateEventError     error

	NumCalledGetPlatforms int
	GetPlatformsError     error

	NumCalledGetControlTypes int
	GetControlTypesError     error

	NumCalledCreateModuleVersion int
	CreateModuleVersionError     error

	NumCalledGetMultiplayerServerConfigs     int
	GetMultiplayerServerConfigsError         error
	GetMultiplayerServerConfigsEmpty         bool
	GetMultiplayerServerConfigsEmptyVersions bool

	NumCalledGetMultiplayerServerVersions int
	GetMultiplayerServerVersionsError     error
	GetMultiplayerServerVersionsEmpty     bool

	NumCalledGetMultiplayerServerVersion int
	GetMultiplayerServerVersionError     error
	GetMultiplayerServerVersionEmpty     bool

	NumCalledCreateMultiplayerServerVersion int
	CreateMultiplayerServerVersionError     error
}

func (m *MockClient) Path() string {
	return "v2"
}

func (m *MockClient) ActiveUserID() int {
	return 1
}

func (m *MockClient) ActiveOrgID() int {
	return 1
}

func (m *MockClient) GetUserByUsername(ctx context.Context, username string) (*platform.User, error) {
	m.NumCalledGetUser++

	if username == "" {
		return nil, commonerrors.ErrorRequired("username")
	}

	if m.GetUserError != nil {
		return nil, m.GetUserError
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

func (m *MockClient) CreateUser(ctx context.Context, user platform.User) (*platform.User, error) {
	m.NumCalledCreateUser++

	if user.Username == "" {
		return nil, commonerrors.ErrorRequired("username")
	}

	if user.Password == "" {
		return nil, commonerrors.ErrorRequired("password")
	}

	if user.OrgID <= 0 {
		return nil, errors.New("invalid org id")
	}

	if m.CreateUserError != nil {
		return nil, m.CreateUserError
	}

	user.ID = 1
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()
	return &user, nil
}

func (m *MockClient) UpdateUser(ctx context.Context, user platform.User) (*platform.User, error) {
	m.NumCalledUpdateUser++

	if user.ID <= 0 {
		return nil, errors.New("invalid user id")
	}

	if user.Username == "" {
		return nil, commonerrors.ErrorRequired("username")
	}

	if user.OrgID <= 0 {
		return nil, errors.New("invalid org id")
	}

	if m.UpdateUserError != nil {
		return nil, m.UpdateUserError
	}

	user.UpdatedAt = time.Now().UTC()
	return &user, nil
}

func (m *MockClient) DeleteUser(ctx context.Context, id int) error {
	m.NumCalledDeleteUser++

	if id <= 0 {
		return errors.New("invalid user id")
	}

	if m.DeleteUserError != nil {
		return m.DeleteUserError
	}

	return nil
}

func (m *MockClient) GetOrg(ctx context.Context, id int) (*Org, error) {
	m.NumCalledGetOrg++

	if id <= 0 {
		return nil, errors.New("invalid org id")
	}

	if m.GetOrgError != nil {
		return nil, m.GetOrgError
	}

	return &Org{
		ID:   1,
		Name: faker.Name(),
	}, nil
}

func (m *MockClient) CreateOrg(ctx context.Context, org Org) (*Org, error) {
	m.NumCalledCreateOrg++

	if org.Name == "" {
		return nil, commonerrors.ErrorRequired("name")
	}

	if m.CreateOrgError != nil {
		return nil, m.CreateOrgError
	}

	org.ID = 1
	org.CreatedAt = time.Now().UTC()
	org.UpdatedAt = time.Now().UTC()
	return &org, nil
}

func (m *MockClient) UpdateOrg(ctx context.Context, org Org) (*Org, error) {
	m.NumCalledUpdateOrg++

	if org.ID <= 0 {
		return nil, errors.New("org id is required")
	}

	if m.UpdateOrgError != nil {
		return nil, m.UpdateOrgError
	}

	org.UpdatedAt = time.Now().UTC()
	return &org, nil
}

func (m *MockClient) DeleteOrg(ctx context.Context, id int) error {
	m.NumCalledDeleteOrg++

	if id <= 0 {
		return errors.New("invalid org id")
	}

	if m.DeleteOrgError != nil {
		return m.DeleteOrgError
	}

	return nil
}

func (m *MockClient) CreateAPIKey(ctx context.Context, input platform.APIKey) (*platform.APIKey, error) {
	m.NumCalledCreateAPIKey++

	if m.CreateAPIKeyError != nil {
		return nil, m.CreateAPIKeyError
	}

	input.ID = 1
	input.Key = faker.UUIDHyphenated()
	input.CreatedAt = time.Now().UTC()
	input.UpdatedAt = time.Now().UTC()

	return &input, nil
}

func (m *MockClient) GetAPIKeys(ctx context.Context, params *APIKeyQueryParams) ([]*platform.APIKey, error) {
	m.NumCalledGetAPIKeys++

	if m.GetAPIKeysError != nil {
		return nil, m.GetAPIKeysError
	}

	if m.GetAPIKeysEmpty {
		return []*platform.APIKey{}, nil
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

func (m *MockClient) DeleteAPIKey(ctx context.Context, id int) error {
	m.NumCalledDeleteAPIKey++

	if id <= 0 {
		return errors.New("invalid api key id")
	}

	if m.DeleteAPIKeyError != nil {
		return m.DeleteAPIKeyError
	}

	return nil
}

func (m *MockClient) GetWebhooks(ctx context.Context, params *WebhookParams) ([]Webhook, error) {
	m.NumCalledGetWebhooks++

	if m.GetWebhooksError != nil {
		return nil, m.GetWebhooksError
	}

	return []Webhook{
		{
			ID:    1,
			OrgID: 1,
			URL:   "http://example.com",
			Token: "token",
		},
		{
			ID:    2,
			OrgID: 2,
			URL:   "http://example-2.com",
			Token: "token-2",
		},
	}, nil
}

func (m *MockClient) GetWebhook(ctx context.Context, id int) (*Webhook, error) {
	m.NumCalledGetWebhook++

	if id <= 0 {
		return nil, errors.New("invalid webhook id")
	}

	if m.GetWebhookError != nil {
		return nil, m.GetWebhookError
	}

	return &Webhook{
		ID:    id,
		OrgID: 1,
		URL:   "http://example.com",
		Token: "token",
	}, nil
}

func (m *MockClient) CreateWebhook(ctx context.Context, input Webhook) (*Webhook, error) {
	m.NumCalledCreateWebhook++

	if input.OrgID <= 0 {
		return nil, errors.New("invalid org id")
	}

	if input.URL == "" {
		return nil, commonerrors.ErrorRequired("url")
	}

	if input.Token == "" {
		return nil, commonerrors.ErrorRequired("token")
	}

	if m.CreateWebhookError != nil {
		return nil, m.CreateWebhookError
	}

	input.ID = 1
	return &input, nil
}

func (m *MockClient) UpdateWebhook(ctx context.Context, input Webhook) (*Webhook, error) {
	m.NumCalledUpdateWebhook++

	if input.ID <= 0 {
		return nil, errors.New("invalid webhook id")
	}

	if input.OrgID <= 0 {
		return nil, errors.New("invalid org id")
	}

	if input.URL == "" {
		return nil, commonerrors.ErrorRequired("url")
	}

	if input.Token == "" {
		return nil, commonerrors.ErrorRequired("token")
	}

	if m.UpdateWebhookError != nil {
		return nil, m.UpdateWebhookError
	}

	return &input, nil
}

func (m *MockClient) DeleteWebhook(ctx context.Context, id int) error {
	m.NumCalledDeleteWebhook++

	if id <= 0 {
		return errors.New("invalid webhook id")
	}

	if m.DeleteWebhookError != nil {
		return m.DeleteWebhookError
	}

	return nil
}

func (m *MockClient) GetSession(ctx context.Context, id int) (*Session, error) {
	m.NumCalledGetSession++

	if id <= 0 {
		return nil, errors.New("invalid session id")
	}

	if m.GetSessionError != nil {
		return nil, m.GetSessionError
	}

	return &Session{
		ID:        id,
		UserID:    1,
		ModuleID:  1,
		IPAddress: "127.0.0.1",
		DeviceID:  "1234567890",
	}, nil
}

func (m *MockClient) CreateSession(ctx context.Context, moduleID int, ipAddress, deviceId string) (*Session, error) {
	m.NumCalledCreateSession++

	if moduleID <= 0 {
		return nil, errors.New("invalid module id")
	}

	if ipAddress == "" {
		return nil, commonerrors.ErrorRequired("ip address")
	}

	if m.CreateSessionError != nil {
		return nil, m.CreateSessionError
	}

	return &Session{
		ID:        1,
		UserID:    1,
		ModuleID:  moduleID,
		IPAddress: ipAddress,
		DeviceID:  deviceId,
	}, nil
}

func (m *MockClient) UpdateSession(ctx context.Context, session Session) (*Session, error) {
	m.NumCalledUpdateSession++

	if session.ID <= 0 {
		return nil, errors.New("invalid session id")
	}

	if m.UpdateSessionError != nil {
		return nil, m.UpdateSessionError
	}

	if session.MaxScore > 0 {
		session.ScaledScore = session.RawScore / session.MaxScore
	}

	session.Duration = "1s"

	return &session, nil
}

func (m *MockClient) CreateEvent(ctx context.Context, sessionID int, uuid string, eventType string, data string) (*platform.Event, error) {
	m.NumCalledCreateEvent++

	if sessionID <= 0 {
		return nil, InvalidSessionError
	}

	if eventType == "" {
		return nil, commonerrors.ErrorRequired("event type")
	}

	if m.CreateEventError != nil {
		return nil, m.CreateEventError
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

func (m *MockClient) GetPlatforms(ctx context.Context) ([]*Platform, error) {
	m.NumCalledGetPlatforms++

	if m.GetPlatformsError != nil {
		return nil, m.GetPlatformsError
	}

	return []*Platform{{
		ID:   1,
		Name: "android",
	}}, nil
}

func (m *MockClient) GetControlTypes(ctx context.Context) ([]*ControlType, error) {
	m.NumCalledGetControlTypes++

	if m.GetControlTypesError != nil {
		return nil, m.GetControlTypesError
	}

	return []*ControlType{{
		ID:   1,
		Name: "keyboard/mouse",
	}}, nil
}

func (m *MockClient) CreateModuleVersion(ctx context.Context, input ModuleVersion) (*ModuleVersion, error) {
	m.NumCalledCreateModuleVersion++

	if input.ModuleID == 0 {
		return nil, errors.New("invalid module id")
	}

	if input.LocalFilePath == "" {
		return nil, errors.New("local file path required")
	}

	if input.SemanticVersion == "" {
		return nil, errors.New("invalid version")
	}

	if m.CreateModuleVersionError != nil {
		return nil, m.CreateModuleVersionError
	}

	return &ModuleVersion{
		ModuleID:        input.ModuleID,
		SemanticVersion: input.SemanticVersion,
		Package:         input.Package,
	}, nil
}

func (m *MockClient) GetMultiplayerServerConfigs(ctx context.Context, params *MultiplayerServerConfigParams) ([]*MultiplayerServerConfigQueryParams, error) {
	m.NumCalledGetMultiplayerServerConfigs++

	if m.GetMultiplayerServerConfigsError != nil {
		return nil, m.GetMultiplayerServerConfigsError
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

func (m *MockClient) GetMultiplayerServerVersions(ctx context.Context, params *MultiplayerServerVersionQueryParams) ([]*MultiplayerServerVersion, error) {
	m.NumCalledGetMultiplayerServerVersions++

	if params.ModuleID == 0 {
		return nil, errors.New("invalid module id")
	}

	if params.SemanticVersion == "" {
		return nil, errors.New("invalid semantic version")
	}

	if m.GetMultiplayerServerVersionsError != nil {
		return nil, m.GetMultiplayerServerVersionsError
	}

	if m.GetMultiplayerServerVersionsEmpty {
		return []*MultiplayerServerVersion{}, nil
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

func (m *MockClient) GetMultiplayerServerVersion(ctx context.Context, versionID int) (*MultiplayerServerVersion, error) {
	m.NumCalledGetMultiplayerServerVersion++

	if m.GetMultiplayerServerVersionError != nil {
		return nil, m.GetMultiplayerServerVersionError
	}

	if m.GetMultiplayerServerVersionEmpty {
		return nil, nil
	}

	return &MultiplayerServerVersion{
		ID:              versionID,
		ModuleID:        1,
		SemanticVersion: "1.0.0",
		Status:          "enabled",
		Engine:          "unreal",
	}, nil
}

func (m *MockClient) CreateMultiplayerServerVersion(ctx context.Context, input MultiplayerServerVersion) (*MultiplayerServerVersion, error) {
	m.NumCalledCreateMultiplayerServerVersion++

	if input.ModuleID == 0 {
		return nil, errors.New("invalid module id")
	}

	if input.ImageRegistry == "" && input.LocalFilePath == "" {
		return nil, errors.New("image or file path required")
	}

	if input.SemanticVersion == "" {
		return nil, errors.New("invalid semantic version")
	}

	if m.CreateMultiplayerServerVersionError != nil {
		return nil, m.CreateMultiplayerServerVersionError
	}

	return &MultiplayerServerVersion{
		ModuleID:        input.ModuleID,
		SemanticVersion: input.SemanticVersion,
		Status:          input.Status,
		Engine:          input.Engine,
	}, nil
}
