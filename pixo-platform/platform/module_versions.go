package platform

import (
	"context"
)

type ModuleVersionParams struct {
	ModuleID           *int     `json:"moduleId,omitempty"`
	ModulePlayerID     *int     `json:"modulePlayerId,omitempty"`
	DistributorID      *int     `json:"distributorId,omitempty"`
	Version            *string  `json:"version,omitempty"`
	PlatformIds        []int    `json:"platformIds,omitempty"`
	PlatformShortNames []string `json:"platformShortNames,omitempty"`
	Lifecycles         []string `json:"lifecycles,omitempty"`
	Model              *string  `json:"model,omitempty"`
	SortField          *string  `json:"sortField,omitempty"`
	SortOrder          *string  `json:"sortOrder,omitempty"`
}

type ModuleVersionsResponse struct {
	ModuleVersions []ModuleVersion `json:"moduleVersions"`
}

func (p *clientImpl) GetModuleVersions(ctx context.Context, params *ModuleVersionParams) ([]ModuleVersion, error) {
	if params == nil {
		params = &ModuleVersionParams{}
	}

	query := `query moduleVersions($params: ModuleVersionParamsInput) { moduleVersions(params: $params) { id moduleId lifecycleId lifecycle { id name } version fileLink package module { id abbreviation description } } }`
	variables := map[string]interface{}{
		"params": params,
	}

	var moduleVersionsReponse ModuleVersionsResponse
	if err := p.Exec(ctx, query, &moduleVersionsReponse, variables); err != nil {
		return nil, err
	}

	return moduleVersionsReponse.ModuleVersions, nil
}
