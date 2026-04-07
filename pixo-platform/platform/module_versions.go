package platform

import (
	"context"
)

type ModuleVersionParams struct {
	Status             []string `json:"status,omitempty"`
	ModuleID           *int     `json:"moduleId,omitempty"`
	PlatformShortNames []string `json:"platformShortNames,omitempty"`
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

	query := `query moduleVersions($params: ModuleVersionParams) { moduleVersions(params: $params) { id moduleId status version fileLink package module { id abbreviation description } } }`
	variables := map[string]interface{}{
		"params": params,
	}

	var moduleVersionsReponse ModuleVersionsResponse
	if err := p.Exec(ctx, query, &moduleVersionsReponse, variables); err != nil {
		return nil, err
	}

	return moduleVersionsReponse.ModuleVersions, nil
}
