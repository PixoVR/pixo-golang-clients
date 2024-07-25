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

var (
	generateToken bool
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
		}

		if !generateToken {
			configValue := config.Value{
				Question: forms.Question{
					Type: forms.Input,
					Key:  "webhook-token",
				},
			}
			questions = append(questions, configValue)
		}

		answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd)
		if err != nil {
			Ctx.Printer.Printf(":exclamation: %v\n", err)
		}

		url := forms.String(answers["url"])
		if url == "" {
			Ctx.Printer.Println(":exclamation: URL not provided")
			return nil
		}

		description := forms.String(answers["description"])
		if description == "" {
			Ctx.Printer.Println(":exclamation: DESCRIPTION not provided")
			return nil
		}

		webhookToken := forms.String(answers["webhook-token"])
		if !generateToken && webhookToken == "" {
			if err = Ctx.FormHandler.Confirm("Generate token automatically?", &generateToken); err != nil || !generateToken {
				Ctx.Printer.Println(":warning: No token provided. Webhook will be insecure")
			}
		}

		spinner := loader.NewLoader(cmd.Context(), "Creating webhook...", Ctx.Printer)
		webhook, err := Ctx.PlatformClient.CreateWebhook(cmd.Context(), platform.Webhook{
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
	webhooksCreateCmd.Flags().BoolVarP(&generateToken, "generate-token", "g", false, "Description of the webhook")
}
