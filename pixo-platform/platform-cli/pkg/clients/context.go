package clients

import (
	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/editor"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/forms"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/forms/charm"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/loader"
	primary_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/spf13/cobra"
	"os"
)

type CLIContext struct {
	FormHandler       forms.FormHandler
	ConfigManager     config.Manager
	OldAPIClient      abstract_client.AbstractClient
	PlatformClient    platform.PlatformClient
	MatchmakingClient matchmaker.Matchmaker
	FileOpener        editor.FileOpener
}

func NewCLIContextWithConfig(configFiles ...string) *CLIContext {
	var configFile string
	if len(configFiles) > 0 {
		for _, file := range configFiles {
			if _, err := os.Stat(file); err == nil {
				configFile = file
				break
			}
		}
	}

	formHandler := charm.NewFormHandler()

	configManager := config.NewFileManager(configFile, formHandler)

	token, _ := configManager.GetConfigValue("token")

	clientConfig := urlfinder.ClientConfig{
		Token:     token,
		Lifecycle: configManager.Lifecycle(),
		Region:    configManager.Region(),
	}

	return &CLIContext{
		FormHandler:       formHandler,
		ConfigManager:     configManager,
		OldAPIClient:      primary_api.NewClient(clientConfig),
		PlatformClient:    platform.NewClient(clientConfig),
		MatchmakingClient: matchmaker.NewMatchmaker(clientConfig),
		FileOpener:        editor.NewFileOpener(""),
	}
}

func (p *CLIContext) SetIO(cmd *cobra.Command) {
	p.ConfigManager.SetReader(cmd.InOrStdin())
	p.ConfigManager.SetWriter(cmd.OutOrStdout())
}

func (p *CLIContext) Authenticate(cmd *cobra.Command) error {
	if p.PlatformClient.IsAuthenticated() {
		return nil
	}

	token, ok := p.ConfigManager.GetFlagOrConfigValue("token", cmd)
	if ok {
		p.PlatformClient.SetToken(token)
		p.ConfigManager.SetConfigValue("token", token)
		return nil
	}

	apiKey, ok := p.ConfigManager.GetFlagOrConfigValue("api-key", cmd)
	if ok {
		p.PlatformClient.SetAPIKey(apiKey)
		p.ConfigManager.SetConfigValue("api-key", apiKey)
		return nil
	}

	username, ok := p.ConfigManager.GetFlagOrConfigValueOrAskUser("username", cmd)
	if !ok {
		p.ConfigManager.Println(":exclamation: Login failed. Username is required.")
		return nil
	}
	p.ConfigManager.SetConfigValue("username", username)

	password, ok := p.ConfigManager.GetSensitiveFlagOrConfigValueOrAskUser("password", cmd)
	if !ok {
		p.ConfigManager.Println(":exclamation: Login failed. Password is required.")
		return nil
	}
	p.ConfigManager.SetConfigValue("password", password)

	spinner := loader.NewLoader(cmd.Context(), "Logging into the Pixo Platform...", p.ConfigManager)
	defer spinner.Stop()

	if err := p.PlatformClient.Login(username, password); err != nil {
		p.ConfigManager.Println(":exclamation: Login failed. Please check your credentials and try again.\nError: ", err)
		return err
	}

	p.ConfigManager.SetConfigValue("token", p.PlatformClient.GetToken())
	p.ConfigManager.SetIntConfigValue("user-id", p.PlatformClient.ActiveUserID())

	return nil
}
