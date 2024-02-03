/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/loader"
	"github.com/spf13/cobra"
)

// apiKeyCmd represents the apiKey command
var deleteApiKeyCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deleting an API key",
	Long:  `Delete API key with the following command:`,
	RunE: func(cmd *cobra.Command, args []string) error {
		spinner := loader.NewSpinner(cmd.OutOrStdout())

		apiKeyID, ok := Ctx.ConfigManager.GetIntFlagOrConfigValue("key-id", cmd)
		if !ok {
			cmd.Println("Error: API key ID is required")
			return nil
		}

		if err := Ctx.PlatformClient.DeleteAPIKey(cmd.Context(), apiKeyID); err != nil {
			cmd.Println("Error creating API key: ", err.Error())
			return err
		}

		spinner.Stop()
		cmd.Printf("Deleted API key: %d\n", apiKeyID)
		return nil
	},
}

func init() {
	apiKeyCmd.AddCommand(deleteApiKeyCmd)

	deleteApiKeyCmd.Flags().IntP("key-id", "k", 0, "API key ID")
}
