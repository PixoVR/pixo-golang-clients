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

// apiKeyCmd represents the apiKey command
var deleteApiKeyCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deleting an API key",
	Long:  `Delete API key with the following command:`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var keys []platform.APIKey
		if _, ok := Ctx.ConfigManager.GetFlagValue("key-ids", cmd); !ok {
			keys, err = Ctx.PlatformClient.GetAPIKeys(cmd.Context(), nil)
			if err != nil {
				return err
			}
		}

		options := make([]forms.Option, len(keys))
		for i, key := range keys {
			labelPrefix := fmt.Sprintf("Key ID %d: ", key.ID)

			if key.User != nil {
				if key.User.Email != "" {
					labelPrefix = fmt.Sprintf("%s%s", labelPrefix, key.User.Email)
				} else if key.User.Username != "" {
					labelPrefix = fmt.Sprintf("%s%s", labelPrefix, key.User.Username)
				}
				if key.User.Role != "" {
					labelPrefix = fmt.Sprintf("%s - %s", labelPrefix, key.User.Role)
				}
			}

			options[i] = forms.Option{
				Label: labelPrefix,
				Value: fmt.Sprint(key.ID),
			}
		}

		questions := []config.Value{
			{Question: forms.Question{
				Type:    forms.MultiSelectIDs,
				Key:     "key-ids",
				Prompt:  "Select API keys to delete",
				Options: options,
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
