/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"fmt"

	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/config"
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
		token := input.GetConfigValue("secret-key", "SECRET_KEY")
		if token != "" {
			log.Debug().Msgf("Found secret key in config: %s", token)
			viper.Set("token", token)
		} else {

			username := input.GetStringValueOrAskUser(cmd, "username", config.PixoUsernameEnvVarKey)
			viper.Set("username", username)
			log.Debug().Msgf("Attempting to login as user: %s", username)

			password := input.GetStringValueOrAskUser(cmd, "password", config.PixoPasswordEnvVarKey)
			viper.Set("password", password)
			log.Debug().Msgf("Attempting to login with password: %s", password)

			var client *platform.PrimaryAPIClient
			if client = platform.NewClientWithBasicAuth(username, password, input.GetConfigValue("legacy-api-url", config.PixoLegacyAPIURLEnvVarKey)); client == nil {
				return
			}

			cmd.Println(fmt.Sprintf("Login successful. Here is your API token: \n%s", client.GetToken()))

			viper.Set("token", client.GetToken())
		}

		if err := viper.WriteConfigAs(cfgFile); err != nil {
			log.Error().Err(err).Msg("Could not write config file")
		}
	},
}

func init() {
	authCmd.AddCommand(loginCmd)
}
