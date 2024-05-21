/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// versionCmd represents the config rootCmd
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Configure the CLI settings",
	Long: `Manage settings like region, org, and module ID.  This commands will prompt you for the settings if they are not already set.
`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("pixo version %s\n", cliVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
