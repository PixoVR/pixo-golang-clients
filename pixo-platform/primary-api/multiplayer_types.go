package primary_api

type MultiplayerServerConfig struct {
	ID       int   `json:"id"`
	ModuleID int   `json:"moduleId"`
	Capacity int32 `json:"capacity"`
	Disabled bool  `json:"disabled"`

	CreatedBy string `json:"createdBy"`
	UpdatedBy string `json:"updatedBy"`

	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type MultiplayerServerVersion struct {
	ID               int    `json:"id,omitempty" graphql:"id"`
	ModuleID         int    `json:"moduleId,omitempty" graphql:"moduleId"`
	Engine           string `json:"engine,omitempty" graphql:"engine"`
	Status           string `json:"status,omitempty" graphql:"status"`
	ImageRegistry    string `json:"imageRegistry" graphql:"imageRegistry"`
	SemanticVersion  string `json:"semanticVersion,omitempty" graphql:"semanticVersion"`
	MinClientVersion string `json:"minClientVersion,omitempty" graphql:"minClientVersion"`
	Filename         string `json:"filename,omitempty" graphql:"filename"`

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
