package platform

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/legacy"
)

type APIKeyQueryParams struct {
	UserID *int `json:"userId" graphql:"userId"`
}

type GetAPIKeysResponse struct {
	APIKeys []*legacy.APIKey `json:"apiKeys"`
}

type CreateAPIKeyResponse struct {
	APIKey legacy.APIKey `json:"createApiKey"`
}

type DeleteAPIKeyResponse struct {
	Success bool `json:"deleteApiKey"`
}

func (g *PlatformClient) CreateAPIKey(ctx context.Context, input legacy.APIKey) (*legacy.APIKey, error) {
	query := `mutation createApiKey($input: ApiKeyInput!) { createApiKey(input: $input) { id key userId user { role } } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{},
	}

	if input.UserID > 0 {
		variables["input"] = map[string]interface{}{
			"userId": input.UserID,
		}
	}

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var apiKeyResponse CreateAPIKeyResponse
	if err = json.Unmarshal(res, &apiKeyResponse); err != nil {
		return nil, err
	}

	return &apiKeyResponse.APIKey, nil
}

func (g *PlatformClient) GetAPIKeys(ctx context.Context, params *APIKeyQueryParams) ([]*legacy.APIKey, error) {
	query := `query apiKeys($params: ApiKeyParams) { apiKeys(params: $params) { id key userId user { role } } }`

	variables := map[string]interface{}{
		"params": map[string]interface{}{},
	}

	if params.UserID != nil {
		variables["params"] = map[string]interface{}{
			"userId": *params.UserID,
		}
	}

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var apiKeysResponse GetAPIKeysResponse
	if err = json.Unmarshal(res, &apiKeysResponse); err != nil {
		return nil, err
	}

	return apiKeysResponse.APIKeys, nil
}

func (g *PlatformClient) DeleteAPIKey(ctx context.Context, id int) error {
	query := `mutation deleteApiKey($id: ID!) { deleteApiKey(id: $id) }`

	variables := map[string]interface{}{
		"id": id,
	}

	res, err := g.Client.ExecRaw(ctx, query, variables)
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
