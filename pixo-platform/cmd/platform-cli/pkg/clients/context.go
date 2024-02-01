package clients

import (
	"fmt"
	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/editor"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/loader"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/rs/zerolog/log"
	"io"
)

type PlatformContext struct {
	ConfigManager     config.Manager
	OldAPIClient      abstract_client.AuthClient
	PlatformClient    platform.PlatformClient
	MatchmakingClient matchmaker.Matchmaker
	FileOpener        editor.FileOpener
}

func (p *PlatformContext) Authenticate(reader io.Reader, writer io.Writer) error {
	token := p.ConfigManager.GetConfigValue("token")
	if token != "" {
		p.PlatformClient.SetToken(token)
		return nil
	}

	apiKey := p.ConfigManager.GetConfigValue("apiKey")
	if apiKey != "" {
		p.PlatformClient.SetAPIKey(apiKey)
		return nil
	}

	username := p.ConfigManager.GetConfigValue("username")
	if username == "" {
		if writer == nil {
			return nil
		}
		username = input.ReadFromUser(reader, writer, "Enter username: ")
	}

	password := p.ConfigManager.GetConfigValue("password")
	if password == "" {
		if writer == nil {
			return nil
		}
		password = input.ReadSensitiveFromUser(writer, "Enter password: ")
	}

	if writer != nil {
		spinner := loader.NewSpinner(writer)
		defer spinner.Stop()
	}

	if err := p.ConfigManager.SetConfigValue("username", username); err != nil {
		log.Error().Err(err).Msg("Could not set username")
		return err
	}

	if err := p.ConfigManager.SetConfigValue("password", password); err != nil {
		log.Error().Err(err).Msg("Could not set password")
		return err
	}

	if err := p.PlatformClient.Login(username, password); err != nil {
		log.Error().Err(err).Msg("Failed to authenticate")
		return err
	}

	if err := p.ConfigManager.SetConfigValue("token", p.PlatformClient.GetToken()); err != nil {
		log.Error().Err(err).Msg("Could not set userID")
		return err
	}

	if err := p.ConfigManager.SetConfigValue("userId", fmt.Sprint(p.PlatformClient.ActiveUserID())); err != nil {
		log.Error().Err(err).Msg("Could not set userID")
		return err
	}

	return nil
}
