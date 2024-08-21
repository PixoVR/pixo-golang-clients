package platform

import (
	"context"
	abstract "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
)

type Client interface {
	abstract.AbstractClient

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

	GetPlatforms(ctx context.Context) ([]*Platform, error)
	GetControlTypes(ctx context.Context) ([]*ControlType, error)

	GetModules(ctx context.Context, params ...ModuleParams) ([]Module, error)
	CreateModuleVersion(ctx context.Context, input ModuleVersion) (*ModuleVersion, error)

	GetOrgs(ctx context.Context, params ...*OrgParams) ([]Org, error)
	GetOrg(ctx context.Context, id int) (*Org, error)
	CreateOrg(ctx context.Context, org Org) (*Org, error)
	UpdateOrg(ctx context.Context, org Org) (*Org, error)
	DeleteOrg(ctx context.Context, id int) error

	GetSession(ctx context.Context, id int) (*Session, error)
	CreateSession(ctx context.Context, session *Session) error
	UpdateSession(ctx context.Context, session Session) (*Session, error)
	CreateEvent(ctx context.Context, event *Event) error

	GetMultiplayerServerConfigs(ctx context.Context, params *MultiplayerServerConfigParams) ([]*MultiplayerServerConfigQueryParams, error)
	GetMultiplayerServerVersions(ctx context.Context, params *MultiplayerServerVersionQueryParams) ([]*MultiplayerServerVersion, error)
	GetMultiplayerServerVersion(ctx context.Context, id int) (*MultiplayerServerVersion, error)
	CreateMultiplayerServerVersion(ctx context.Context, input MultiplayerServerVersion) (*MultiplayerServerVersion, error)
}
