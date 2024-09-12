package platform

import (
	"context"
	"encoding/json"
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

	res, err := p.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var apiKeyResponse CreateAPIKeyResponse
	if err = json.Unmarshal(res, &apiKeyResponse); err != nil {
		return nil, err
	}

	return &apiKeyResponse.APIKey, nil
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

	res, err := p.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var apiKeysResponse GetAPIKeysResponse
	if err = json.Unmarshal(res, &apiKeysResponse); err != nil {
		return nil, err
	}

	return apiKeysResponse.APIKeys, nil
}

func (p *clientImpl) DeleteAPIKey(ctx context.Context, id int) error {
	query := `mutation deleteApiKey($id: ID!) { deleteApiKey(id: $id) }`

	variables := map[string]interface{}{
		"id": id,
	}

	res, err := p.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return err
	}

	var deleteResponse DeleteAPIKeyResponse
	if err = json.Unmarshal(res, &deleteResponse); err != nil {
		return err
	}

	if !deleteResponse.Success {
		return errors.New("failed to delete api key")
	}

	return nil
}
