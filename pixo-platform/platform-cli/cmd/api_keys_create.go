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
		input := platform.APIKey{}

		username, _ := Ctx.ConfigManager.GetFlagValue("username", cmd)
		if username != "" {
			user, err := Ctx.PlatformClient.GetUserByUsername(cmd.Context(), username)
			if err != nil {
				return err
			}
			input.UserID = user.ID
		}

		spinner := loader.NewLoader(cmd.Context(), ":key: Creating API Key...", Ctx.Printer)
		defer spinner.Stop()

		apiKey, err := Ctx.PlatformClient.CreateAPIKey(cmd.Context(), input)
		if err != nil {
			return err
		}

		forUser := ""
		if input.UserID > 0 {
			forUser = " for " + username
		}

		Ctx.ConfigManager.SetConfigValue("api-key", apiKey.Key)
		Ctx.Printer.Printf(":heavy_check_mark: API key created%s: %s\n", forUser, apiKey.Key)
		//Ctx.Printer.Println(":heavy_check_mark: API key created: ", "*********") // use this when running the tape, --safe flag ?
		return nil
	},
}

func init() {
	apiKeyCmd.AddCommand(createApiKeyCmd)
}
