package primary_api

type MultiplayerServerConfig struct {
	ID       int   `json:"id"`
	Capacity int32 `json:"capacity,omitempty"`
	Disabled bool  `json:"disabled,omitempty"`

	ModuleID int     `json:"moduleId,omitempty"`
	Module   *Module `json:"module,omitempty"`

	CreatedBy string `json:"createdBy,omitempty"`
	UpdatedBy string `json:"updatedBy,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

type MultiplayerServerTrigger struct {
	ID         int    `json:"id,omitempty"`
	Revision   string `json:"revision,omitempty"`
	Dockerfile string `json:"dockerfile,omitempty"`
	Context    string `json:"context,omitempty"`
	Config     string `json:"config,omitempty"`

	Module   *Module `json:"module,omitempty"`
	ModuleID int     `json:"moduleId,omitempty"`

	CreatedBy string `json:"createdBy,omitempty"`
	UpdatedBy string `json:"updatedBy,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

type MultiplayerServerVersion struct {
	ID               int    `json:"id,omitempty" graphql:"id"`
	Engine           string `json:"engine,omitempty" graphql:"engine"`
	Status           string `json:"status,omitempty" graphql:"status"`
	ImageRegistry    string `json:"imageRegistry" graphql:"imageRegistry"`
	SemanticVersion  string `json:"semanticVersion,omitempty" graphql:"semanticVersion"`
	MinClientVersion string `json:"minClientVersion,omitempty" graphql:"minClientVersion"`
	Filename         string `json:"filename,omitempty" graphql:"filename"`

	ModuleID int     `json:"moduleId,omitempty" graphql:"moduleId"`
	Module   *Module `json:"module,omitempty" graphql:"module"`

	CreatedBy string `json:"createdBy" graphql:"createdBy"`
	UpdatedBy string `json:"updatedBy" graphql:"updatedBy"`

	CreatedAt string `json:"createdAt" graphql:"createdAt"`
	UpdatedAt string `json:"updatedAt" graphql:"updatedAt"`
}

type GameProfileMetaDataResponse struct {
	Error      bool                   `json:"error"`
	StatusCode int                    `json:"statusCode"`
	Message    string                 `json:"message"`
	Profiles   []*GameProfileMetadata `json:"profiles"`

	CreatedBy string `json:"createdBy"`
	UpdatedBy string `json:"updatedBy"`

	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type GameProfileMetadata struct {
	OrgID         int                      `json:"orgId"`
	Config        MultiplayerServerConfig  `json:"config"`
	ServerVersion MultiplayerServerVersion `json:"serverVersion"`

	CreatedBy string `json:"createdBy"`
	UpdatedBy string `json:"updatedBy"`

	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
