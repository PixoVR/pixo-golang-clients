package ctx

import (
	"context"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/headset"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/editor"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms/charm"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/loader"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/printer"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/spf13/cobra"
	"os"
)

type CLIContext struct {
	ConfigManager     *config.ConfigManager
	FileManager       *config.FileManager
	FormHandler       forms.FormHandler
	Printer           printer.Printer
	FileOpener        editor.FileOpener
	HeadsetClient     headset.Client
	PlatformClient    platform.Client
	MatchmakingClient matchmaker.Matchmaker
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
	//formHandler := basic.NewFormHandler(os.Stdin, os.Stdout)

	fileManager := config.NewFileConfigManager(configFile, formHandler)
	configManager := config.NewConfigManager(fileManager, formHandler)

	token, _ := configManager.GetConfigValue("token")

	clientConfig := urlfinder.ClientConfig{
		Token:     token,
		Lifecycle: configManager.Lifecycle(),
		Region:    configManager.Region(),
	}

	return &CLIContext{
		FormHandler:       formHandler,
		ConfigManager:     configManager,
		FileManager:       fileManager,
		Printer:           printer.NewEmojiPrinter(os.Stdout),
		HeadsetClient:     headset.NewClient(clientConfig),
		PlatformClient:    platform.NewClient(clientConfig),
		MatchmakingClient: matchmaker.NewClient(clientConfig),
		FileOpener:        editor.NewFileOpener(""),
	}
}

func (p *CLIContext) SetIO(cmd *cobra.Command) {
	p.Printer.SetWriter(cmd.OutOrStdout())
	p.FormHandler.SetReader(cmd.InOrStdin())
	p.FormHandler.SetWriter(cmd.OutOrStdout())
}

func (p *CLIContext) Authenticate(cmd *cobra.Command) error {
	if p.PlatformClient.IsAuthenticated() {
		return nil
	}

	token, ok := p.ConfigManager.GetFlagOrConfigValue("token", cmd)
	if ok {
		p.PlatformClient.SetToken(token)
		p.HeadsetClient.SetToken(token)
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
		p.Printer.Println(":exclamation: Login failed. Username is required.")
		return nil
	}
	p.ConfigManager.SetConfigValue("username", username)

	password, ok := p.ConfigManager.GetSensitiveFlagOrConfigValueOrAskUser("password", cmd)
	if !ok {
		p.Printer.Println(":exclamation: Login failed. Password is required.")
		return nil
	}
	p.ConfigManager.SetConfigValue("password", password)

	ctx := context.Background()
	if cmd != nil {
		ctx = cmd.Context()
	}
	spinner := loader.NewLoader(ctx, "Logging into the Pixo Platform...", p.Printer)
	defer spinner.Stop()

	if err := p.PlatformClient.Login(username, password); err != nil {
		p.Printer.Println(":exclamation: Login failed. Please check your credentials and try again.\nError: ", err)
		return err
	}

	p.ConfigManager.SetConfigValue("token", p.PlatformClient.GetToken())
	p.ConfigManager.SetIntConfigValue("user-id", p.PlatformClient.ActiveUserID())
	p.ConfigManager.SetIntConfigValue("org-id", p.PlatformClient.ActiveUserID())
	return nil
}
