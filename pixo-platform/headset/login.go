package headset

import (
	"encoding/json"
	"errors"
)

// LoginRequest is the request body for the login endpoint
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// LoginResponse is the response body for the login endpoint
type LoginResponse struct {
	Token string `json:"token"`
}

// Login logs in the user with the given username and password
func (c *client) Login(username, password string) error {
	url := c.GetURLWithPath("login")

	loginInput := LoginRequest{
		Login:    username,
		Password: password,
	}

	res, err := c.NewRequest().
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

	c.SetToken(loginResponse.Token)
	return nil
}
