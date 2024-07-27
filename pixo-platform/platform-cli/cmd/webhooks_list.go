/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/loader"
	"github.com/spf13/cobra"
)

// webhooksListCmd represents the sessions start command
var webhooksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List webhooks",
	Long:  `List webhooks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		spinner := loader.NewLoader(cmd.Context(), "Getting webhooks...", Ctx.Printer)
		webhooks, err := Ctx.PlatformClient.GetWebhooks(cmd.Context(), &platform.WebhookParams{OrgID: Ctx.PlatformClient.ActiveOrgID()})
		spinner.Stop()
		if err != nil {
			Ctx.Printer.Println(":exclamation: Failed to get webhooks")
			return err
		}

		for _, webhook := range webhooks {
			Ctx.Printer.Println(fmt.Sprintf("%d. Description: %s\n    URL: %s", webhook.ID, webhook.Description, webhook.URL))
		}

		return nil
	},
}

func init() {
	webhooksCmd.AddCommand(webhooksListCmd)
}
