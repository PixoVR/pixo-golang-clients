package platform

import (
	"context"
	"fmt"
)

// moduleVersionFields is the selection every module version read shares.
const moduleVersionFields = `
	id
	moduleId
	version
	notes
	package
	public
	fileLink
	filePath
	fileSize
	uploadStatus
	externalId
	lifecycleId
	lifecycle { id name }
	platforms { id name shortName }
	createdAt
	updatedAt
`

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

	query := fmt.Sprintf(
		`query moduleVersions($params: ModuleVersionParamsInput) { moduleVersions(params: $params) { %s module { %s } } }`,
		moduleVersionFields,
		moduleFields,
	)
	variables := map[string]interface{}{
		"params": params,
	}

	var moduleVersionsReponse ModuleVersionsResponse
	if err := p.Exec(ctx, query, &moduleVersionsReponse, variables); err != nil {
		return nil, err
	}

	return moduleVersionsReponse.ModuleVersions, nil
}
