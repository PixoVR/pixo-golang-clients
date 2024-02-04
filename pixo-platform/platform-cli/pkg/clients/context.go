package clients

import (
	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/editor"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/loader"
	primary_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
	"os"
)

type CLIContext struct {
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

	configManager := config.NewFileManager(configFile)

	token, _ := configManager.GetConfigValue("token")

	clientConfig := urlfinder.ClientConfig{
		Token:     token,
		Lifecycle: configManager.Lifecycle(),
		Region:    configManager.Region(),
	}

	return &CLIContext{
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
		return nil
	}

	apiKey, ok := p.ConfigManager.GetFlagOrConfigValue("api-key", cmd)
	if ok {
		p.PlatformClient.SetAPIKey(apiKey)
		return nil
	}

	username, ok := p.ConfigManager.GetFlagOrConfigValueOrAskUser("username", cmd)
	if !ok {
		return nil
	}

	password, ok := p.ConfigManager.GetSensitiveFlagOrConfigValueOrAskUser("password", cmd)
	if !ok {
		return nil
	}

	if p.ConfigManager.Writer() != nil {
		spinner := loader.NewSpinner(p.ConfigManager.Writer())
		defer spinner.Stop()
	}

	p.ConfigManager.SetConfigValue("username", username)
	p.ConfigManager.SetConfigValue("password", password)

	if err := p.PlatformClient.Login(username, password); err != nil {
		cmd.Println(emoji.Sprintf(":exclamation: Login failed. Please check your credentials and try again.\n%s", err))
		return err
	}

	p.ConfigManager.SetConfigValue("token", p.PlatformClient.GetToken())
	p.ConfigManager.SetIntConfigValue("user-id", p.PlatformClient.ActiveUserID())

	return nil
}