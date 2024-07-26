/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// modulesCmd represents the modules rootCmd
var modulesCmd = &cobra.Command{
	Use:   "modules",
	Short: "Manage modules",
	Long: `Manage modules and their versions.
`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(modulesCmd)
	modulesCmd.PersistentFlags().IntP("module-id", "m", 0, "Module ID")
}
