/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
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

		region := Ctx.ConfigManager.Region()
		if region != "" {
			cmd.Println(emoji.Sprintf(":earth_americas: Region: %s", region))
		}

		lifecycle := Ctx.ConfigManager.Lifecycle()
		if lifecycle != "" {
			cmd.Println(emoji.Sprintf(":gear:  Lifecycle: %s", lifecycle))
		}

		username, ok := Ctx.ConfigManager.GetConfigValue("username")
		if ok {
			cmd.Println(emoji.Sprintf(":bust_in_silhouette: Username: %s", username))
		}

	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.PersistentFlags().BoolVarP(&edit, "edit", "e", false, "Edit the config file in your default editor")
}
