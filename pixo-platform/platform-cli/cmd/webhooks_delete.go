/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/loader"
	"github.com/spf13/cobra"
)

// webhooksDeleteCmd represents the sessions start command
var webhooksDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete webhooks",
	Long:  `Delete webhooks`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var webhooks []platform.Webhook
		if _, ok := Ctx.ConfigManager.GetFlagValue("webhook-ids", cmd); !ok {
			webhooks, err = Ctx.PlatformClient.GetWebhooks(cmd.Context(), nil)
			if err != nil {
				return err
			}
		}

		options := make([]forms.Option, len(webhooks))
		for i, webhook := range webhooks {
			labelPrefix := fmt.Sprintf("Org ID %d - ", webhook.OrgID)
			if webhook.Org != nil {
				labelPrefix = fmt.Sprintf("%s - ", webhook.Org.Name)
			}
			options[i] = forms.Option{
				Label: fmt.Sprintf("%d: %s%s", webhook.ID, labelPrefix, webhook.URL),
				Value: fmt.Sprint(webhook.ID),
			}
		}

		questions := []config.Value{
			{Question: forms.Question{
				Type:    forms.MultiSelectIDs,
				Key:     "webhook-ids",
				Prompt:  "Select webhooks to delete",
				Options: options,
			}},
		}

		answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd)
		if err != nil {
			return err
		}

		ids := forms.IntSlice(answers["webhook-ids"])

		spinner := loader.NewLoader(cmd.Context(), "Getting webhooks...", Ctx.Printer)
		for _, id := range ids {
			if err := Ctx.PlatformClient.DeleteWebhook(cmd.Context(), id); err != nil {
				Ctx.Printer.Printf(":exclamation: Unable to delete webhook %d: %s\n", id, err.Error())
			} else {
				Ctx.Printer.Printf(":white_check_mark: Webhook %d deleted\n", id)
			}
		}

		spinner.Stop()
		return nil
	},
}

func init() {
	webhooksCmd.AddCommand(webhooksDeleteCmd)

	webhooksDeleteCmd.Flags().String("webhook-ids", "", "Comma separated list of webhook ids to delete")
}
