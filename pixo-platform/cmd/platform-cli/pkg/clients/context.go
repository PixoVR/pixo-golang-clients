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
	if p.PlatformClient.IsAuthenticated() {
		return nil
	}

	token, ok := p.ConfigManager.GetConfigValue("token")
	if ok {
		p.PlatformClient.SetToken(token)
		return nil
	}

	apiKey, ok := p.ConfigManager.GetConfigValue("api-key")
	if ok {
		p.PlatformClient.SetAPIKey(apiKey)
		return nil
	}

	username, ok := p.ConfigManager.GetConfigValue("username")
	if !ok {
		if writer == nil {
			return nil
		}
		username = input.ReadFromUser(reader, writer, "username")
	}

	password, ok := p.ConfigManager.GetConfigValue("password")
	if !ok {
		if writer == nil {
			return nil
		}
		password = input.ReadSensitiveFromUser(writer, "password")
	}

	if writer != nil {
		spinner := loader.NewSpinner(writer)
		defer spinner.Stop()
	}

	p.ConfigManager.SetConfigValue("username", username)
	p.ConfigManager.SetConfigValue("password", password)

	if err := p.PlatformClient.Login(username, password); err != nil {
		log.Error().Err(err).Msg("Failed to authenticate")
		return err
	}

	p.ConfigManager.SetConfigValue("token", p.PlatformClient.GetToken())
	p.ConfigManager.SetConfigValue("user-id", fmt.Sprint(p.PlatformClient.ActiveUserID()))

	return nil
}
