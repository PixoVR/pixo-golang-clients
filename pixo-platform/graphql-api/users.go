package graphql_api

import (
	"context"
	"encoding/json"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
)

type UsersClient interface {
	CreateUser(ctx context.Context, username, password string, orgID int) (*platform.User, error)
}

type CreateUserResponse struct {
	User platform.User `json:"createUser"`
}

func (g *GraphQLAPIClient) CreateUser(ctx context.Context, username, password string, orgID int) (*platform.User, error) {
	query := `mutation createUser($input: UserInput!) { createUser(input: $input) { id } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"username": username,
			"password": password,
			"orgId":    orgID,
		},
	}

	res, err := g.gqlClient.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var userResponse CreateUserResponse
	if err = json.Unmarshal(res, &userResponse); err != nil {
		return nil, err
	}

	return &userResponse.User, nil
}
