package legacy

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
)

// Login performs a login request to the API
func (p *Client) Login(username, password string) error {
	url := p.GetURLWithPath("login")

	loginInput := LoginRequest{
		Login:    username,
		Password: password,
	}

	res, err := p.NewRequest().
		SetHeader("Content-Type", "application/json").
		SetBody(loginInput).
		Post(url)
	if err != nil {
		log.Error().Err(err).Msg("Failed to perform login request")
		return err
	}

	if res.IsError() {
		log.Error().Bytes("body", res.Body()).Msg("Login attempt failed")
		return errors.New(string(res.Body()))
	}

	var loginResponse LegacyAuthResponse
	if err = json.Unmarshal(res.Body(), &loginResponse); err != nil {
		log.Error().Err(err).Msg("Failed to login")
		return errors.New(string(res.Body()))
	}

	p.SetToken(loginResponse.User.Token)
	p.SetHeader("x-access-token", loginResponse.User.Token)
	return nil
}
