/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/loader"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/kyokomi/emoji"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// apiKeyCmd represents the apiKey command
var createApiKeyCmd = &cobra.Command{
	Use:   "create",
	Short: "Creating API keys",
	Long:  `Create API key with the following command:`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println(emoji.Sprintf(":key: Creating API Key"))

		spinner := loader.NewSpinner(cmd.OutOrStdout())

		input := platform.APIKey{
			//UserID: input.GetIntValue(cmd, "user-id", "PIXO_USER_ID"),
		}

		apiKey, err := PlatformCtx.PlatformClient.CreateAPIKey(cmd.Context(), input)
		if err != nil {
			cmd.Println("Error creating API key: ", err.Error())
			return err
		}

		if err := PlatformCtx.ConfigManager.SetAPIKey(apiKey.Key); err != nil {
			log.Error().Err(err).Msg("Could not set API key")
			return err
		}

		spinner.Stop()

		cmd.Println(emoji.Sprintf(":heavy_check_mark: API key created: %s", apiKey.Key))
		return nil
	},
}

func init() {
	apiKeyCmd.AddCommand(createApiKeyCmd)
}
