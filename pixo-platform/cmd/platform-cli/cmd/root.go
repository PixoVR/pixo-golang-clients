/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/parser"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	platformAPI "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
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

	apiClient *platformAPI.GraphQLAPIClient

	cfgFile string
)

func NewRootCmd() *cobra.Command {
	return rootCmd
}

// rootCmd represents the base rootCmd when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "pixo",
	Version: "0.0.146",
	Short:   "A CLI for the Pixo Platform",
	Long:    `A CLI tool used to streamline interactions with the Pixo Platform`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: rootCmd.OutOrStdout(), TimeFormat: time.RFC1123})

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug logging")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is %s)", globalConfigFile))
	rootCmd.PersistentFlags().StringP("ini", "c", parser.DefaultConfigFilepath, "Path to the ini file to use for the rootCmd")
	rootCmd.PersistentFlags().StringP("module-id", "m", "", "Module ID to use for the rootCmd")

	if cfgFile == "" {
		cfgFile = globalConfigFile
	}
	viper.AddConfigPath(cfgDir)
	viper.AddConfigPath(".pixo")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Error().Err(err).Msg("Unable to read config file")
	}

	clientConfig := urlfinder.ClientConfig{
		Token:     viper.GetString("token"),
		Lifecycle: input.GetConfigValue("lifecycle", "PIXO_LIFECYCLE"),
		Region:    input.GetConfigValue("region", "PIXO_REGION"),
	}
	apiClient = platformAPI.NewClient(clientConfig)

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		initLogger(cmd)
	}
}

func initLogger(cmd *cobra.Command) {

	if isDebug(cmd) {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Debug logging enabled")
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

}

func isDebug(cmd *cobra.Command) bool {
	return cmd.Flag("debug").Value.String() == "true"
}
