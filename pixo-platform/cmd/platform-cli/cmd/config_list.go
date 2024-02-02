/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
)

// listCmd represents the list rootCmd
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists configuration settings",
	Long:  `Lists configuration settings like region, org, and module ID.`,
	Run: func(cmd *cobra.Command, args []string) {
		region := PlatformCtx.ConfigManager.Region()
		if region != "" {
			cmd.Println(emoji.Sprintf(":earth_americas: Region: %s", region))
		}

		lifecycle := PlatformCtx.ConfigManager.Lifecycle()
		if lifecycle != "" {
			cmd.Println(emoji.Sprintf(":gear: Lifecycle: %s", lifecycle))
		}

		username, ok := PlatformCtx.ConfigManager.GetConfigValue("username")
		if ok {
			cmd.Println(emoji.Sprintf(":bust_in_silhouette: Username: %s", username))
		}

	},
}

func init() {
	configCmd.AddCommand(listCmd)
}
