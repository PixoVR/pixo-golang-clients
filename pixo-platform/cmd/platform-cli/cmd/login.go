/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"errors"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/loader"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	primary_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/kyokomi/emoji"

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

			spinner := loader.NewSpinner(cmd.OutOrStdout())

			client, err := platform.NewClientWithBasicAuth(
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

			spinner.Stop()
			cmd.Println(emoji.Sprintf(":rocket:Login successful. Here is your API token: \n%s", client.GetToken()))

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

func getAuthenticatedClient() *platform.GraphQLAPIClient {
	token := input.GetConfigValue("token", "SECRET_KEY")
	if token != "" {
		oldAPIClient := primary_api.NewClient(urlfinder.ClientConfig{
			Token:     token,
			Lifecycle: input.GetConfigValue("lifecycle", "PIXO_LIFECYCLE"),
			Region:    input.GetConfigValue("region", "PIXO_REGION"),
		})
		return platform.NewClient(urlfinder.ClientConfig{
			Token:     oldAPIClient.GetToken(),
			Lifecycle: input.GetConfigValue("lifecycle", "PIXO_LIFECYCLE"),
			Region:    input.GetConfigValue("region", "PIXO_REGION"),
		})
	}

	apiKey := input.GetConfigValue("api-key", "PIXO_API_KEY")
	if apiKey != "" {
		apiClient.SetAPIKey(apiKey)
		return apiClient
	}

	username := input.GetConfigValue("username", "PIXO_USERNAME")
	password := input.GetConfigValue("password", "PIXO_PASSWORD")

	clientConfig := urlfinder.ClientConfig{
		Lifecycle: input.GetConfigValue("lifecycle", "PIXO_LIFECYCLE"),
		Region:    input.GetConfigValue("region", "PIXO_REGION"),
	}

	oldAPIClient, err := primary_api.NewClientWithBasicAuth(username, password, clientConfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to authenticate")
		return nil
	}

	apiClient = platform.NewClient(urlfinder.ClientConfig{
		Token:     oldAPIClient.GetToken(),
		Lifecycle: input.GetConfigValue("lifecycle", "PIXO_LIFECYCLE"),
		Region:    input.GetConfigValue("region", "PIXO_REGION"),
	})

	//if err = apiClient.Login(username, password); err != nil {
	//	log.Error().Err(err).Msg("Failed to authenticate")
	//	return nil
	//}

	return apiClient
}
