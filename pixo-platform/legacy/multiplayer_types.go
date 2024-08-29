package legacy

type MultiplayerServerConfig struct {
	ID              int    `json:"id"`
	Capacity        int    `json:"capacity,omitempty"`
	StandbyReplicas string `json:"standbyReplicas,omitempty"`
	Disabled        bool   `json:"disabled,omitempty"`

	ModuleID int     `json:"moduleId,omitempty"`
	Module   *Module `json:"module,omitempty"`

	ServerVersions []*MultiplayerServerVersion `json:"serverVersions,omitempty"`

	CreatedBy string `json:"createdBy,omitempty"`
	UpdatedBy string `json:"updatedBy,omitempty"`

	//CreatedAt time.Time `json:"createdAt,omitempty"`
	//UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type MultiplayerServerVersion struct {
	ID              int    `json:"id,omitempty" graphql:"id"`
	Engine          string `json:"engine,omitempty" graphql:"engine"`
	Status          string `json:"status,omitempty" graphql:"status"`
	ImageRegistry   string `json:"imageRegistry" graphql:"imageRegistry"`
	SemanticVersion string `json:"semanticVersion,omitempty" graphql:"semanticVersion"`
	FileLink        string `json:"fileLink,omitempty" graphql:"fileLink"`

	ModuleID int     `json:"moduleId,omitempty" graphql:"moduleId"`
	Module   *Module `json:"module,omitempty" graphql:"module"`

	CreatedBy string `json:"createdBy" graphql:"createdBy"`
	UpdatedBy string `json:"updatedBy" graphql:"updatedBy"`

	//CreatedAt time.Time `json:"createdAt,omitempty" graphql:"createdAt"`
	//UpdatedAt time.Time `json:"updatedAt,omitempty" graphql:"updatedAt"`
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
