package platform

import (
	"context"
	abstract "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract"
)

type Client interface {
	abstract.AbstractClient

	CheckAuth(ctx context.Context) (User, error)
	ActiveUserID() int
	ActiveOrgID() int

	GetUser(ctx context.Context, id int) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int) error

	GetAPIKeys(ctx context.Context, params *APIKeyQueryParams) ([]APIKey, error)
	CreateAPIKey(ctx context.Context, input APIKey) (*APIKey, error)
	DeleteAPIKey(ctx context.Context, id int) error

	GetWebhooks(ctx context.Context, params *WebhookParams) ([]Webhook, error)
	GetWebhook(ctx context.Context, id int) (*Webhook, error)
	CreateWebhook(ctx context.Context, input Webhook) (*Webhook, error)
	UpdateWebhook(ctx context.Context, input Webhook) (*Webhook, error)
	DeleteWebhook(ctx context.Context, id int) error

	GetRoles(ctx context.Context) ([]Role, error)

	// GetPlatforms retrieves platforms from the platform using the GraphQL interface
	GetPlatforms(ctx context.Context) ([]Platform, error)
	// GetControlTypes retrieves control types from the platform using the GraphQL interface
	GetControlTypes(ctx context.Context) ([]ControlType, error)

	// GetModules retrieves modules from the platform using the GraphQL interface
	GetModules(ctx context.Context, params ...ModuleParams) ([]Module, error)
	// CreateModuleVersion retrieves a module from the platform using the GraphQL interface
	CreateModuleVersion(ctx context.Context, input ModuleVersion) (*ModuleVersion, error)

	// GetOrgs retrieves orgs from the platform using the GraphQL interface
	GetOrgs(ctx context.Context, params ...*OrgParams) ([]Org, error)
	// GetOrg retrieves an org from the platform using the GraphQL interface
	GetOrg(ctx context.Context, id int) (*Org, error)
	// CreateOrg creates an org on the platform using the GraphQL interface
	CreateOrg(ctx context.Context, org Org) (*Org, error)
	// UpdateOrg updates an org on the platform using the GraphQL interface
	UpdateOrg(ctx context.Context, org Org) (*Org, error)
	// DeleteOrg deletes an org on the platform using the GraphQL interface
	DeleteOrg(ctx context.Context, id int) error

	// GetAsset retrieves an asset from the platform using the GraphQL interface
	GetAsset(ctx context.Context, id int) (*Asset, error)
	// GetAssets retrieves assets from the platform using the GraphQL interface
	GetAssets(ctx context.Context, params AssetParams) ([]Asset, error)
	// CreateAsset creates an asset on the platform using the GraphQL interface
	CreateAsset(ctx context.Context, asset *Asset) error
	// CreateAssetVersion creates an asset version on the platform using the GraphQL interface
	CreateAssetVersion(ctx context.Context, assetVersion *AssetVersion) error
	// UpdateAssetVersion updates an asset version on the platform using the GraphQL interface
	UpdateAssetVersion(ctx context.Context, assetVersion *AssetVersion) error

	// PostAsset posts an asset to the platform using the REST API
	PostAsset(ctx context.Context, assetVersion *AssetVersion) error
	// RetrieveAssets retrieves assets from the platform using the REST API
	RetrieveAssets(ctx context.Context, params AssetParams) ([]Asset, error)

	// GetSession retrieves a session from the platform using the GraphQL interface
	GetSession(ctx context.Context, id int) (*Session, error)
	// CreateSession creates a session on the platform using the GraphQL interface
	CreateSession(ctx context.Context, session *Session) error
	// UpdateSession updates a session on the platform using the GraphQL interface
	UpdateSession(ctx context.Context, session Session) (*Session, error)
	// CreateEvent creates an event on the platform using the GraphQL interface
	CreateEvent(ctx context.Context, event *Event) error

	// GetMultiplayerServerConfigs retrieves multiplayer server configs from the platform using the GraphQL interface
	GetMultiplayerServerConfigs(ctx context.Context, params *MultiplayerServerConfigParams) ([]MultiplayerServerConfigQueryParams, error)
	// GetMultiplayerServerVersions retrieves multiplayer server versions from the platform using the GraphQL interface
	GetMultiplayerServerVersions(ctx context.Context, params *MultiplayerServerVersionParams) ([]MultiplayerServerVersion, error)
	// GetMultiplayerServerVersionsWithConfig retrieves multiplayer server versions with config from the platform using the GraphQL interface
	GetMultiplayerServerVersionsWithConfig(ctx context.Context, params *MultiplayerServerVersionParams) ([]MultiplayerServerVersion, error)
	// GetMultiplayerServerVersion retrieves a multiplayer server version from the platform using the GraphQL interface
	GetMultiplayerServerVersion(ctx context.Context, id int) (*MultiplayerServerVersion, error)
	// CreateMultiplayerServerVersion creates a multiplayer server version on the platform using the GraphQL interface
	CreateMultiplayerServerVersion(ctx context.Context, input MultiplayerServerVersion) (*MultiplayerServerVersion, error)
	// UpdateMultiplayerServerVersion updates a multiplayer server version on the platform using the GraphQL interface
	UpdateMultiplayerServerVersion(ctx context.Context, input MultiplayerServerVersion) (*MultiplayerServerVersion, error)
}
