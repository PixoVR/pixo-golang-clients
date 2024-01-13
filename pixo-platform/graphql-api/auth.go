package graphql_api

import (
	"encoding/json"
	"errors"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		Role string `json:"role"`
	}
}

// Login performs a login request to the API
func (g *GraphQLAPIClient) Login(username, password string) error {
	url := g.GetURLWithPath("auth/login")

	loginInput := LoginRequest{
		Login:    username,
		Password: password,
	}

	res, err := g.FormatRequest().
		SetHeader("Content-Type", "application/json").
		SetBody(loginInput).
		Post(url)
	if err != nil {
		return err
	}

	if res.IsError() {
		return errors.New(string(res.Body()))
	}

	var loginResponse LoginResponse
	if err = json.Unmarshal(res.Body(), &loginResponse); err != nil {
		return errors.New(string(res.Body()))
	}

	g.SetToken(loginResponse.Token)

	return nil
}
