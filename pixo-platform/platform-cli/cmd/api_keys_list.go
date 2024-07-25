/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	graphql_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/loader"
	"github.com/spf13/cobra"
)

// apiKeyCmd represents the apiKey command
var listApiKeyCmd = &cobra.Command{
	Use:   "list",
	Short: "List API keys",
	Long:  `List API key with the following command:`,
	RunE: func(cmd *cobra.Command, args []string) error {
		spinner := loader.NewLoader(cmd.Context(), ":key: Getting API Keys...", Ctx.Printer)

		apiKeyParams := &graphql_api.APIKeyQueryParams{}

		userIDVal, _ := Ctx.ConfigManager.GetIntConfigValue("user-id")
		if userIDVal != 0 {
			apiKeyParams.UserID = &userIDVal
		}

		apiKeys, err := Ctx.PlatformClient.GetAPIKeys(cmd.Context(), apiKeyParams)
		spinner.Stop()
		if err != nil {
			Ctx.Printer.Println("Error getting API keys: ", err)
			return err
		}

		if len(apiKeys) == 0 {
			Ctx.Printer.Println("No API keys found")
			return nil
		}

		Ctx.Printer.Println("API keys:")
		for _, apiKey := range apiKeys {
			Ctx.Printer.Println("Key ID: ", apiKey.ID)
		}

		return nil
	},
}

func init() {
	apiKeyCmd.AddCommand(listApiKeyCmd)
}
