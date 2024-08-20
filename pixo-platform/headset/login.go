package headset

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
}

func (c *client) Login(username, password string) error {
	url := c.GetURLWithPath("login")

	loginInput := LoginRequest{
		Login:    username,
		Password: password,
	}

	res, err := c.FormatRequest().
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