package platform

import (
	"context"
	"encoding/json"
	"errors"
	"io"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// Login performs a login request to the API
func (p *clientImpl) Login(username, password string) error {
	loginInput := LoginRequest{
		Login:    username,
		Password: password,
	}

	payload, _ := json.Marshal(loginInput)

	res, err := p.Post(context.TODO(), "auth/login", payload)
	if err != nil {
		return err
	}

	resBody, _ := io.ReadAll(res.Body)

	if res.StatusCode > 299 {
		return errors.New(string(resBody))
	}

	var loginResponse LoginResponse
	if err = json.Unmarshal(resBody, &loginResponse); err != nil {
		return errors.New(string(resBody))
	}

	p.SetToken(loginResponse.Token)
	return nil
}
