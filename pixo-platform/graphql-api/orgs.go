package graphql_api

import (
	"context"
	"encoding/json"
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

func (g *GraphQLAPIClient) GetOrg(ctx context.Context, id int) (*Org, error) {
	query := `query org($id: ID!) { org(id: $id) { id name type openAccess logoLink hubLogoLink } }`

	variables := map[string]interface{}{
		"id": id,
	}

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var orgResponse GetOrgResponse
	if err = json.Unmarshal(res, &orgResponse); err != nil {
		return nil, err
	}

	return &orgResponse.Org, nil
}

func (g *GraphQLAPIClient) CreateOrg(ctx context.Context, org Org) (*Org, error) {
	query := `mutation createOrg($input: OrgInput!) { createOrg(input: $input) { id name } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"name":       org.Name,
			"type":       org.Type,
			"openAccess": org.OpenAccess,
		},
	}

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var orgResponse CreateOrgResponse
	if err = json.Unmarshal(res, &orgResponse); err != nil {
		return nil, err
	}

	return &orgResponse.Org, nil
}

func (g *GraphQLAPIClient) UpdateOrg(ctx context.Context, org Org) (*Org, error) {

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

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var userResponse UpdateOrgResponse
	if err = json.Unmarshal(res, &userResponse); err != nil {
		return nil, err
	}

	return &userResponse.Org, nil
}

func (g *GraphQLAPIClient) DeleteOrg(ctx context.Context, id int) error {
	query := `mutation deleteOrg($id: ID!) { deleteOrg(id: $id) }`

	variables := map[string]interface{}{
		"id": id,
	}

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return err
	}

	var deleteResponse DeleteOrgResponse
	if err = json.Unmarshal(res, &deleteResponse); err != nil {
		return err
	}

	if !deleteResponse.Success {
		return errors.New("failed to delete user")
	}

	return nil
}
