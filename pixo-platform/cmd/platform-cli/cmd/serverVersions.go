/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// serverVersionsCmd represents the serverVersions command
var serverVersionsCmd = &cobra.Command{
	Use:   "serverVersions",
	Short: "Used to list and manage multiplayer server versions",
	Long: `Used to list and manage multiplayer server versions.
`,
	Run: func(cmd *cobra.Command, args []string) {
		initLogger(cmd)
	
		if err := cmd.Help(); err != nil {
			log.Error().Err(err).Msg("Could not display help")
			return
		}
	},
}

func init() {
	mpCmd.AddCommand(serverVersionsCmd)
}
