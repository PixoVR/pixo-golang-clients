/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// webhooksCmd represents the webhooks command
var webhooksCmd = &cobra.Command{
	Use:   "webhooks",
	Short: "Manage webhooks",
	Long:  `Create, list, and delete webhooks`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(webhooksCmd)
}
