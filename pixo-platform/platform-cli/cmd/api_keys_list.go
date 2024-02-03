/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	graphql_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/loader"
	"github.com/spf13/cobra"
)

// apiKeyCmd represents the apiKey command
var listApiKeyCmd = &cobra.Command{
	Use:   "list",
	Short: "List API keys",
	Long:  `List API key with the following command:`,
	RunE: func(cmd *cobra.Command, args []string) error {
		spinner := loader.NewSpinner(cmd.OutOrStdout())

		apiKeyParams := &graphql_api.APIKeyQueryParams{}

		userIDVal, _ := Ctx.ConfigManager.GetIntConfigValue("user-id")
		if userIDVal != 0 {
			apiKeyParams.UserID = &userIDVal
		}

		apiKeys, err := Ctx.PlatformClient.GetAPIKeys(cmd.Context(), apiKeyParams)
		if err != nil {
			cmd.Println("Error listing API keys: ", err.Error())
			return err
		}

		spinner.Stop()

		if len(apiKeys) == 0 {
			cmd.Println("No API Keys found")
			return nil
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
