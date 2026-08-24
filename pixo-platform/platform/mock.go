package platform

import (
	"context"
	"errors"
	"time"

	abstract "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract"
	commonerrors "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/commonerrors"
	faker "github.com/go-faker/faker/v4"
)

var _ Client = (*MockClient)(nil)

type MockClient struct {
	abstract.MockAbstractClient

	NumCalledCheckConnection int
	CheckConnectionError     error

	NumCalledGetUser int
	GetUserResponse  *User
	GetUserError     error

	NumCalledGetUserByUsername int
	GetUserByUsernameError     error

	NumCalledGetUsers int
	GetUsersError     error

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

	NumCalledGetOrgModules  int
	GetOrgModulesError      error
	GetOrgModulesParameters *OrgModuleParams
	GetOrgModulesReturn     []OrgModule

	NumCalledCreateOrgModule int
	CreateOrgModuleError     error

	NumCalledDeleteOrgModule int
	DeleteOrgModuleError     error

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

	NumCalledGetModules  int
	GetModulesEmpty      bool
	GetModulesError      error
	GetModulesParameters *ModuleParams
	GetModulesReturn     []Module

	NumCalledGetModulesWithAssociations  int
	GetModulesWithAssociationsEmpty      bool
	GetModulesWithAssociationsError      error
	GetModulesWithAssociationsParameters *ModuleParams
	GetModulesWithAssociationsReturn     []Module

	NumCalledGetModule int
	GetModuleError     error
	GetModuleIDParam   int
	GetModuleReturn    *Module

	NumCalledGetAPIKeys int
	GetAPIKeysEmpty     bool
	GetAPIKeysError     error

	NumCalledCreateAPIKey int
	CreateAPIKeyError     error

	NumCalledDeleteAPIKey int
	DeleteAPIKeyError     error

	NumCalledGetSession int
	GetSessionError     error

	CalledGetAssetWith []*Asset
	GetAssetReturns    *Asset
	GetAssetError      error

	CalledGetAssetsWith []*AssetParams
	GetAssetsReturns    []Asset
	GetAssetsError      error

	CalledCreateAssetWith []*Asset
	CreateAssetError      error

	CalledCreateAssetVersionWith []*AssetVersion
	CreateAssetVersionError      error

	CalledUpdateAssetVersionWith []AssetVersion
	UpdateAssetVersionError      error

	CalledPostAssetWith []AssetVersion
	PostAssetError      error

	CalledRetrieveAssetsWith []AssetParams
	RetrieveAssetsReturns    []Asset
	RetrieveAssetsError      error

	CalledCreateSessionWith []*Session
	CreateSessionError      error

	CalledUpdateSessionWith []*Session
	UpdateSessionError      error

	CalledCreateEventWith []*Event
	CreateEventError      error

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

	NumCalledGetMultiplayerServerVersionsWithConfig int
	GetMultiplayerServerVersionsWithConfigError     error
	GetMultiplayerServerVersionsWithConfigEmpty     bool

	NumCalledGetMultiplayerServerVersions int
	GetMultiplayerServerVersionsError     error
	GetMultiplayerServerVersionsEmpty     bool
	GetMultiplayerServerVersionsResponse  []MultiplayerServerVersion

	NumCalledGetMultiplayerServerVersion int
	GetMultiplayerServerVersionError     error
	GetMultiplayerServerVersionEmpty     bool

	NumCalledCreateMultiplayerServerVersion int
	CreateMultiplayerServerVersionError     error

	NumCalledUpdateMultiplayerServerVersion int
	UpdateMultiplayerServerVersionError     error

	NumCalledGetLearningHistoryRecords int
	GetLearningHistoryRecordsError     error
	LearningHistoryRecordsToReturn     []LearningHistory

	NumCalledGetCourseDataRecords int
	GetCourseDataRecordsError     error
	CourseDataRecordsToReturn     []CourseData

	NumCalledGetOrgSuccessFactors int
	GetOrgSuccessFactorsError     error
	OrgSuccessFactorsToReturn     []OrgSuccessFactor

	NumCalledGetExpiringModules int
	GetExpiringModulesError     error
	ExpiringModulesToReturn     []ExpiringModule

	NumCalledGetModuleVersions  int
	GetModuleVersionsError      error
	GetModuleVersionsParameters *ModuleVersionParams
	ModuleVersionsReturns       []ModuleVersion

	NumCalledGetVersionLifecycles int
	GetVersionLifecyclesError     error
	VersionLifecyclesReturns      []VersionLifecycle

	NumCalledGetModulesForUser int
	GetModulesForUserError     error
	GetModulesForUserIDParam   int
	GetModulesForUserReturn    []Module

	NumCalledGetUsersWithModuleAccess int
	GetUsersWithModuleAccessError     error
	GetUsersWithModuleAccessModuleID  int
	GetUsersWithModuleAccessOrgID     int
	GetUsersWithModuleAccessReturn    []User

	NumCalledGetModulePlayers  int
	GetModulePlayersError      error
	GetModulePlayersParameters *ModulePlayerParams
	GetModulePlayersReturn     []ModulePlayer

	NumCalledGetModulePlayerVersions  int
	GetModulePlayerVersionsError      error
	GetModulePlayerVersionsParameters *ModulePlayerVersionParams
	GetModulePlayerVersionsReturn     []ModulePlayerVersion
}

func (m *MockClient) Reset() {
	m.CalledCreateEventWith = nil
	m.CreateEventError = nil

	m.NumCalledCheckConnection = 0
	m.CheckConnectionError = nil

	m.CreateSessionError = nil

	m.NumCalledGetUser = 0
	m.GetUserError = nil

	m.NumCalledGetUserByUsername = 0
	m.GetUserByUsernameError = nil

	m.NumCalledGetUsers = 0
	m.GetUsersError = nil

	m.NumCalledCreateUser = 0
	m.CreateUserError = nil

	m.NumCalledUpdateUser = 0
	m.UpdateUserError = nil

	m.NumCalledDeleteUser = 0
	m.DeleteUserError = nil

	m.NumCalledGetRoles = 0
	m.GetRolesError = nil

	m.NumCalledGetOrgs = 0
	m.GetOrgsError = nil

	m.NumCalledGetOrg = 0
	m.GetOrgError = nil

	m.NumCalledCreateOrg = 0
	m.CreateOrgError = nil

	m.NumCalledUpdateOrg = 0
	m.UpdateOrgError = nil

	m.NumCalledDeleteOrg = 0
	m.DeleteOrgError = nil
	m.NumCalledGetOrgModules = 0
	m.GetOrgModulesError = nil
	m.GetOrgModulesParameters = nil
	m.GetOrgModulesReturn = nil

	m.NumCalledCreateOrgModule = 0
	m.CreateOrgModuleError = nil
	m.NumCalledDeleteOrgModule = 0
	m.DeleteOrgModuleError = nil

	m.NumCalledGetAPIKeys = 0
	m.GetAPIKeysError = nil

	m.NumCalledCreateAPIKey = 0
	m.CreateAPIKeyError = nil

	m.NumCalledDeleteAPIKey = 0
	m.DeleteAPIKeyError = nil

	m.NumCalledGetWebhooks = 0
	m.GetWebhooksError = nil

	m.NumCalledGetWebhook = 0
	m.GetWebhookError = nil

	m.NumCalledCreateWebhook = 0
	m.CreateWebhookError = nil

	m.NumCalledUpdateWebhook = 0
	m.UpdateWebhookError = nil

	m.NumCalledDeleteWebhook = 0
	m.DeleteWebhookError = nil

	m.NumCalledGetModules = 0
	m.GetModulesError = nil
	m.GetModulesEmpty = false
	m.GetModulesParameters = nil
	m.GetModulesReturn = nil

	m.NumCalledGetModulesWithAssociations = 0
	m.GetModulesWithAssociationsError = nil
	m.GetModulesWithAssociationsEmpty = false
	m.GetModulesWithAssociationsParameters = nil
	m.GetModulesWithAssociationsReturn = nil

	m.NumCalledGetModule = 0
	m.GetModuleError = nil
	m.GetModuleIDParam = 0
	m.GetModuleReturn = nil

	m.CalledGetAssetWith = nil
	m.GetAssetReturns = nil
	m.GetAssetError = nil

	m.CalledGetAssetsWith = nil
	m.GetAssetsReturns = nil
	m.GetAssetsError = nil

	m.CalledCreateAssetWith = nil
	m.CreateAssetError = nil

	m.CalledCreateAssetWith = nil
	m.CreateAssetError = nil

	m.CalledUpdateAssetVersionWith = nil
	m.UpdateAssetVersionError = nil

	m.NumCalledGetSession = 0
	m.GetSessionError = nil

	m.CalledCreateSessionWith = nil
	m.CreateSessionError = nil

	m.CalledUpdateSessionWith = nil
	m.UpdateSessionError = nil

	m.NumCalledGetPlatforms = 0
	m.GetPlatformsError = nil

	m.NumCalledGetControlTypes = 0
	m.GetControlTypesError = nil

	m.NumCalledCreateModuleVersion = 0
	m.CreateModuleVersionError = nil

	m.NumCalledGetVersionLifecycles = 0
	m.GetVersionLifecyclesError = nil
	m.VersionLifecyclesReturns = nil

	m.NumCalledGetMultiplayerServerConfigs = 0
	m.GetMultiplayerServerConfigsError = nil

	m.NumCalledGetMultiplayerServerVersions = 0
	m.GetMultiplayerServerVersionsError = nil
	m.GetMultiplayerServerVersionsResponse = nil

	m.NumCalledGetMultiplayerServerVersionsWithConfig = 0
	m.GetMultiplayerServerVersionsWithConfigError = nil

	m.NumCalledGetMultiplayerServerVersion = 0
	m.GetMultiplayerServerVersionError = nil

	m.NumCalledCreateMultiplayerServerVersion = 0
	m.CreateMultiplayerServerVersionError = nil

	m.NumCalledUpdateMultiplayerServerVersion = 0
	m.UpdateMultiplayerServerVersionError = nil

	m.NumCalledGetLearningHistoryRecords = 0
	m.GetLearningHistoryRecordsError = nil
	m.LearningHistoryRecordsToReturn = nil

	m.NumCalledGetCourseDataRecords = 0
	m.GetCourseDataRecordsError = nil
	m.CourseDataRecordsToReturn = nil

	m.NumCalledGetOrgSuccessFactors = 0
	m.GetOrgSuccessFactorsError = nil
	m.OrgSuccessFactorsToReturn = nil

	m.NumCalledGetExpiringModules = 0
	m.GetExpiringModulesError = nil
	m.ExpiringModulesToReturn = nil

	m.NumCalledGetModuleVersions = 0
	m.GetModuleVersionsError = nil
	m.GetModuleVersionsParameters = nil
	m.ModuleVersionsReturns = nil

	m.NumCalledGetModulesForUser = 0
	m.GetModulesForUserError = nil
	m.GetModulesForUserIDParam = 0
	m.GetModulesForUserReturn = nil

	m.NumCalledGetUsersWithModuleAccess = 0
	m.GetUsersWithModuleAccessError = nil
	m.GetUsersWithModuleAccessModuleID = 0
	m.GetUsersWithModuleAccessOrgID = 0
	m.GetUsersWithModuleAccessReturn = nil

	m.NumCalledGetModulePlayers = 0
	m.GetModulePlayersError = nil
	m.GetModulePlayersParameters = nil
	m.GetModulePlayersReturn = nil

	m.NumCalledGetModulePlayerVersions = 0
	m.GetModulePlayerVersionsError = nil
	m.GetModulePlayerVersionsParameters = nil
	m.GetModulePlayerVersionsReturn = nil
}

func (m *MockClient) CheckAuth(ctx context.Context) (User, error) {
	if !m.IsAuthenticated() {
		return User{}, errors.New("not authenticated")
	}

	return User{ID: m.ActiveUserID()}, nil
}

func (m *MockClient) CheckConnection(ctx context.Context) error {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledCheckConnection++

	if m.CheckConnectionError != nil {
		return m.CheckConnectionError
	}

	// IsAuthenticated is hardcoded to false on the embedded mock, so read the credentials the
	// caller configured instead - otherwise the mock can never report a working connection.
	if m.APIKey == "" && m.Token == "" {
		return ErrNoCredentials
	}

	return nil
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

func (m *MockClient) GetUser(ctx context.Context, id int) (*User, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledGetUser++

	if id <= 0 {
		return nil, errors.New("invalid user id")
	}

	if m.GetUserError != nil {
		return nil, m.GetUserError
	}

	if m.GetUserResponse != nil {
		return m.GetUserResponse, nil
	}

	return &User{
		ID:        id,
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		Username:  faker.Username(),
		OrgID:     1,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}

func (m *MockClient) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledGetUserByUsername++

	if username == "" {
		return nil, commonerrors.ErrorRequired("username")
	}

	if m.GetUserByUsernameError != nil {
		return nil, m.GetUserByUsernameError
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
	m.Lock.Lock()
	defer m.Lock.Unlock()

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
	m.Lock.Lock()
	defer m.Lock.Unlock()

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
	m.Lock.Lock()
	defer m.Lock.Unlock()

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
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledGetModules++

	if len(params) > 0 {
		m.GetModulesParameters = &params[0]
	}

	if m.GetModulesError != nil {
		return nil, m.GetModulesError
	}

	if m.GetModulesEmpty {
		return []Module{}, nil
	}

	if m.GetModulesReturn != nil {
		return m.GetModulesReturn, nil
	}

	return []Module{
		{
			ID:           1,
			Abbreviation: "TST",
		},
		{
			ID:           2,
			Abbreviation: "TST-2",
		},
	}, nil
}

func (m *MockClient) GetModulesWithAssociations(ctx context.Context, params ModuleParams) ([]Module, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledGetModulesWithAssociations++
	m.GetModulesWithAssociationsParameters = &params

	if m.GetModulesWithAssociationsError != nil {
		return nil, m.GetModulesWithAssociationsError
	}

	if m.GetModulesWithAssociationsEmpty {
		return []Module{}, nil
	}

	if m.GetModulesWithAssociationsReturn != nil {
		return m.GetModulesWithAssociationsReturn, nil
	}

	return []Module{
		{
			ID:           1,
			Abbreviation: "TST",
			Versions: []ModuleVersion{
				{
					ID:              1,
					ModuleID:        1,
					SemanticVersion: "1.0.0",
					LifecycleID:     1,
					FilePath:        "ModuleVersions/1/zips/module.zip",
					Platforms:       []Platform{{ID: 1, ShortName: "quest"}},
				},
			},
		},
	}, nil
}

func (m *MockClient) GetModule(ctx context.Context, id int) (*Module, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledGetModule++
	m.GetModuleIDParam = id

	if id == 0 {
		return nil, errors.New("module id is required")
	}

	if m.GetModuleError != nil {
		return nil, m.GetModuleError
	}

	if m.GetModuleReturn != nil {
		return m.GetModuleReturn, nil
	}

	return &Module{
		ID:           id,
		Abbreviation: "TST",
		Versions: []ModuleVersion{
			{
				ID:              1,
				ModuleID:        id,
				SemanticVersion: "1.0.0",
				LifecycleID:     1,
				FilePath:        "ModuleVersions/1/zips/module.zip",
				Platforms:       []Platform{{ID: 1, ShortName: "quest"}},
			},
		},
	}, nil
}

func (m *MockClient) GetRoles(ctx context.Context) ([]Role, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

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
	m.Lock.Lock()
	defer m.Lock.Unlock()

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
	m.Lock.Lock()
	defer m.Lock.Unlock()

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
	m.Lock.Lock()
	defer m.Lock.Unlock()

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
	m.Lock.Lock()
	defer m.Lock.Unlock()

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
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledDeleteOrg++

	if id <= 0 {
		return errors.New("invalid org id")
	}

	if m.DeleteOrgError != nil {
		return m.DeleteOrgError
	}

	return nil
}

func (m *MockClient) GetOrgModules(ctx context.Context, params OrgModuleParams) ([]OrgModule, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledGetOrgModules++
	m.GetOrgModulesParameters = &params

	if params.OrgID == 0 {
		return nil, errors.New("org id is required")
	}

	if m.GetOrgModulesError != nil {
		return nil, m.GetOrgModulesError
	}

	if m.GetOrgModulesReturn != nil {
		return m.GetOrgModulesReturn, nil
	}

	return []OrgModule{
		{
			ID:       1,
			OrgID:    params.OrgID,
			Lifetime: true,
			Module:   &Module{ID: 1, Abbreviation: "TST"},
		},
	}, nil
}

func (m *MockClient) CreateOrgModule(ctx context.Context, input OrgModule) (*OrgModule, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledCreateOrgModule++

	if input.OrgID == 0 || input.ModuleID == 0 {
		return nil, errors.New("org id and module id are required")
	}

	if m.CreateOrgModuleError != nil {
		return nil, m.CreateOrgModuleError
	}

	orgModule := input
	orgModule.ID = 1

	if input.ExpirDate != nil {
		orgModule.ExpiresAt = input.ExpirDate.Format(time.RFC3339)
	}

	return &orgModule, nil
}

func (m *MockClient) DeleteOrgModule(ctx context.Context, orgID, moduleID int) error {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledDeleteOrgModule++

	if orgID == 0 || moduleID == 0 {
		return errors.New("org id and module id are required")
	}

	return m.DeleteOrgModuleError
}

func (m *MockClient) CreateAPIKey(ctx context.Context, input APIKey) (*APIKey, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

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
	m.Lock.Lock()
	defer m.Lock.Unlock()

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
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledDeleteAPIKey++

	if id <= 0 {
		return errors.New("invalid api key id")
	}

	if m.DeleteAPIKeyError != nil {
		return m.DeleteAPIKeyError
	}

	return nil
}

func (m *MockClient) GetAsset(ctx context.Context, id int) (*Asset, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.CalledGetAssetWith = append(m.CalledGetAssetWith, &Asset{ID: id})

	if m.GetAssetError != nil {
		return nil, m.GetAssetError
	}

	if m.GetAssetReturns != nil {
		return m.GetAssetReturns, nil
	}

	return &Asset{
		ID:        id,
		ModuleID:  1,
		Name:      faker.Name(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}

func (m *MockClient) GetAssets(ctx context.Context, params AssetParams) ([]Asset, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.CalledGetAssetsWith = append(m.CalledGetAssetsWith, &params)

	if m.GetAssetsError != nil {
		return nil, m.GetAssetsError
	}

	assets := []Asset{
		{
			ID:       1,
			ModuleID: 1,
			Name:     faker.Name(),
			Type:     "text",
			Versions: []AssetVersion{
				{
					ID:           1,
					AssetID:      1,
					Status:       "stage",
					LanguageCode: "en",
					CreatedAt:    time.Now().UTC(),
					UpdatedAt:    time.Now().UTC(),
				},
			},
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
	}

	if m.GetAssetsReturns != nil {
		assets = m.GetAssetsReturns
	}

	var filteredAssets []Asset
	for i := range assets {
		if params.ModuleID > 0 && assets[i].ModuleID != params.ModuleID {
			continue
		}

		if params.Name != "" && assets[i].Name != params.Name {
			continue
		}

		if params.LanguageCode != "" {
			var versions []AssetVersion
			for j := range assets[i].Versions {
				if assets[i].Versions[j].Status == params.Status &&
					assets[i].Versions[j].LanguageCode == params.LanguageCode {

					versions = append(versions, assets[i].Versions[j])
				}
			}

			assets[i].Versions = versions
		}

		filteredAssets = append(filteredAssets, assets[i])
	}

	return filteredAssets, nil
}

func (m *MockClient) CreateAsset(ctx context.Context, asset *Asset) error {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.CalledCreateAssetWith = append(m.CalledCreateAssetWith, asset)
	return m.CreateAssetError
}

func (m *MockClient) CreateAssetVersion(ctx context.Context, assetVersion *AssetVersion) error {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.CalledCreateAssetVersionWith = append(m.CalledCreateAssetVersionWith, assetVersion)

	assetVersion.ID = len(m.CalledCreateAssetVersionWith)

	return m.CreateAssetVersionError
}

func (m *MockClient) PostAsset(ctx context.Context, assetVersion *AssetVersion) error {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.CalledPostAssetWith = append(m.CalledPostAssetWith, *assetVersion)

	if assetVersion.LanguageCode == "" {
		assetVersion.LanguageCode = "en"
	}

	return m.CreateAssetVersionError
}

func (m *MockClient) RetrieveAssets(ctx context.Context, params AssetParams) ([]Asset, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.CalledRetrieveAssetsWith = append(m.CalledRetrieveAssetsWith, params)

	if m.RetrieveAssetsError != nil {
		return nil, m.RetrieveAssetsError
	}

	if m.RetrieveAssetsReturns != nil {
		return m.RetrieveAssetsReturns, nil
	}

	return []Asset{
		{
			ID:       1,
			ModuleID: 1,
			Name:     faker.Name(),
			Type:     "text",
			Versions: []AssetVersion{
				{
					ID:           1,
					AssetID:      1,
					Status:       "stage",
					LanguageCode: "en",
					CreatedAt:    time.Now().UTC(),
					UpdatedAt:    time.Now().UTC(),
				},
			},
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
	}, nil
}

func (m *MockClient) UpdateAssetVersion(ctx context.Context, assetVersion *AssetVersion) error {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.CalledUpdateAssetVersionWith = append(m.CalledUpdateAssetVersionWith, *assetVersion)

	if assetVersion.LanguageCode == "" {
		assetVersion.LanguageCode = "en"
	}

	return m.UpdateAssetVersionError
}

func (m *MockClient) GetWebhooks(ctx context.Context, params *WebhookParams) ([]Webhook, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

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
	m.Lock.Lock()
	defer m.Lock.Unlock()

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
	m.Lock.Lock()
	defer m.Lock.Unlock()

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
	m.Lock.Lock()
	defer m.Lock.Unlock()

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
	m.Lock.Lock()
	defer m.Lock.Unlock()

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
	m.Lock.Lock()
	defer m.Lock.Unlock()

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
	m.Lock.Lock()
	defer m.Lock.Unlock()

	if session == nil {
		return errors.New("session can not be nil")
	}

	sessionCopy := *session
	m.CalledCreateSessionWith = append(m.CalledCreateSessionWith, &sessionCopy)

	if session.ModuleID <= 0 {
		return errors.New("invalid module id")
	}

	if m.CreateSessionError != nil {
		return m.CreateSessionError
	}

	session.ID = 1
	session.UserID = m.ActiveUserID()
	session.Module.Abbreviation = "TST"
	session.IPAddress = "127.0.0.1"

	session.CreatedAt = time.Now().UTC()
	session.UpdatedAt = time.Now().UTC()

	return nil
}

func (m *MockClient) UpdateSession(ctx context.Context, session Session) (*Session, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.CalledUpdateSessionWith = append(m.CalledUpdateSessionWith, &session)

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
	m.Lock.Lock()
	defer m.Lock.Unlock()

	eventCopy := *event
	m.CalledCreateEventWith = append(m.CalledCreateEventWith, &eventCopy)

	noSessionID := event.SessionID == nil || *event.SessionID <= 0
	noSessionUUID := event.SessionUUID == nil || *event.SessionUUID == ""
	if noSessionID && noSessionUUID {
		return errors.New("session id or session uuid required")
	}

	if m.CreateEventError != nil {
		return m.CreateEventError
	}

	event.CreatedAt = time.Now().UTC()
	return nil
}

func (m *MockClient) GetPlatforms(ctx context.Context) ([]Platform, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledGetPlatforms++

	if m.GetPlatformsError != nil {
		return nil, m.GetPlatformsError
	}

	return []Platform{{
		ID:   1,
		Name: "android",
	}}, nil
}

func (m *MockClient) GetControlTypes(ctx context.Context) ([]ControlType, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledGetControlTypes++

	if m.GetControlTypesError != nil {
		return nil, m.GetControlTypesError
	}

	return []ControlType{{
		ID:   1,
		Name: "keyboard/mouse",
	}}, nil
}

func (m *MockClient) CreateModuleVersion(ctx context.Context, input ModuleVersion) (*ModuleVersion, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

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

func (m *MockClient) GetMultiplayerServerConfigs(ctx context.Context, params *MultiplayerServerConfigParams) ([]MultiplayerServerConfigQueryParams, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledGetMultiplayerServerConfigs++

	if m.GetMultiplayerServerConfigsError != nil {
		return nil, m.GetMultiplayerServerConfigsError
	}

	if m.GetMultiplayerServerConfigsEmpty {
		return []MultiplayerServerConfigQueryParams{}, nil
	}

	if m.GetMultiplayerServerConfigsEmptyVersions {
		return []MultiplayerServerConfigQueryParams{
			{
				ModuleID: 1,
				Capacity: 5,
			},
		}, nil
	}

	return []MultiplayerServerConfigQueryParams{
		{
			ModuleID: 1,
			Capacity: 5,
			ServerVersions: []MultiplayerServerVersion{
				{
					Engine:          "unreal",
					ImageRegistry:   "gcr.io/pixo-bootstrap/multiplayer/gameservers/simple-server:latest",
					Status:          "enabled",
					SemanticVersion: "1.0.0",
				},
			},
		},
	}, nil
}

func (m *MockClient) GetMultiplayerServerVersions(ctx context.Context, params *MultiplayerServerVersionParams) ([]MultiplayerServerVersion, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledGetMultiplayerServerVersions++

	if m.GetMultiplayerServerVersionsError != nil {
		return nil, m.GetMultiplayerServerVersionsError
	}

	if m.GetMultiplayerServerVersionsEmpty {
		return []MultiplayerServerVersion{}, nil
	}

	if m.GetMultiplayerServerVersionsResponse != nil {
		return m.GetMultiplayerServerVersionsResponse, nil
	}

	return []MultiplayerServerVersion{
		{
			ModuleID:        1,
			SemanticVersion: "1.0.0",
			Status:          "enabled",
			Engine:          "unreal",
			FilePath:        "module-1/version-1/file.zip",
		},
	}, nil
}

func (m *MockClient) GetMultiplayerServerVersionsWithConfig(ctx context.Context, params *MultiplayerServerVersionParams) ([]MultiplayerServerVersion, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledGetMultiplayerServerVersionsWithConfig++

	if m.GetMultiplayerServerVersionsWithConfigError != nil {
		return nil, m.GetMultiplayerServerVersionsWithConfigError
	}

	if m.GetMultiplayerServerVersionsWithConfigEmpty {
		return []MultiplayerServerVersion{}, nil
	}

	return []MultiplayerServerVersion{
		{
			ModuleID:        1,
			SemanticVersion: "1.0.0",
			Status:          "enabled",
			Engine:          "unreal",
		},
	}, nil
}

func (m *MockClient) GetMultiplayerServerVersion(ctx context.Context, versionID int) (*MultiplayerServerVersion, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

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
	m.Lock.Lock()
	defer m.Lock.Unlock()

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
		Module:          &Module{ID: input.ModuleID, Abbreviation: "TST"},
		SemanticVersion: input.SemanticVersion,
		Status:          input.Status,
		Engine:          input.Engine,
	}, nil
}

func (m *MockClient) UpdateMultiplayerServerVersion(ctx context.Context, input MultiplayerServerVersion) (*MultiplayerServerVersion, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledUpdateMultiplayerServerVersion++

	if input.ID == 0 && input.ModuleID == 0 && input.SemanticVersion == "" {
		return nil, errors.New("id, or module id and semantic version is required")
	}

	if m.UpdateMultiplayerServerVersionError != nil {
		return nil, m.UpdateMultiplayerServerVersionError
	}

	return &MultiplayerServerVersion{
		ID:              input.ID,
		ModuleID:        input.ModuleID,
		SemanticVersion: input.SemanticVersion,
		ImageRegistry:   input.ImageRegistry,
		Status:          input.Status,
		Engine:          input.Engine,
		Module: &Module{
			Abbreviation: "TST",
		},
	}, nil
}

func (m *MockClient) GetLearningHistoryRecords(ctx context.Context, params LearningHistoryParams) ([]LearningHistory, error) {
	m.NumCalledGetLearningHistoryRecords++

	if m.GetLearningHistoryRecordsError != nil {
		return nil, m.GetLearningHistoryRecordsError
	}

	if len(m.LearningHistoryRecordsToReturn) > 0 {
		return m.LearningHistoryRecordsToReturn, nil
	}

	return []LearningHistory{}, nil
}

func (m *MockClient) GetCourseDataRecords(ctx context.Context, orgID int) ([]CourseData, error) {
	m.NumCalledGetCourseDataRecords++

	if m.GetCourseDataRecordsError != nil {
		return nil, m.GetCourseDataRecordsError
	}

	if len(m.CourseDataRecordsToReturn) > 0 {
		return m.CourseDataRecordsToReturn, nil
	}

	return []CourseData{}, nil
}

func (m *MockClient) GetOrgSuccessFactors(ctx context.Context) ([]OrgSuccessFactor, error) {
	m.NumCalledGetOrgSuccessFactors++

	if m.GetOrgSuccessFactorsError != nil {
		return nil, m.GetOrgSuccessFactorsError
	}

	if len(m.OrgSuccessFactorsToReturn) > 0 {
		return m.OrgSuccessFactorsToReturn, nil
	}

	return []OrgSuccessFactor{}, nil
}

func (m *MockClient) GetExpiringModules(ctx context.Context) ([]ExpiringModule, error) {
	m.NumCalledGetExpiringModules++

	if m.GetExpiringModulesError != nil {
		return nil, m.GetExpiringModulesError
	}

	if len(m.ExpiringModulesToReturn) > 0 {
		return m.ExpiringModulesToReturn, nil
	}

	return []ExpiringModule{}, nil
}

func (m *MockClient) GetModuleVersions(ctx context.Context, params *ModuleVersionParams) ([]ModuleVersion, error) {
	m.NumCalledGetModuleVersions++
	m.GetModuleVersionsParameters = params

	if m.GetModuleVersionsError != nil {
		return nil, m.GetModuleVersionsError
	}

	return m.ModuleVersionsReturns, nil
}

func (m *MockClient) GetVersionLifecycles(ctx context.Context) ([]VersionLifecycle, error) {
	m.NumCalledGetVersionLifecycles++

	if m.GetVersionLifecyclesError != nil {
		return nil, m.GetVersionLifecyclesError
	}

	if m.VersionLifecyclesReturns != nil {
		return m.VersionLifecyclesReturns, nil
	}

	return []VersionLifecycle{
		{ID: 1, Name: LifecycleDevelopment},
		{ID: 2, Name: LifecycleQA},
		{ID: 3, Name: LifecycleReleased},
		{ID: 4, Name: LifecycleArchived},
	}, nil
}

func (m *MockClient) GetModulesForUser(ctx context.Context, userID int) ([]Module, error) {
	m.NumCalledGetModulesForUser++
	m.GetModulesForUserIDParam = userID
	if m.GetModulesForUserError != nil {
		return nil, m.GetModulesForUserError
	}
	return m.GetModulesForUserReturn, nil
}

func (m *MockClient) GetUsersWithModuleAccess(ctx context.Context, moduleID, orgID int) ([]User, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledGetUsersWithModuleAccess++
	m.GetUsersWithModuleAccessModuleID = moduleID
	m.GetUsersWithModuleAccessOrgID = orgID
	if m.GetUsersWithModuleAccessError != nil {
		return nil, m.GetUsersWithModuleAccessError
	}

	if m.GetUsersWithModuleAccessReturn != nil {
		return m.GetUsersWithModuleAccessReturn, nil
	}

	return []User{{ID: 1, Username: "test-user", Role: "user", OrgID: orgID}}, nil
}

func (m *MockClient) GetModulePlayers(ctx context.Context, params ...*ModulePlayerParams) ([]ModulePlayer, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledGetModulePlayers++
	if len(params) > 0 {
		m.GetModulePlayersParameters = params[0]
	}

	if m.GetModulePlayersError != nil {
		return nil, m.GetModulePlayersError
	}

	if m.GetModulePlayersReturn != nil {
		return m.GetModulePlayersReturn, nil
	}

	return []ModulePlayer{
		{
			ID:             1,
			Name:           "test-module-player",
			DistributorID:  1,
			LaunchProtocol: "pixo",
			Distributor:    Org{ID: 1, Name: "test-org"},
		},
	}, nil
}

func (m *MockClient) GetModulePlayerVersions(ctx context.Context, params *ModulePlayerVersionParams) ([]ModulePlayerVersion, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledGetModulePlayerVersions++
	m.GetModulePlayerVersionsParameters = params

	if m.GetModulePlayerVersionsError != nil {
		return nil, m.GetModulePlayerVersionsError
	}

	if m.GetModulePlayerVersionsReturn != nil {
		return m.GetModulePlayerVersionsReturn, nil
	}

	return []ModulePlayerVersion{
		{
			ID:             1,
			ModulePlayerID: 1,
			Status:         "enabled",
			Version:        "1.0.0",
			Package:        "com.pixovr.test",
		},
	}, nil
}
