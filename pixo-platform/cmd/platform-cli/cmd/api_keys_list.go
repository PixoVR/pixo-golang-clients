/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	graphql_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/spf13/cobra"
)

// apiKeyCmd represents the apiKey command
var listApiKeyCmd = &cobra.Command{
	Use:   "list",
	Short: "List API keys",
	Long:  `List API key with the following command:`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.Login(
			input.GetConfigValue("username", "PIXO_USERNAME"),
			input.GetConfigValue("password", "PIXO_PASSWORD"),
		); err != nil {
			return err
		}

		apiKeyParams := &graphql_api.APIKeyQueryParams{}
		apiKeys, err := apiClient.GetAPIKeys(cmd.Context(), apiKeyParams)
		if err != nil {
			cmd.Println("Error listing API keys: ", err.Error())
			return err
		}

		cmd.Println("API Keys:")
		for _, apiKey := range apiKeys {
			cmd.Printf("Key ID: %d\n", apiKey.ID)
		}

		return nil
	},
}

func init() {
	apiKeyCmd.AddCommand(listApiKeyCmd)
}
