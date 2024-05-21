/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// serversCmd represents the serverVersions rootCmd
var serversCmd = &cobra.Command{
	Use:   "servers",
	Short: "Used to list and manage multiplayer servers",
	Long: `Used to manage multiplayer server versions and live servers
`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	mpCmd.AddCommand(serversCmd)
}
