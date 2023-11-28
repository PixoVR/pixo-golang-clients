/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/parser"
	platformAPI "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	secretKey = config.GetEnvOrReturn("SECRET_KEY", "fake-key")
	apiClient = platformAPI.NewClient(secretKey, "")
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pixo-cli",
	Short: "A CLI for Pixo Platform",
	Long: `A CLI for Pixo Platform for managing platform resources
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	log.Logger = log.Level(zerolog.InfoLevel)

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.platform-cli.yaml)")

	deployCmd.PersistentFlags().StringP("ini", "c", parser.DefaultConfigFilepath, "Path to the ini file to use for the command")
	rootCmd.PersistentFlags().StringP("module-id", "m", "", "Module ID to use for the command")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
