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

	query := `fragment ModuleVersionFragment on ModuleVersion { id moduleId status version fileLink createdAt updatedAt package controls { id name } platforms { id name shortName } module { id abbreviation isMultiplayer description details createdBy updatedBy isPublic status distributorId distributor { id name } createdAt updatedAt } createdBy updatedBy } query moduleVersions($params: ModuleVersionParamsInput) { moduleVersions(params: $params) { ...ModuleVersionFragment } }`
	variables := map[string]interface{}{
		"params": map[string]interface{}{
			"status":             params.Status,
			"moduleId":           params.ModuleID,
			"platformShortNames": params.PlatformShortNames,
			"sortField":          params.SortField,
			"sortOrder":          params.SortOrder,
		},
	}

	var moduleVersionsReponse ModuleVersionsResponse
	if err := p.Exec(ctx, query, &moduleVersionsReponse, variables); err != nil {
		return nil, err
	}

	return moduleVersionsReponse.ModuleVersions, nil
}
