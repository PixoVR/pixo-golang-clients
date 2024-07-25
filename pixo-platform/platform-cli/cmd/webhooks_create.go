/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/loader"
	"github.com/spf13/cobra"
	"strings"
)

// webhooksCreateCmd represents the sessions start command
var webhooksCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a webhook",
	Long:  `Create a webhook`,
	RunE: func(cmd *cobra.Command, args []string) error {
		url, ok := Ctx.ConfigManager.GetConfigValueOrAskUser("url", cmd)
		if !ok {
			Ctx.Printer.Println(":exclamation: URL not provided")
			return nil
		}

		description, ok := Ctx.ConfigManager.GetFlagOrConfigValueOrAskUser("description", cmd)
		if !ok {
			Ctx.Printer.Println(":exclamation: DESCRIPTION not provided")
			return nil
		}

		generateTokenRes, err := Ctx.FormHandler.GetResponseFromUser("Generate token? (yes/no)")
		if err != nil {
			Ctx.Printer.Println(":exclamation: Unable to get generate token response from user: ", err)
		}

		generateToken := strings.ToLower(generateTokenRes) == "yes" || strings.ToLower(generateTokenRes) == "y"

		var webhookToken string
		if !generateToken {
			webhookToken, ok = Ctx.ConfigManager.GetConfigValueOrAskUser("webhook-token", cmd)
			if !ok {
				Ctx.Printer.Println(":warning: No token provided. Webhook will be insecure")
			}
		}

		spinner := loader.NewLoader(cmd.Context(), "Creating webhook...", Ctx.Printer)
		_, err = Ctx.PlatformClient.CreateWebhook(cmd.Context(), platform.Webhook{
			OrgID:         Ctx.PlatformClient.ActiveOrgID(),
			URL:           url,
			Description:   description,
			GenerateToken: &generateToken,
			Token:         webhookToken,
		})
		spinner.Stop()
		if err != nil {
			Ctx.Printer.Println(":exclamation: Unable to create webhook: ", err)
			return nil
		}

		Ctx.Printer.Println(":white_check_mark: Webhook created")
		return nil
	},
}

func init() {
	webhooksCmd.AddCommand(webhooksCreateCmd)

	webhooksCreateCmd.Flags().String("url", "", "URL of the webhook")
	webhooksCreateCmd.Flags().String("description", "", "Description of the webhook")
}
