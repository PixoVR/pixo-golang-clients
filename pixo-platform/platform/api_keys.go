package platform

import (
	"context"
	"errors"
	"time"
)

type APIKey struct {
	ID     int    `json:"id,omitempty"`
	Key    string `json:"key,omitempty"`
	UserID int    `json:"userId,omitempty"`
	User   *User  `json:"user,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type APIKeyQueryParams struct {
	UserID *int `json:"userId" graphql:"userId"`
}

type GetAPIKeysResponse struct {
	APIKeys []APIKey `json:"apiKeys"`
}

type CreateAPIKeyResponse struct {
	APIKey APIKey `json:"createApiKey"`
}

type DeleteAPIKeyResponse struct {
	Success bool `json:"deleteApiKey"`
}

func (p *clientImpl) CreateAPIKey(ctx context.Context, input APIKey) (*APIKey, error) {
	query := `mutation createApiKey($input: ApiKeyInput!) { createApiKey(input: $input) { id key userId user { role } } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{},
	}

	if input.UserID > 0 {
		variables["input"] = map[string]interface{}{
			"userId": input.UserID,
		}
	}

	var res CreateAPIKeyResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return &res.APIKey, nil
}

func (p *clientImpl) GetAPIKeys(ctx context.Context, params *APIKeyQueryParams) ([]APIKey, error) {
	query := `query apiKeys($params: ApiKeyParams) { apiKeys(params: $params) { id key userId user { username email role } } }`

	variables := map[string]interface{}{
		"params": map[string]interface{}{},
	}

	if params != nil && params.UserID != nil {
		variables["params"] = map[string]interface{}{
			"userId": *params.UserID,
		}
	}

	var res GetAPIKeysResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return res.APIKeys, nil
}

func (p *clientImpl) DeleteAPIKey(ctx context.Context, id int) error {
	query := `mutation deleteApiKey($id: ID!) { deleteApiKey(id: $id) }`

	variables := map[string]interface{}{
		"id": id,
	}

	var res DeleteAPIKeyResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return err
	}

	if !res.Success {
		return errors.New("failed to delete api key")
	}

	return nil
}
