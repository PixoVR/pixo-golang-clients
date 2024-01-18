/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"errors"
	"fmt"
	graphql_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"

	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the login rootCmd
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the Pixo Platform",
	Long: `Your username and password can be provided in multiple ways:
	- rootCmd line flags --username and --password
	- local config file ./config.yaml
	- environment variables PIXO_USERNAME and PIXO_PASSWORD
	- global config file ~/.pixo/config.yaml
	Will prioritize in order of the above list, and will prompt the user if none is found.	
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		token := input.GetConfigValue("secret-key", "SECRET_KEY")
		if token != "" {
			log.Debug().Msgf("Found secret key in config: %s", token)
			viper.Set("token", token)

		} else {
			username := input.GetStringValueOrAskUser(cmd, "username", config.PixoUsernameEnvVarKey)
			viper.Set("username", username)
			log.Debug().Msgf("Attempting to login as user: %s", username)

			password := input.GetSensitiveStringValueOrAskUser(cmd, "password", config.PixoPasswordEnvVarKey)
			viper.Set("password", password)
			log.Debug().Msgf("Attempting to login with password: %s", password)

			clientConfig := urlfinder.ClientConfig{
				Lifecycle: input.GetConfigValue("lifecycle", "PIXO_LIFECYCLE"),
				Region:    input.GetConfigValue("region", "PIXO_REGION"),
			}
			client, err := graphql_api.NewClientWithBasicAuth(
				username,
				password,
				clientConfig,
			)
			if err != nil {
				log.Error().Err(err).Msg("Could not create platform client")
				return err
			} else if client == nil {
				log.Error().Msg("Could not create platform client")
				return errors.New("could not create platform client")
			}

			cmd.Println(fmt.Sprintf("Login successful. Here is your API token: \n%s", client.GetToken()))

			viper.Set("token", client.GetToken())
		}

		if err := viper.WriteConfigAs(cfgFile); err != nil {
			log.Error().Err(err).Msg("Could not write config file")
			return err
		}

		return nil
	},
}

func init() {
	authCmd.AddCommand(loginCmd)
}
