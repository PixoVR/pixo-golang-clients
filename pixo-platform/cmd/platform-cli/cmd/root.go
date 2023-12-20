/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/parser"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	platformAPI "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	cc "github.com/ivanpirog/coloredcobra"
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

func RootCmd() *cobra.Command {
	return rootCmd
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "pixo",
	Version: "0.0.96",
	Short:   "A CLI for the Pixo Platform",
	Long:    `A CLI tool used to simplify interactions with the Pixo Platform`,
}

func Execute() {

	cc.Init(&cc.Config{
		RootCmd:         rootCmd,
		Headings:        cc.HiYellow + cc.Bold + cc.Underline,
		ExecName:        cc.HiRed + cc.Bold,
		Commands:        cc.HiRed + cc.Bold,
		CmdShortDescr:   cc.White + cc.Italic,
		Example:         cc.Italic,
		Flags:           cc.HiRed + cc.Bold,
		FlagsDataType:   cc.HiCyan,
		FlagsDescr:      cc.White + cc.Italic,
		NoExtraNewlines: true,
	})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
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
