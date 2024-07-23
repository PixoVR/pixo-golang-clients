/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	primary_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/legacy"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/loader"
	"github.com/spf13/cobra"
)

func oldAPILogin() {
	username, _ := Ctx.ConfigManager.GetConfigValue("username")
	password, _ := Ctx.ConfigManager.GetConfigValue("password")
	if err := Ctx.OldAPIClient.Login(username, password); err != nil {
		Ctx.ConfigManager.Println(":exclamation: Unable to login to old API: ", err)
	}
}

// webhooksCreateCmd represents the sessions start command
var webhooksCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a webhook",
	Long:  `Create a webhook`,
	RunE: func(cmd *cobra.Command, args []string) error {
		oldAPILogin()

		url, ok := Ctx.ConfigManager.GetConfigValueOrAskUser("url", cmd)
		if !ok {
			Ctx.ConfigManager.Println(":exclamation: URL not provided")
		}

		description, ok := Ctx.ConfigManager.GetFlagOrConfigValueOrAskUser("description", cmd)
		if !ok {
			Ctx.ConfigManager.Println(":exclamation: DESCRIPTION not provided")
		}

		spinner := loader.NewLoader(cmd.Context(), "Creating webhook...", Ctx.ConfigManager)
		err := Ctx.OldAPIClient.CreateWebhook(primary_api.Webhook{
			OrgID:       Ctx.PlatformClient.ActiveOrgID(),
			URL:         url,
			Description: description,
		})
		spinner.Stop()
		if err != nil {
			Ctx.ConfigManager.Println(":exclamation: Unable to create webhook: ", err)
			return nil
		}

		Ctx.ConfigManager.Println(":white_check_mark: Webhook created")
		return nil
	},
}

func init() {
	webhooksCmd.AddCommand(webhooksCreateCmd)

	webhooksCreateCmd.Flags().String("url", "", "URL of the webhook")
	webhooksCreateCmd.Flags().String("description", "", "Description of the webhook")
}
