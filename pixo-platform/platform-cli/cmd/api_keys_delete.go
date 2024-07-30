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

// apiKeyCmd represents the apiKey command
var deleteApiKeyCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deleting an API key",
	Long:  `Delete API key with the following command:`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		questions := []config.Value{
			{Question: forms.Question{
				Type: forms.MultiSelectIDs,
				Key:  "key-ids",
				LabelFunc: func(value interface{}) string {
					item := value.(platform.APIKey)
					label := fmt.Sprintf("Key ID %d", item.ID)

					if item.User != nil {
						label = fmt.Sprintf("%s: ", label)
						if item.User.Email != "" {
							label = fmt.Sprintf("%s%s", label, item.User.Email)
						} else if item.User.Username != "" {
							label = fmt.Sprintf("%s%s", label, item.User.Username)
						}
						if item.User.Role != "" {
							label = fmt.Sprintf("%s - %s", label, item.User.Role)
						}
					}

					return label
				},
				GetItemsFunc: func(ctx context.Context) (interface{}, error) {
					return Ctx.PlatformClient.GetAPIKeys(cmd.Context(), nil)
				},
			}},
		}

		answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd)
		if err != nil {
			return err
		}

		ids := forms.IntSlice(answers["key-ids"])

		spinner := loader.NewLoader(cmd.Context(), ":key: Deleting API Key...", Ctx.Printer)
		defer spinner.Stop()
		for _, id := range ids {
			if err := Ctx.PlatformClient.DeleteAPIKey(cmd.Context(), id); err != nil {
				Ctx.Printer.Printf("Error deleting API key %d: %s\n", id, err.Error())
			} else {
				Ctx.Printer.Printf(":white_check_mark: Deleted API key: %d\n", id)
			}
		}

		return nil
	},
}

func init() {
	apiKeyCmd.AddCommand(deleteApiKeyCmd)

	deleteApiKeyCmd.Flags().StringP("key-ids", "k", "", "Comma separated list of API key IDs to delete")
}
