package primary_api

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
)

// Login performs a login request to the API
func (p *PrimaryAPIClient) Login(username, password string) error {
	url := p.GetURLWithPath("login")

	loginInput := LoginRequest{
		Login:    username,
		Password: password,
	}

	res, err := p.FormatRequest().
		SetHeader("Content-Type", "application/json").
		SetBody(loginInput).
		Post(url)
	if err != nil {
		log.Error().Err(err).Msg("Failed to perform login request")
		return err
	}

	if res.IsError() {
		log.Error().Err(err).Msg("Login attempt failed")
		return errors.New(string(res.Body()))
	}

	var loginResponse AuthResponse
	if err = json.Unmarshal(res.Body(), &loginResponse); err != nil {
		log.Error().Err(err).Msg("Failed to login")
		return errors.New(string(res.Body()))
	}

	p.SetToken(loginResponse.User.Token)
	return nil
}
