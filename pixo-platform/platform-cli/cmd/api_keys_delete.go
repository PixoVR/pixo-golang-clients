/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/loader"
	"github.com/spf13/cobra"
)

// apiKeyCmd represents the apiKey command
var deleteApiKeyCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deleting an API key",
	Long:  `Delete API key with the following command:`,
	RunE: func(cmd *cobra.Command, args []string) error {
		spinner := loader.NewLoader(cmd.Context(), ":key: Deleting API Key...", Ctx.Printer)
		defer spinner.Stop()

		apiKeyID, ok := Ctx.ConfigManager.GetIntFlagOrConfigValue("key-id", cmd)
		if !ok {
			Ctx.Printer.Println("Error: API key ID is required")
			return nil
		}

		if err := Ctx.PlatformClient.DeleteAPIKey(cmd.Context(), apiKeyID); err != nil {
			Ctx.Printer.Println("Error deleting API key: ", err)
			return err
		}

		Ctx.Printer.Println(":check_mark: Deleted API key: ", apiKeyID)
		return nil
	},
}

func init() {
	apiKeyCmd.AddCommand(deleteApiKeyCmd)

	deleteApiKeyCmd.Flags().IntP("key-id", "k", 0, "API key ID")
}
