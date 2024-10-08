package ctx

import (
	"context"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/headset"
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
	printer.Printer
	ConfigManager     *config.ConfigManager
	FileManager       *config.FileManager
	FormHandler       forms.FormHandler
	FileOpener        editor.FileOpener
	HeadsetClient     headset.Client
	PlatformClient    platform.Client
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
	emojiPrinter := printer.NewEmojiPrinter(os.Stdout)
	configManager := config.NewConfigManager(fileManager, emojiPrinter, formHandler)

	token, _ := configManager.GetConfigValue("token")

	clientConfig := urlfinder.ClientConfig{
		Token:     token,
		Lifecycle: configManager.Lifecycle(),
		Region:    configManager.Region(),
	}

	return &Context{
		Printer:           emojiPrinter,
		FormHandler:       formHandler,
		ConfigManager:     configManager,
		FileManager:       fileManager,
		HeadsetClient:     headset.NewClient(clientConfig),
		PlatformClient:    platform.NewClient(clientConfig),
		MatchmakingClient: matchmaker.NewClient(clientConfig),
		FileOpener:        editor.NewFileOpener(""),
	}
}

func (p *Context) SetIO(cmd *cobra.Command) {
	p.Printer.SetWriter(cmd.OutOrStdout())
	p.FormHandler.SetReader(cmd.InOrStdin())
	p.FormHandler.SetWriter(cmd.OutOrStdout())
}

func (p *Context) Authenticate(cmd *cobra.Command) error {
	ctx := context.Background()
	if cmd != nil {
		ctx = cmd.Context()
	}

	token, ok := p.ConfigManager.GetFlagOrConfigValue("auth-token", cmd)
	if ok {
		p.PlatformClient.SetToken(token)
		if user, err := p.PlatformClient.CheckAuth(ctx); err == nil {
			p.MatchmakingClient.SetToken(token)
			p.HeadsetClient.SetToken(token)
			p.ConfigManager.SetConfigValue("auth-token", token)
			p.ConfigManager.SetIntConfigValue("auth-user-id", user.ID)
			return nil
		} else {
			p.PlatformClient.SetToken("")
		}
	}

	apiKey, ok := p.ConfigManager.GetFlagOrConfigValue("api-key", cmd)
	if ok {
		p.PlatformClient.SetAPIKey(apiKey)
		if user, err := p.PlatformClient.CheckAuth(ctx); err == nil {
			p.MatchmakingClient.SetAPIKey(apiKey)
			p.HeadsetClient.SetAPIKey(apiKey)
			p.ConfigManager.SetConfigValue("api-key", apiKey)
			p.ConfigManager.SetIntConfigValue("auth-user-id", user.ID)
			return nil
		}
	}

	username, ok := p.ConfigManager.GetFlagOrConfigValue("auth-username", cmd)
	if ok {
		p.ConfigManager.SetConfigValue("auth-username", username)
	}

	password, ok := p.ConfigManager.GetFlagOrConfigValue("auth-password", cmd)
	if ok {
		p.ConfigManager.SetConfigValue("auth-password", password)
	}

	if cmd != nil {
		spinner := loader.NewLoader(ctx, "Logging into the Pixo Platform...", cmd.OutOrStdout())
		defer spinner.Stop()
	}

	if username != "" && password != "" {
		if err := p.PlatformClient.Login(username, password); err != nil {
			return err
		}
		p.ConfigManager.SetConfigValue("auth-token", p.PlatformClient.GetToken())
		p.ConfigManager.SetIntConfigValue("auth-user-id", p.PlatformClient.ActiveUserID())
	}

	return nil
}
