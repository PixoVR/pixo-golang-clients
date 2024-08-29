package platform

type MultiplayerServerConfigQuery struct {
	MultiplayerServerConfigs []MultiplayerServerConfigQueryParams `graphql:"multiplayerServerConfigs(params: $params)"`
}

type MultiplayerServerVersionQuery struct {
	MultiplayerServerVersions []MultiplayerServerVersion `graphql:"multiplayerServerVersions(params: $params)"`
}

type MultiplayerServerConfigParams struct {
	ModuleID       int                        `json:"moduleId,omitempty"`
	OrgID          int                        `json:"orgId,omitempty"`
	ServerVersion  string                     `json:"serverVersion,omitempty"`
	Capacity       int                        `json:"capacity,omitempty"`
	ServerVersions []MultiplayerServerVersion `json:"serverVersions,omitempty"`
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
	ServerVersions []MultiplayerServerVersion `json:"serverVersions" graphql:"serverVersions"`
}

type MultiplayerServerVersionParams struct {
	ModuleID        int    `json:"moduleId" graphql:"moduleId"`
	SemanticVersion string `json:"semanticVersion" graphql:"semanticVersion"`
}
