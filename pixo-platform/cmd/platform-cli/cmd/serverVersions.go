/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// serverVersionsCmd represents the serverVersions rootCmd
var serverVersionsCmd = &cobra.Command{
	Use:   "serverVersions",
	Short: "Used to list and manage multiplayer server versions",
	Long: `Used to list and manage multiplayer server versions.
`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	mpCmd.AddCommand(serverVersionsCmd)
}
