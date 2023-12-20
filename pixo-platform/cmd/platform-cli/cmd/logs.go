/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Retrieve logs from the platform",
	Long: `Retrieve logs from different components of the platform:
		- gameservers
		- build workflows
`,
	Run: func(cmd *cobra.Command, args []string) {
		initLogger(cmd)
		if err := cmd.Help(); err != nil {
			log.Debug().Err(err).Msg("Could not display logs help")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)

	// logsCmd.PersistentFlags().String("foo", "", "A help for foo")
	// logsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
