package platform

import (
	"context"
	"errors"
	abstract "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	commonerrors "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/commonerrors"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/agones"
	"github.com/go-faker/faker/v4"
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

	NumCalledGetRoles int
	GetRolesError     error

	NumCalledGetOrgs int
	GetOrgsError     error

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

	NumCalledGetModules int
	GetModulesEmpty     bool
	GetModulesError     error

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

func (m *MockClient) Reset() {
	m.NumCalledCreateEvent = 0

	m.NumCalledCreateSession = 0

	m.NumCalledGetUser = 0
	m.NumCalledCreateUser = 0
	m.NumCalledUpdateOrg = 0
	m.NumCalledDeleteUser = 0

	m.NumCalledGetRoles = 0

	m.NumCalledGetOrgs = 0
	m.NumCalledGetOrg = 0
	m.NumCalledCreateOrg = 0
	m.NumCalledUpdateUser = 0
	m.NumCalledDeleteOrg = 0

	m.NumCalledGetAPIKeys = 0
	m.NumCalledCreateAPIKey = 0
	m.NumCalledDeleteAPIKey = 0

	m.NumCalledGetWebhooks = 0
	m.NumCalledGetWebhook = 0
	m.NumCalledCreateWebhook = 0
	m.NumCalledUpdateWebhook = 0
	m.NumCalledDeleteWebhook = 0

	m.NumCalledGetSession = 0
	m.NumCalledUpdateSession = 0

	m.NumCalledGetPlatforms = 0

	m.NumCalledGetControlTypes = 0

	m.NumCalledCreateModuleVersion = 0

	m.NumCalledGetMultiplayerServerConfigs = 0
	m.NumCalledGetMultiplayerServerVersions = 0
	m.NumCalledGetMultiplayerServerVersion = 0
	m.NumCalledCreateMultiplayerServerVersion = 0
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

func (m *MockClient) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	m.NumCalledGetUser++

	if username == "" {
		return nil, commonerrors.ErrorRequired("username")
	}

	if m.GetUserError != nil {
		return nil, m.GetUserError
	}

	return &User{
		ID:        1,
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		Username:  username,
		OrgID:     1,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}

func (m *MockClient) CreateUser(ctx context.Context, user *User) error {
	m.NumCalledCreateUser++

	if user == nil {
		return errors.New("user is nil")
	}

	if user.Username == "" && user.Email == "" {
		return commonerrors.ErrorRequired("username or email")
	}

	if user.Password == "" {
		return commonerrors.ErrorRequired("password")
	}

	if user.OrgID <= 0 {
		return errors.New("invalid org id")
	}

	if m.CreateUserError != nil {
		return m.CreateUserError
	}

	user.ID = 1
	user.Org = Org{ID: user.OrgID, Name: "test-org"}
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()
	return nil
}

func (m *MockClient) UpdateUser(ctx context.Context, user *User) error {
	m.NumCalledUpdateUser++

	if user == nil {
		return errors.New("user is nil")
	}

	if user.ID <= 0 {
		return errors.New("invalid user id")
	}

	if user.Username == "" {
		return commonerrors.ErrorRequired("username")
	}

	if user.OrgID <= 0 {
		return errors.New("invalid org id")
	}

	if m.UpdateUserError != nil {
		return m.UpdateUserError
	}

	user.UpdatedAt = time.Now().UTC()
	return nil
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

func (m *MockClient) GetModules(ctx context.Context, params ...ModuleParams) ([]Module, error) {
	m.NumCalledGetModules++

	if m.GetModulesError != nil {
		return nil, m.GetModulesError
	}

	if m.GetModulesEmpty {
		return []Module{}, nil
	}

	return []Module{
		{
			ID:           1,
			Name:         "test",
			Abbreviation: "TST",
		},
		{
			ID:           2,
			Name:         "test-2",
			Abbreviation: "TST-2",
		},
	}, nil
}

func (m *MockClient) GetRoles(ctx context.Context) ([]Role, error) {
	m.NumCalledGetRoles++

	if m.GetRolesError != nil {
		return nil, m.GetRolesError
	}

	return []Role{
		{
			ID:   1,
			Name: "admin",
		},
		{
			ID:   2,
			Name: "user",
		},
	}, nil
}

func (m *MockClient) GetOrgs(ctx context.Context, params ...*OrgParams) ([]Org, error) {
	m.NumCalledGetOrgs++

	if m.GetOrgsError != nil {
		return nil, m.GetOrgsError
	}

	orgs := []Org{
		{
			ID:   1,
			Name: "test-org",
		},
		{
			ID:   2,
			Name: "test-org-2",
		},
	}

	if len(params) > 0 && params[0] != nil {
		orgs[0].Name = params[0].Name
	}

	return orgs, nil
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

func (m *MockClient) CreateAPIKey(ctx context.Context, input APIKey) (*APIKey, error) {
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

func (m *MockClient) GetAPIKeys(ctx context.Context, params *APIKeyQueryParams) ([]APIKey, error) {
	m.NumCalledGetAPIKeys++

	if m.GetAPIKeysError != nil {
		return nil, m.GetAPIKeysError
	}

	if m.GetAPIKeysEmpty {
		return []APIKey{}, nil
	}

	if params == nil {
		params = &APIKeyQueryParams{}
	}

	if params.UserID == nil {
		params.UserID = &[]int{1}[0]
	}

	return []APIKey{
		{
			ID:        1,
			UserID:    *params.UserID,
			Key:       faker.UUIDHyphenated(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:        2,
			UserID:    *params.UserID,
			Key:       faker.UUIDHyphenated(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:        3,
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
			ID:          1,
			OrgID:       1,
			URL:         "https://example.com",
			Description: "test",
			Token:       "token",
		},
		{
			ID:          2,
			OrgID:       2,
			URL:         "https://example-2.com",
			Description: "test-2",
			Token:       "token-2",
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

	if input.GenerateToken != nil && *input.GenerateToken {
		input.Token = faker.UUIDHyphenated()
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

func (m *MockClient) CreateSession(ctx context.Context, session *Session) error {
	m.NumCalledCreateSession++

	if session == nil {
		return errors.New("session details are required")
	}

	if session.ModuleID <= 0 {
		return errors.New("invalid module id")
	}

	if session.IPAddress == "" {
		return commonerrors.ErrorRequired("ip address")
	}

	if m.CreateSessionError != nil {
		return m.CreateSessionError
	}

	*session = Session{
		ID:        1,
		UserID:    m.ActiveUserID(),
		ModuleID:  session.ModuleID,
		Module:    Module{ID: session.ModuleID, Name: "test", Abbreviation: "TST"},
		IPAddress: session.IPAddress,
		DeviceID:  session.DeviceID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	return nil
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

func (m *MockClient) CreateEvent(ctx context.Context, event *Event) error {
	m.NumCalledCreateEvent++

	noSessionID := event.SessionID == nil || *event.SessionID <= 0
	noSessionUUID := event.SessionUUID == nil || *event.SessionUUID == ""
	if noSessionID && noSessionUUID {
		return errors.New("session id or session uuid required")
	}

	if event.Type == "" {
		return commonerrors.ErrorRequired("type")
	}

	if m.CreateEventError != nil {
		return m.CreateEventError
	}

	event.CreatedAt = time.Now().UTC()
	return nil
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
