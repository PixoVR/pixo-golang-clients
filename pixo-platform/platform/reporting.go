package platform

import (
	"context"
)

type ExpiringModule struct {
	OrgModuleID int    `json:"orgModuleId"`
	ExpiresAt   string `json:"expiresAt"`
	ModuleName  string `json:"moduleName"`
	Description string `json:"description"`
	ShortDesc   string `json:"shortDesc"`
	OrgID       int    `json:"orgId"`
	OrgName     string `json:"orgName"`
	Type        string `json:"type"`
}

type GetExpiringModulesResponse struct {
	ExpiringModules []ExpiringModule `json:"expiringModules"`
}

func (p *clientImpl) GetExpiringModules(ctx context.Context) ([]ExpiringModule, error) {
	query := `
query expiringModules {
    expiringModules {
        orgModuleId
        expiresAt
        moduleName
        description
        shortDesc
        orgId
        orgName
        type
    }
}`

	var res GetExpiringModulesResponse
	if err := p.Exec(ctx, query, &res, nil); err != nil {
		return nil, err
	}

	return res.ExpiringModules, nil
}
