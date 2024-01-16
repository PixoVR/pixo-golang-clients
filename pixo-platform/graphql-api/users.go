package graphql_api

import (
	"context"
	"encoding/json"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
)

type UsersClient interface {
	CreateUser(ctx context.Context, user platform.User) (*platform.User, error)
}

type CreateUserResponse struct {
	User platform.User `json:"createUser"`
}

func (g *GraphQLAPIClient) CreateUser(ctx context.Context, user platform.User) (*platform.User, error) {
	query := `mutation createUser($input: UserInput!) { createUser(input: $input) { id } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"username":  user.Username,
			"password":  user.Password,
			"orgId":     user.OrgID,
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
