/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// apiKeyCmd represents the apiKey command
var apiKeyCmd = &cobra.Command{
	Use:   "apiKeys",
	Short: "Handling API keys",
	Long:  `Manage API key with the following command:`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(apiKeyCmd)

	apiKeyCmd.Flags().StringP("user-id", "u", "", "User ID")
}
