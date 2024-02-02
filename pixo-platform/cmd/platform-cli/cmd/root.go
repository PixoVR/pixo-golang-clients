/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/clients"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/editor"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	primary_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	homeDir          = os.Getenv("HOME")
	cfgDir           = fmt.Sprintf("%s/.pixo", homeDir)
	globalConfigFile = fmt.Sprintf("%s/config.yaml", cfgDir)
	isDebug          bool

	PlatformCtx *clients.PlatformContext

	cfgFile string
)

func NewRootCmd(platformContext *clients.PlatformContext) *cobra.Command {
	PlatformCtx = platformContext
	return rootCmd
}

// rootCmd represents the base rootCmd when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "pixo",
	Version: "0.0.148",
	Short:   "A CLI for the Pixo Platform",
	Long:    `A CLI tool used to streamline interactions with the Pixo Platform`,
}

func Execute() {
	configManager := config.NewFileManager(cfgDir)
	if err := configManager.ReadConfigFile(cfgFile); err != nil {
		log.Error().Err(err).Msg("Could not read config file")
	}

	token, _ := configManager.GetConfigValue("token")

	clientConfig := urlfinder.ClientConfig{
		Token:     token,
		Lifecycle: configManager.Lifecycle(),
		Region:    configManager.Region(),
	}

	PlatformCtx = &clients.PlatformContext{
		ConfigManager:     configManager,
		OldAPIClient:      primary_api.NewClient(clientConfig),
		PlatformClient:    platform.NewClient(clientConfig),
		MatchmakingClient: matchmaker.NewMatchmaker(clientConfig),
		FileOpener:        editor.NewFileOpener(""),
	}

	if err := PlatformCtx.Authenticate(nil, nil); err != nil {
		log.Error().Err(err).Msg("Failed to authenticate")
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: rootCmd.OutOrStdout(), TimeFormat: time.RFC1123})

	rootCmd.PersistentFlags().BoolVarP(&isDebug, "debug", "d", false, "Enable debug logging")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is %s)", globalConfigFile))
	rootCmd.PersistentFlags().StringP("module-id", "m", "", "Module ID to use for the rootCmd")

	if cfgFile == "" {
		cfgFile = globalConfigFile
	}

	viper.AddConfigPath(".pixo")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		initLogger()
	}
}

func initLogger() {

	if isDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Debug logging enabled")
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

}
