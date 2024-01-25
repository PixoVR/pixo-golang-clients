/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/spf13/cobra"
)

// apiKeyCmd represents the apiKey command
var createApiKeyCmd = &cobra.Command{
	Use:   "create",
	Short: "Creating API keys",
	Long:  `Create API key with the following command:`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.Login(
			input.GetConfigValue("username", "PIXO_USERNAME"),
			input.GetConfigValue("password", "PIXO_PASSWORD"),
		); err != nil {
			return err
		}

		input := platform.APIKey{
			//UserID: input.GetIntValue(cmd, "user-id", "PIXO_USER_ID"),
		}

		apiKey, err := apiClient.CreateAPIKey(cmd.Context(), input)
		if err != nil {
			cmd.Println("Error creating API key: ", err.Error())
			return err
		}

		cmd.Println("Created API key: " + apiKey.Key)
		return nil
	},
}

func init() {
	apiKeyCmd.AddCommand(createApiKeyCmd)
}
