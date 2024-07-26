package ctx

import (
	"context"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/headset"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/legacy"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/editor"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms/fancy"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/loader"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/printer"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/spf13/cobra"
	"os"
)

type Context struct {
	ConfigManager     *config.ConfigManager
	FileManager       *config.FileManager
	FormHandler       forms.FormHandler
	Printer           printer.Printer
	FileOpener        editor.FileOpener
	HeadsetClient     headset.Client
	PlatformClient    platform.Client
	LegacyClient      legacy.LegacyClient
	MatchmakingClient matchmaker.Matchmaker
}

func NewContext(configFiles ...string) *Context {
	var configFile string
	if len(configFiles) > 0 {
		for _, file := range configFiles {
			if _, err := os.Stat(file); err == nil {
				configFile = file
				break
			}
		}
	}

	formHandler := fancy.NewFormHandler()
	//formHandler := basic.NewFormHandler(os.Stdin, os.Stdout)

	fileManager := config.NewFileConfigManager(configFile, formHandler)
	configManager := config.NewConfigManager(fileManager, formHandler)

	token, _ := configManager.GetConfigValue("token")

	clientConfig := urlfinder.ClientConfig{
		Token:     token,
		Lifecycle: configManager.Lifecycle(),
		Region:    configManager.Region(),
	}

	return &Context{
		FormHandler:       formHandler,
		ConfigManager:     configManager,
		FileManager:       fileManager,
		Printer:           printer.NewEmojiPrinter(os.Stdout),
		HeadsetClient:     headset.NewClient(clientConfig),
		PlatformClient:    platform.NewClient(clientConfig),
		MatchmakingClient: matchmaker.NewClient(clientConfig),
		LegacyClient:      legacy.NewClient(clientConfig),
		FileOpener:        editor.NewFileOpener(""),
	}
}

func (p *Context) SetIO(cmd *cobra.Command) {
	p.Printer.SetWriter(cmd.OutOrStdout())
	p.FormHandler.SetReader(cmd.InOrStdin())
	p.FormHandler.SetWriter(cmd.OutOrStdout())
}

func (p *Context) Authenticate(cmd *cobra.Command) error {
	token, ok := p.ConfigManager.GetFlagOrConfigValue("token", cmd)
	if ok && token != "" {
		p.PlatformClient.SetToken(token)
		p.HeadsetClient.SetToken(token)
		p.ConfigManager.SetConfigValue("token", token)
		return nil
	}

	apiKey, ok := p.ConfigManager.GetFlagOrConfigValue("api-key", cmd)
	if ok && apiKey != "" {
		p.PlatformClient.SetAPIKey(apiKey)
		p.ConfigManager.SetConfigValue("api-key", apiKey)
		return nil
	}

	username, ok := p.ConfigManager.GetFlagOrConfigValue("username", cmd)
	if ok && username != "" {
		p.ConfigManager.SetConfigValue("username", username)
	}

	password, ok := p.ConfigManager.GetFlagOrConfigValue("password", cmd)
	if ok && password != "" {
		p.ConfigManager.SetConfigValue("password", password)
	}

	ctx := context.Background()
	if cmd != nil {
		ctx = cmd.Context()
	}
	spinner := loader.NewLoader(ctx, "Logging into the Pixo Platform...", p.Printer)
	defer spinner.Stop()

	if username != "" && password != "" {
		if err := p.PlatformClient.Login(username, password); err != nil {
			return err
		}
		p.ConfigManager.SetConfigValue("token", p.PlatformClient.GetToken())
		p.ConfigManager.SetIntConfigValue("user-id", p.PlatformClient.ActiveUserID())
		p.ConfigManager.SetIntConfigValue("org", p.PlatformClient.ActiveOrgID())
	}

	return nil
}
