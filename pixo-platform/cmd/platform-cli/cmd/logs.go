/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// logsCmd represents the logs rootCmd
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Retrieve logs from the platform",
	Long:  `Retrieve logs from different components of the platform`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
