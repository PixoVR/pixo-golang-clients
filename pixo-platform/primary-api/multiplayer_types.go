package primary_api

import "time"

type MultiplayerServerConfig struct {
	ID       int   `json:"id"`
	ModuleID int   `json:"moduleId"`
	Capacity int32 `json:"capacity"`
	Disabled bool  `json:"disabled"`

	CreatedBy string `json:"createdBy"`
	UpdatedBy string `json:"updatedBy"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type MultiplayerServerVersion struct {
	ModuleID         int    `json:"moduleId,omitempty"`
	VersionID        int    `json:"versionId,omitempty"`
	Engine           string `json:"engine,omitempty"`
	Status           string `json:"status,omitempty"`
	ImageRegistry    string `json:"imageRegistry"`
	Version          string `json:"version,omitempty"`
	MinClientVersion string `json:"minClientVersion,omitempty"`
	Filename         string `json:"filename,omitempty"`

	CreatedBy string `json:"createdBy"`
	UpdatedBy string `json:"updatedBy"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GameProfileMetaDataResponse struct {
	Error      bool                   `json:"error"`
	StatusCode int                    `json:"statusCode"`
	Message    string                 `json:"message"`
	Profiles   []*GameProfileMetadata `json:"profiles"`
}

type GameProfileMetadata struct {
	OrgID         int                      `json:"orgId"`
	Config        MultiplayerServerConfig  `json:"config"`
	ServerVersion MultiplayerServerVersion `json:"serverVersion"`
}
