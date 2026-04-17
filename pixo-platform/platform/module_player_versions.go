package platform

import (
	"context"
	"time"
)

type ModulePlayerVersion struct {
	ID             int           `json:"id,omitempty"`
	ModulePlayerID int           `json:"modulePlayerId,omitempty"`
	ModulePlayer   *ModulePlayer `json:"modulePlayer,omitempty"`
	Status         string        `json:"status,omitempty"`
	Version        string        `json:"version,omitempty"`
	FileLink       string        `json:"fileLink,omitempty"`
	Package        string        `json:"package,omitempty"`
	Controls       []ControlType `json:"controls,omitempty"`
	Platforms      []Platform    `json:"platforms,omitempty"`
	CreatedBy      int           `json:"createdBy,omitempty"`
	UpdatedBy      int           `json:"updatedBy,omitempty"`
	CreatedAt      time.Time     `json:"createdAt,omitempty"`
	UpdatedAt      time.Time     `json:"updatedAt,omitempty"`
}

type ModulePlayerVersionParams struct {
	Status             []string `json:"status,omitempty"`
	ModulePlayerID     *int     `json:"modulePlayerId,omitempty"`
	PlatformShortNames []string `json:"platformShortNames,omitempty"`
	SortField          *string  `json:"sortField,omitempty"`
	SortOrder          *string  `json:"sortOrder,omitempty"`
}

type ModulePlayerVersionsResponse struct {
	ModulePlayerVersions []ModulePlayerVersion `json:"modulePlayerVersions"`
}

func (p *clientImpl) GetModulePlayerVersions(ctx context.Context, params *ModulePlayerVersionParams) ([]ModulePlayerVersion, error) {
	if params == nil {
		params = &ModulePlayerVersionParams{}
	}

	query := `query modulePlayerVersions($params: ModulePlayerVersionParamsInput) { modulePlayerVersions(params: $params) { id modulePlayerId status version fileLink package controls { id name } platforms { id name shortName } modulePlayer { id name description distributor { id name } createdAt updatedAt } createdBy updatedBy createdAt updatedAt } }`
	variables := map[string]interface{}{
		"params": params,
	}

	var res ModulePlayerVersionsResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return res.ModulePlayerVersions, nil
}
