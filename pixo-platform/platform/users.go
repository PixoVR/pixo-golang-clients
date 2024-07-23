package platform

import (
	"context"
	"encoding/json"
	"errors"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/legacy"
)

type GetUserResponse struct {
	User platform.User `json:"user"`
}

type CreateUserResponse struct {
	User platform.User `json:"createUser"`
}

type UpdateUserResponse struct {
	User platform.User `json:"updateUser"`
}

type DeleteUserResponse struct {
	Success bool `json:"deleteUser"`
}

func (g *GraphQLAPIClient) CreateUser(ctx context.Context, user platform.User) (*platform.User, error) {
	query := `mutation createUser($input: UserInput!) { createUser(input: $input) { id orgId firstName lastName username role } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"username":  user.Username,
			"password":  user.Password,
			"orgId":     user.OrgID,
			"role":      user.Role,
		},
	}

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var userResponse CreateUserResponse
	if err = json.Unmarshal(res, &userResponse); err != nil {
		return nil, err
	}

	return &userResponse.User, nil
}

func (g *GraphQLAPIClient) UpdateUser(ctx context.Context, user platform.User) (*platform.User, error) {

	if user.ID == 0 {
		return nil, errors.New("user id is required")
	}

	query := `mutation updateUser($input: UserInput!) { updateUser(input: $input) { id firstName lastName username role orgId } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"id": user.ID,
		},
	}

	if user.FirstName != "" {
		variables["input"].(map[string]interface{})["firstName"] = user.FirstName
	}

	if user.LastName != "" {
		variables["input"].(map[string]interface{})["lastName"] = user.LastName
	}

	if user.Username != "" {
		variables["input"].(map[string]interface{})["username"] = user.Username
	}

	if user.Password != "" {
		variables["input"].(map[string]interface{})["password"] = user.Password
	}

	if user.OrgID != 0 {
		variables["input"].(map[string]interface{})["orgId"] = user.OrgID
	}

	if user.Role != "" {
		variables["input"].(map[string]interface{})["role"] = user.Role
	}

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var userResponse UpdateUserResponse
	if err = json.Unmarshal(res, &userResponse); err != nil {
		return nil, err
	}

	return &userResponse.User, nil
}

func (g *GraphQLAPIClient) DeleteUser(ctx context.Context, id int) error {
	query := `mutation deleteUser($id: ID!) { deleteUser(id: $id) }`

	variables := map[string]interface{}{
		"id": id,
	}

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return err
	}

	var deleteResponse DeleteUserResponse
	if err = json.Unmarshal(res, &deleteResponse); err != nil {
		return err
	}

	if !deleteResponse.Success {
		return errors.New("failed to delete user")
	}

	return nil
}

func (g *GraphQLAPIClient) GetUserByUsername(ctx context.Context, username string) (*platform.User, error) {
	query := `query user($id: ID, $username: String) { user(id: $id, username: $username) { id username firstName lastName orgId role } }`

	variables := map[string]interface{}{
		"username": username,
	}

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var userResponse GetUserResponse
	if err = json.Unmarshal(res, &userResponse); err != nil {
		return nil, err
	}

	return &userResponse.User, nil
}
