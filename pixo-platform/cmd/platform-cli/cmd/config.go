/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure the CLI settings",
	Long: `Manage settings like region, org, and module ID.  This commands will prompt you for the settings if they are not already set.
`,
	Run: func(cmd *cobra.Command, args []string) {
		initLogger(cmd)
		if err := cmd.Help(); err != nil {
			log.Debug().Err(err).Msg("Could not display help")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
