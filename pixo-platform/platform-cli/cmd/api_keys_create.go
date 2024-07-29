/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
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

		apiKey, err := Ctx.PlatformClient.CreateAPIKey(cmd.Context(), platform.APIKey{})
		if err != nil {
			return err
		}

		Ctx.ConfigManager.SetConfigValue("api-key", apiKey.Key)
		Ctx.Printer.Println(":heavy_check_mark: API key created: ", apiKey.Key)
		//Ctx.Printer.Println(":heavy_check_mark: API key created: ", "*********") // use this when running the tape, --safe flag ?
		return nil
	},
}

func init() {
	apiKeyCmd.AddCommand(createApiKeyCmd)
}
