/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/ctx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	cliVersion = "0.1.34"

	homeDir          = os.Getenv("HOME")
	configDirName    = ".pixo"
	configFileName   = "config.yaml"
	globalCfgDir     = fmt.Sprintf("%s/%s", homeDir, configDirName)
	globalConfigFile = fmt.Sprintf("%s/%s", globalCfgDir, configFileName)
	localConfigFile  = fmt.Sprintf("./%s/%s", configDirName, configFileName)

	activeConfigFile string

	isDebug          bool
	cfgFileFlagInput string
	Ctx              *ctx.Context
)

var rootCmd = &cobra.Command{
	Use:     "pixo",
	Version: cliVersion,
	Short:   "A CLI for the Pixo Platform",
	Long:    `A CLI tool used to streamline interactions with the Pixo Platform`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initLogger()
		Ctx.SetIO(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		//if _, err := tea.NewProgram(tui.NewModel()).Run(); err != nil {
		//	return err
		//}
		return nil
	},
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}

func Execute() {
	Ctx = ctx.NewContext(localConfigFile, globalConfigFile)
	_ = Ctx.Authenticate(nil)

	activeConfigFile = Ctx.FileManager.ConfigFile()

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

	if cfgFileFlagInput == "" {
		cfgFileFlagInput = globalConfigFile
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
