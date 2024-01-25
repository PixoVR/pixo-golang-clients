/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	"github.com/spf13/cobra"
)

// apiKeyCmd represents the apiKey command
var deleteApiKeyCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deleting API keys",
	Long:  `Delete API key with the following command:`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.Login(
			input.GetConfigValue("username", "PIXO_USERNAME"),
			input.GetConfigValue("password", "PIXO_PASSWORD"),
		); err != nil {
			return err
		}

		apiKeyID := input.GetIntValueOrAskUser(cmd, "key-id", "PIXO_API_KEY_ID")

		if err := apiClient.DeleteAPIKey(cmd.Context(), apiKeyID); err != nil {
			cmd.Println("Error creating API key: ", err.Error())
			return err
		}

		cmd.Printf("Deleted API key: %d\n", apiKeyID)
		return nil
	},
}

func init() {
	apiKeyCmd.AddCommand(deleteApiKeyCmd)

	deleteApiKeyCmd.Flags().IntP("key-id", "k", 0, "API key ID")
}
