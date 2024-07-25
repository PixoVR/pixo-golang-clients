/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/legacy"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/loader"
	"github.com/spf13/cobra"
)

// apiKeyCmd represents the apiKey command
var createApiKeyCmd = &cobra.Command{
	Use:   "create",
	Short: "Creating API keys",
	Long:  `Create API key with the following command:`,
	RunE: func(cmd *cobra.Command, args []string) error {
		spinner := loader.NewLoader(cmd.Context(), ":key: Creating API Key...", Ctx.Printer)
		defer spinner.Stop()

		input := platform.APIKey{
			//UserID: input.GetIntValue(cmd, "user-id", "PIXO_USER_ID"),
		}

		apiKey, err := Ctx.PlatformClient.CreateAPIKey(cmd.Context(), input)
		if err != nil {
			Ctx.Printer.Println("Error creating API key: ", err)

			return err
		}

		Ctx.ConfigManager.SetConfigValue("api-key", apiKey.Key)

		Ctx.Printer.Println(":heavy_check_mark: API key created: ", apiKey.Key)
		return nil
	},
}

func init() {
	apiKeyCmd.AddCommand(createApiKeyCmd)
}
