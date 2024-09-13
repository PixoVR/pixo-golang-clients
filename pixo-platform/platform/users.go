package platform

import (
	"context"
	"errors"
	"time"
)

type User struct {
	ID         int    `json:"id,omitempty"`
	Role       string `json:"role,omitempty"`
	OrgID      int    `json:"orgId,omitempty"`
	Org        Org    `json:"org,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"-"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	Email      string `json:"email,omitempty"`
	ExternalID string `json:"externalId,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type GetUserResponse struct {
	User User `json:"user"`
}

type CreateUserResponse struct {
	User User `json:"createUser"`
}

type UpdateUserResponse struct {
	User User `json:"updateUser"`
}

type DeleteUserResponse struct {
	Success bool `json:"deleteUser"`
}

func (p *clientImpl) CreateUser(ctx context.Context, user *User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	query := `mutation createUser($input: UserInput!) { createUser(input: $input) { id firstName lastName username email role orgId org { id name }  } }`
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"username":  user.Username,
			"email":     user.Email,
			"password":  user.Password,
			"orgId":     user.OrgID,
			"role":      user.Role,
		},
	}

	var userResponse CreateUserResponse
	if err := p.Exec(ctx, query, &userResponse, variables); err != nil {
		return err
	}

	*user = userResponse.User
	return nil
}

func (p *clientImpl) UpdateUser(ctx context.Context, user *User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	if user.ID == 0 {
		return errors.New("user id is required")
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

	var userResponse UpdateUserResponse
	if err := p.Exec(ctx, query, &userResponse, variables); err != nil {
		return err
	}

	*user = userResponse.User
	return nil
}

func (p *clientImpl) DeleteUser(ctx context.Context, id int) error {
	query := `mutation deleteUser($id: ID!) { deleteUser(id: $id) }`

	variables := map[string]interface{}{
		"id": id,
	}

	var deleteResponse DeleteUserResponse
	if err := p.Exec(ctx, query, &deleteResponse, variables); err != nil {
		return err
	}

	if !deleteResponse.Success {
		return errors.New("failed to delete user")
	}

	return nil
}

func (p *clientImpl) GetUser(ctx context.Context, id int) (*User, error) {
	query := `query user($id: ID, $username: String) { user(id: $id, username: $username) { id username email firstName lastName orgId org { id name type } role } }`

	variables := map[string]interface{}{
		"id": id,
	}

	var userResponse GetUserResponse
	if err := p.Exec(ctx, query, &userResponse, variables); err != nil {
		return nil, err
	}

	return &userResponse.User, nil
}

func (p *clientImpl) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	query := `query user($id: ID, $username: String) { user(id: $id, username: $username) { id username firstName lastName orgId role } }`

	variables := map[string]interface{}{
		"username": username,
	}

	var userResponse GetUserResponse
	if err := p.Exec(ctx, query, &userResponse, variables); err != nil {
		return nil, err
	}

	return &userResponse.User, nil
}
