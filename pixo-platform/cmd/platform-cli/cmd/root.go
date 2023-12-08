/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/parser"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	platformAPI "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "pixo",
	Version: "0.0.93",
	Short:   "A CLI for the Pixo Platform",
	Long:    `A CLI tool used to simplify interactions with the Pixo Platform`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initLogger(cmd *cobra.Command) {

	if cmd.Flag("debug").Value.String() == "true" {
		log.Info().Msg("Debug logging enabled")
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

}

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: rootCmd.OutOrStdout(), TimeFormat: time.RFC1123})

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug logging")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is %s)", globalConfigFile))
	rootCmd.PersistentFlags().StringP("ini", "c", parser.DefaultConfigFilepath, "Path to the ini file to use for the command")
	rootCmd.PersistentFlags().StringP("module-id", "m", "", "Module ID to use for the command")

	if cfgFile == "" {
		cfgFile = globalConfigFile
	}
	viper.AddConfigPath(cfgDir)
	viper.AddConfigPath(".pixo")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Error().Err(err).Msg("Failed to read config file")
	}

	apiClient = platformAPI.NewClient(
		viper.GetString("token"),
		input.GetConfigValue("lifecycle", "PIXO_LIFECYCLE"),
		input.GetConfigValue("region", "PIXO_REGION"),
	)
}
