/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/rs/zerolog/log"
	"os"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the Pixo Platform CLI",
	Long:  `Initialize the Pixo Platform CLI by setting up a default configuration file at ~/.pixo/config.yaml`,
	Run: func(cmd *cobra.Command, args []string) {

		initLogger(cmd)

		if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
			if err = os.Mkdir(cfgDir, 0755); err != nil {
				log.Error().Err(err).Msg("unable to create config directory")
				return
			}

			cmd.Println("created config directory")
		}

		if _, err := os.Stat(globalConfigFile); os.IsNotExist(err) {
			f, err := os.Create(globalConfigFile)
			if err != nil {
				log.Error().
					Err(err).
					Str("cfgFile", cfgFile).
					Msg("unable to create config file")
				return
			}

			log.Info().
				Str("cfgFile", cfgFile).
				Msg("created config file")

			defer f.Close()
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
