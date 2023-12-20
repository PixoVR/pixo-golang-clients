/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// initCmd represents the init rootCmd
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the Pixo Platform CLI",
	Long:  `Initialize the Pixo Platform CLI by setting up a default configuration file at ~/.pixo/config.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
			if err = os.Mkdir(cfgDir, 0755); err != nil {
				return err
			}

			cmd.Println("created config directory")
		}

		if _, err := os.Stat(globalConfigFile); os.IsNotExist(err) {
			f, err := os.Create(globalConfigFile)
			if err != nil {
				return err
			}

			defer f.Close()
		}

		cmd.Println("Welcome to the Pixo Platform CLI! Your config file is located at ~/.pixo/config.yaml. Please run `pixo auth login` to get started.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
