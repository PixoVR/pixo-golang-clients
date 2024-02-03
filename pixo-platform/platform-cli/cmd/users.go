/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// usersCmd represents the create command
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Managing users",
	Long:  `Manage users with the following command:`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)
}
