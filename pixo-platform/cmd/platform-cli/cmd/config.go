/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/editor"
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	edit       bool
	FileOpener = editor.NewFileOpener("")
)

// configCmd represents the config rootCmd
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure the CLI settings",
	Long: `Manage settings like region, org, and module ID.  This commands will prompt you for the settings if they are not already set.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if edit {
			cmd.Println(emoji.Sprint(":file_folder:Opening config file in editor"))
			return FileOpener.OpenEditor(viper.ConfigFileUsed())
		}

		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.PersistentFlags().BoolVarP(&edit, "edit", "e", false, "Edit the config file in your default editor")
}
