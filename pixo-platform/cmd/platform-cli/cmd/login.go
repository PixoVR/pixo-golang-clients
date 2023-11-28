/*
Copyright © 2023 NAME HERE walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the Pixo Platform",
	Long: `Your username and password can be provided in multiple ways:
	- command line flags --username and --password
	- local config file ./config.yaml
	- environment variables PIXO_USERNAME and PIXO_PASSWORD
	- global config file ~/.pixo/config.yaml
	Will prioritize in order of the above list, and will prompt the user if none is found.	
`,
	Run: func(cmd *cobra.Command, args []string) {
		var client *platform.PrimaryAPIClient

		username := input.GetStringValue(cmd, "username", "PIXO_USERNAME")
		log.Debug().Msgf("Attempting to login as user: %s", username)

		password := input.GetStringValue(cmd, "password", "PIXO_PASSWORD")
		log.Debug().Msgf("Attempting to login with password: %s", password)

		if client = platform.NewClientWithBasicAuth(username, password, ""); client == nil {
			return
		}

		log.Info().Msgf("Login successful. Here is your API token: %s", client.GetToken())

		viper.Set("username", username)
		viper.Set("password", password)
		viper.Set("token", client.GetToken())

		if err := viper.WriteConfigAs(cfgFile); err != nil {
			log.Error().Err(err).Msg("Could not write config file")
		}

	},
}

func init() {
	authCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
