package graphql_api

import platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"

type MultiplayerServerConfigQuery struct {
	MultiplayerServerConfigs []*MultiplayerServerConfigQueryParams `graphql:"multiplayerServerConfigs(params: $params)"`
}

type MultiplayerServerConfigParams struct {
	ModuleID      int    `json:"moduleId,omitempty"`
	OrgID         int    `json:"orgId,omitempty"`
	ServerVersion string `json:"serverVersion,omitempty"`
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
	ServerVersions []struct {
		Engine          string `json:"engine" graphql:"engine"`
		SemanticVersion string `json:"semanticVersion" graphql:"semanticVersion"`
		ImageRegistry   string `json:"imageRegistry" graphql:"imageRegistry"`
		Status          string `json:"status" graphql:"status"`
	}
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
