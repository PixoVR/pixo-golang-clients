/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
	"strings"
)

var (
	edit bool
)

// configCmd represents the config rootCmd
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure the CLI settings",
	Long: `Manage settings like region, org, and module ID.  This commands will prompt you for the settings if they are not already set.
`,
	Run: func(cmd *cobra.Command, args []string) {

		if edit {
			cmd.Println(emoji.Sprint(":file_folder: Opening config file in editor"))
			if err := Ctx.FileOpener.OpenEditor(Ctx.ConfigManager.ConfigFile()); err != nil {
				cmd.Println(emoji.Sprintf(":warning: Unable to open editor: %s", err))
			}
		}

		cmd.Println(emoji.Sprintf(":file_folder: Config file: %s", Ctx.ConfigManager.ConfigFile()))

		region := Ctx.ConfigManager.Region()
		if region != "" {
			cmd.Println(emoji.Sprintf(":earth_americas: Region: %s", region))
		}

		lifecycle := Ctx.ConfigManager.Lifecycle()
		if lifecycle != "" {
			cmd.Println(emoji.Sprintf(":gear:  Lifecycle: %s", lifecycle))
		}

		cmd.Println()

		activeEnv := Ctx.ConfigManager.GetActiveEnv()

		sensitiveList := []string{
			"password",
			"token",
			"api-key",
		}
		isSensitive := func(k string) bool {
			for _, s := range sensitiveList {
				if strings.Contains(k, s) {
					return true
				}
			}

			return false
		}

		cleanKey := func(k string) string {
			k = strings.Replace(k, "id", "ID", -1)
			k = strings.Replace(k, "api", "API", -1)
			k = strings.Replace(k, "-", " ", -1)
			return strings.Title(k)
		}

		for k, v := range activeEnv.EnvMap {
			if isSensitive(k) {
				v = "********"
			}
			cmd.Println(emoji.Sprintf(":arrow_right: %s: %s", cleanKey(k), v))
		}

	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.PersistentFlags().BoolVarP(&edit, "edit", "e", false, "Edit the config file in your default editor")
}
