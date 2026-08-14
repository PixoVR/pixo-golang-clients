package platform

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type OrgModule struct {
	ID             int        `json:"id,omitempty"`
	OrgID          int        `json:"orgId,omitempty"`
	ModuleID       int        `json:"moduleId,omitempty"`
	Module         *Module    `json:"module,omitempty"`
	Lifetime       bool       `json:"lifetime,omitempty"`
	ExpiresAt      string     `json:"expiresAt,omitempty"`
	ExternalID     string     `json:"externalId,omitempty"`
	LearningType   string     `json:"learningType,omitempty"`
	CompleteStatus string     `json:"completeStatus,omitempty"`
	ExpirDate      *time.Time `json:"-"`
}

// OrgModuleParams identifies the org whose modules should be returned.
type OrgModuleParams struct {
	OrgID          int  `json:"orgId"`
	IncludeExpired bool `json:"includeExpired"`
}

type GetOrgModulesResponse struct {
	OrgModules []OrgModule `json:"orgModules"`
}

type CreateOrgModuleResponse struct {
	OrgModule OrgModule `json:"createOrgModule"`
}

// GetOrgModules retrieves the modules an org has access to, including the
// access type that grants it.
func (p *clientImpl) GetOrgModules(ctx context.Context, params OrgModuleParams) ([]OrgModule, error) {
	if params.OrgID == 0 {
		return nil, errors.New("org id is required")
	}

	query := fmt.Sprintf(
		`query orgModules($params: OrgModuleParams!) { orgModules(params: $params) { id orgId lifetime expiresAt externalId learningType completeStatus accessType module { %s } } }`,
		moduleFields,
	)

	variables := map[string]interface{}{
		"params": params,
	}

	var res GetOrgModulesResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return res.OrgModules, nil
}

type DeleteOrgModuleResponse struct {
	Success bool `json:"deleteOrgModule"`
}

func (p *clientImpl) CreateOrgModule(ctx context.Context, input OrgModule) (*OrgModule, error) {
	if input.OrgID == 0 || input.ModuleID == 0 {
		return nil, errors.New("org id and module id are required")
	}

	query := `mutation createOrgModule($input: OrgModuleInput!) { createOrgModule(input: $input) { id orgId lifetime expiresAt module { id abbreviation } } }`

	orgModuleInput := map[string]interface{}{
		"orgId":    input.OrgID,
		"moduleId": input.ModuleID,
		"lifetime": input.Lifetime,
	}

	if input.ExpirDate != nil {
		orgModuleInput["expirDate"] = input.ExpirDate.Format(time.RFC3339)
	}

	variables := map[string]interface{}{
		"input": orgModuleInput,
	}

	var res CreateOrgModuleResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return &res.OrgModule, nil
}

func (p *clientImpl) DeleteOrgModule(ctx context.Context, orgID, moduleID int) error {
	if orgID == 0 || moduleID == 0 {
		return errors.New("org id and module id are required")
	}

	query := `mutation deleteOrgModule($orgId: ID!, $moduleId: ID!) { deleteOrgModule(orgId: $orgId, moduleId: $moduleId) }`

	variables := map[string]interface{}{
		"orgId":    orgID,
		"moduleId": moduleID,
	}

	var res DeleteOrgModuleResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return err
	}

	if !res.Success {
		return errors.New("failed to delete org module")
	}

	return nil
}
