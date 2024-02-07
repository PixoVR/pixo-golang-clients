/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/clients"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	homeDir          = os.Getenv("HOME")
	configDirName    = ".pixo"
	configFileName   = "config.yaml"
	globalCfgDir     = fmt.Sprintf("%s/%s", homeDir, configDirName)
	globalConfigFile = fmt.Sprintf("%s/%s", globalCfgDir, configFileName)
	localConfigFile  = fmt.Sprintf("./%s/%s", configDirName, configFileName)

	isDebug          bool
	cfgFileFlagInput string
	Ctx              *clients.CLIContext
)

func GetRootCmd() *cobra.Command {
	return rootCmd
}

var rootCmd = &cobra.Command{
	Use:     "pixo",
	Version: "0.0.161",
	Short:   "A CLI for the Pixo Platform",
	Long:    `A CLI tool used to streamline interactions with the Pixo Platform`,
}

func Execute() {

	Ctx = clients.NewCLIContextWithConfig(localConfigFile, globalConfigFile)

	if err := Ctx.Authenticate(nil); err != nil {
		log.Error().Err(err).Msg("Failed to authenticate")
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: rootCmd.OutOrStdout(), TimeFormat: time.RFC1123})

	rootCmd.PersistentFlags().StringVar(&cfgFileFlagInput, "config", "", fmt.Sprintf("config file (default is %s)", globalConfigFile))
	rootCmd.PersistentFlags().BoolVarP(&isDebug, "debug", "d", false, "Enable debug logging")

	configCmd.PersistentFlags().StringP("lifecycle", "l", "", "Lifecycle of Pixo Platform to use (dev, stage, prod)")
	configCmd.PersistentFlags().StringP("region", "r", "", "Region of Pixo Platform to use (na, saudi)")
	rootCmd.PersistentFlags().IntP("module-id", "m", 0, "Module ID")

	if cfgFileFlagInput == "" {
		cfgFileFlagInput = globalConfigFile
	}

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		initLogger()
		Ctx.SetIO(cmd)
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
