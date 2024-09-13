package headset

import (
	"context"
	"encoding/json"
	"errors"
	"io"
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
	loginInput := LoginRequest{
		Login:    username,
		Password: password,
	}

	payload, _ := json.Marshal(loginInput)

	res, err := c.Post(context.TODO(), "login", payload)
	if err != nil {
		return err
	}

	resBody, _ := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		return errors.New(string(resBody))
	}

	var loginResponse LoginResponse
	if err = json.Unmarshal(resBody, &loginResponse); err != nil {
		return errors.New(string(resBody))
	}

	c.SetToken(loginResponse.Token)
	return nil
}
