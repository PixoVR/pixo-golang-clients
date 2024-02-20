package graphql_api

import (
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
)

type MultiplayerServerConfigQuery struct {
	MultiplayerServerConfigs []*MultiplayerServerConfigQueryParams `graphql:"multiplayerServerConfigs(params: $params)"`
}

type MultiplayerServerVersionQuery struct {
	MultiplayerServerVersions []*platform.MultiplayerServerVersion `graphql:"multiplayerServerVersions(params: $params)"`
}

type MultiplayerServerVersionParams struct {
	ModuleID        int    `json:"moduleId,omitempty" graphql:"moduleId"`
	ImageRegistry   string `json:"imageRegistry,omitempty" graphql:"image"`
	SemanticVersion string `json:"semanticVersion,omitempty" graphql:"semanticVersion"`
	Status          string `json:"status,omitempty" graphql:"status"`
	Engine          string `json:"engine,omitempty" graphql:"engine"`
}

type MultiplayerServerConfigParams struct {
	ModuleID       int                         `json:"moduleId,omitempty"`
	OrgID          int                         `json:"orgId,omitempty"`
	ServerVersion  string                      `json:"serverVersion,omitempty"`
	Capacity       int                         `json:"capacity,omitempty"`
	ServerVersions []*MultiplayerServerVersion `json:"serverVersions,omitempty"`
}

type MultiplayerServerConfigQueryParams struct {
	ID       int  `json:"id" graphql:"id"`
	ModuleID int  `json:"moduleId" graphql:"moduleId"`
	Capacity int  `json:"capacity" graphql:"capacity"`
	Disabled bool `json:"disabled" graphql:"disabled"`
	Module   struct {
		ID   int    `json:"id" graphql:"id"`
		Name string `json:"name" graphql:"name"`
	}
	ServerVersions []*MultiplayerServerVersion `json:"serverVersions" graphql:"serverVersions"`
}

type MultiplayerServerVersion struct {
	ID              int             `json:"id" graphql:"id"`
	ModuleID        int             `json:"moduleId" graphql:"moduleId"`
	Module          platform.Module `json:"module" graphql:"module"`
	Engine          string          `json:"engine" graphql:"engine"`
	SemanticVersion string          `json:"semanticVersion" graphql:"semanticVersion"`
	ImageRegistry   string          `json:"imageRegistry" graphql:"imageRegistry"`
	Status          string          `json:"status" graphql:"status"`
}

type MultiplayerServerVersionQueryParams struct {
	ModuleID        int    `json:"moduleId" graphql:"moduleId"`
	SemanticVersion string `json:"semanticVersion" graphql:"semanticVersion"`
}

type MultiplayerServerConfigInput struct {
	Input platform.MultiplayerServerConfig `graphql:"createMultiplayerConfig($input: MultiplayerServerConfigInput!)"`
}

type CreateMultiplayerServerConfigResponse struct {
	CreateMultiplayerServerConfig platform.MultiplayerServerConfig `graphql:"createMultiplayerServerConfig(input: $input)"`
}

type MultiplayerServerVersionInput struct {
	Input platform.MultiplayerServerVersion `graphql:"createMultiplayerServer($input: MultiplayerServerVersionInput!)"`
}

type CreateMultiplayerServerVersionResponse struct {
	CreateMultiplayerServerVersion platform.MultiplayerServerVersion `graphql:"createMultiplayerServerVersion(input: $input)"`
}
