/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/loader"
	"github.com/spf13/cobra"
)

// webhooksCreateCmd represents the sessions start command
var webhooksCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a webhook",
	Long:  `Create a webhook`,
	RunE: func(cmd *cobra.Command, args []string) error {
		questions := []config.Value{
			{Question: forms.Question{Type: forms.Input, Key: "url"}},
			{Question: forms.Question{Type: forms.Input, Key: "description"}},
			{Question: forms.Question{Type: forms.Confirm, Key: "generate-token", Prompt: "Generate token automatically?"}},
		}

		answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd)
		if err != nil {
			return err
		}

		url := forms.String(answers["url"])
		description := forms.String(answers["description"])
		generateToken := forms.Bool(answers["generate-token"])

		if !generateToken {
			Ctx.Printer.Println(":warning: No token provided. Webhook will be insecure")
		}

		spinner := loader.NewLoader(cmd.Context(), "Creating webhook...", Ctx.Printer)
		webhook, err := Ctx.PlatformClient.CreateWebhook(cmd.Context(), platform.Webhook{
			OrgID:         Ctx.PlatformClient.ActiveOrgID(),
			URL:           url,
			Description:   description,
			GenerateToken: &generateToken,
		})
		spinner.Stop()
		if err != nil {
			return err
		}

		Ctx.Printer.Println(":white_check_mark: Webhook created")
		if webhook.Token != "" {
			Ctx.Printer.Println("Token: ", webhook.Token)
		}
		return nil
	},
}

func init() {
	webhooksCmd.AddCommand(webhooksCreateCmd)

	webhooksCreateCmd.Flags().String("url", "", "URL of the webhook")
	webhooksCreateCmd.Flags().String("description", "", "Description of the webhook")
	webhooksCreateCmd.Flags().StringP("generate-token", "g", "", "Description of the webhook")
}
