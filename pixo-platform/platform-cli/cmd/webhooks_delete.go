/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/loader"
	"github.com/spf13/cobra"
)

// webhooksDeleteCmd represents the sessions start command
var webhooksDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete webhooks",
	Long:  `Delete webhooks`,
	Run: func(cmd *cobra.Command, args []string) {
		id, ok := Ctx.ConfigManager.GetIntConfigValueOrAskUser("webhook-id", cmd)
		if !ok {
			Ctx.Printer.Println(":exclamation: ID is required")
			return
		}

		spinner := loader.NewLoader(cmd.Context(), "Getting webhooks...", Ctx.Printer)
		if err := Ctx.PlatformClient.DeleteWebhook(cmd.Context(), id); err != nil {
			Ctx.Printer.Println(":exclamation: Unable to delete webhook: ", err)
		}

		spinner.Stop()

		Ctx.Printer.Println(":white_check_mark: Webhook deleted")
	},
}

func init() {
	webhooksCmd.AddCommand(webhooksDeleteCmd)

	webhooksDeleteCmd.Flags().String("webhook-id", "", "Webhook ID")
}
