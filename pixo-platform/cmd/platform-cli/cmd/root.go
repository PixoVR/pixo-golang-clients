/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/parser"
	platformAPI "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
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
	secretKey        = config.GetEnvOrReturn("SECRET_KEY", "fake-key")
	apiClient        = platformAPI.NewClient(secretKey, "")

	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pixo",
	Short: "A CLI for the Pixo Platform",
	Long: `A CLI tool used to simplify interactions with the Pixo Platform 
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug logging")
	if viper.GetBool("debug") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is %s)", globalConfigFile))
	rootCmd.PersistentFlags().StringP("ini", "c", parser.DefaultConfigFilepath, "Path to the ini file to use for the command")
	rootCmd.PersistentFlags().StringP("org-id", "o", "", "Org ID to use for the command")
	rootCmd.PersistentFlags().StringP("module-id", "m", "", "Module ID to use for the command")

	if cfgFile == "" {
		cfgFile = globalConfigFile
	}
	viper.AddConfigPath(cfgDir)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Error().Err(err).Msg("Failed to read config file")
	}
}
