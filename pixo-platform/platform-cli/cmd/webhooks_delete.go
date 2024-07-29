/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"context"
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
		questions := []config.Value{
			{Question: forms.Question{
				Type: forms.MultiSelectIDs,
				Key:  "webhook-ids",
				LabelFunc: func(item interface{}) string {
					webhook := item.(platform.Webhook)
					label := fmt.Sprintf("Org ID %d: ", webhook.ID)

					if webhook.Org != nil && webhook.Org.Name != "" {
						label = fmt.Sprintf("%s: ", webhook.Org.Name)
					}

					return fmt.Sprintf("%s%s", label, webhook.URL)
				},
				GetItemsFunc: func(ctx context.Context) (interface{}, error) {
					return Ctx.PlatformClient.GetWebhooks(cmd.Context(), nil)
				},
			}},
		}

		answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd)
		if err != nil {
			return err
		}

		ids := forms.IntSlice(answers["webhook-ids"])

		spinner := loader.NewLoader(cmd.Context(), "Deleting webhooks...", Ctx.Printer)
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
