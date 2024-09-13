package platform

import (
	"context"
	"errors"
	"time"
)

type Org struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Status      string    `json:"enabled"`
	LogoLink    string    `json:"logoLink"`
	HubLogoLink string    `json:"hubLogoLink"`
	OpenAccess  bool      `json:"openAccess"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type OrgParams struct {
	Name string `json:"name"`
}

type GetOrgsResponse struct {
	Orgs []Org `json:"orgs"`
}

type GetOrgResponse struct {
	Org Org `json:"org"`
}

type CreateOrgResponse struct {
	Org Org `json:"createOrg"`
}

type UpdateOrgResponse struct {
	Org Org `json:"updateOrg"`
}

type DeleteOrgResponse struct {
	Success bool `json:"deleteOrg"`
}

func (p *clientImpl) GetOrgs(ctx context.Context, params ...*OrgParams) ([]Org, error) {
	query := `query orgs { orgs { id name type openAccess logoLink hubLogoLink } }`

	variables := map[string]interface{}{
		"params": map[string]interface{}{},
	}
	if len(params) > 0 && params[0] != nil {
		if params[0].Name != "" {
			variables["params"].(map[string]interface{})["name"] = params[0].Name
		}
	}

	var res GetOrgsResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return res.Orgs, nil
}

func (p *clientImpl) GetOrg(ctx context.Context, id int) (*Org, error) {
	query := `query org($id: ID!) { org(id: $id) { id name type openAccess logoLink hubLogoLink } }`

	variables := map[string]interface{}{
		"id": id,
	}

	var res GetOrgResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return &res.Org, nil
}

func (p *clientImpl) CreateOrg(ctx context.Context, org Org) (*Org, error) {
	query := `mutation createOrg($input: OrgInput!) { createOrg(input: $input) { id name } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"name":       org.Name,
			"type":       org.Type,
			"openAccess": org.OpenAccess,
		},
	}

	var res CreateOrgResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return &res.Org, nil
}

func (p *clientImpl) UpdateOrg(ctx context.Context, org Org) (*Org, error) {

	if org.ID == 0 {
		return nil, errors.New("org id is required")
	}

	query := `mutation updateOrg($input: OrgInput!) { updateOrg(input: $input) { id name } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"id": org.ID,
		},
	}

	if org.Name != "" {
		variables["input"].(map[string]interface{})["name"] = org.Name
	}

	var res UpdateOrgResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return &res.Org, nil
}

func (p *clientImpl) DeleteOrg(ctx context.Context, id int) error {
	query := `mutation deleteOrg($id: ID!) { deleteOrg(id: $id) }`

	variables := map[string]interface{}{
		"id": id,
	}

	var res DeleteOrgResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return err
	}

	if !res.Success {
		return errors.New("failed to delete user")
	}

	return nil
}
