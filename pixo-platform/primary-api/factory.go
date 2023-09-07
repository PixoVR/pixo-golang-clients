package primary_api

import (
	"encoding/json"
	"errors"
	abstractClient "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/rs/zerolog/log"
)

// PrimaryAPIClient is a struct for the primary API that contains an abstract client
type PrimaryAPIClient struct {
	abstractClient.PixoAbstractAPIClient
}

// NewClient is a function that returns a PixoAbstractAPIClient
func NewClient(token, apiURL string) *PrimaryAPIClient {
	return &PrimaryAPIClient{
		PixoAbstractAPIClient: *abstractClient.NewClient(token, apiURL),
	}
}

// NewClientWithBasicAuth is a function that returns a PixoAbstractAPIClient with basic auth performed
func NewClientWithBasicAuth(username, password, apiURL string) *PrimaryAPIClient {
	primaryClient := &PrimaryAPIClient{
		PixoAbstractAPIClient: *abstractClient.NewClient("", apiURL),
	}

	if err := primaryClient.Login(username, password); err != nil {
		log.Error().Err(err).Msg("Failed to login")
		return nil
	}

	return primaryClient
}

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
