/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// matchmakeCmd represents the matchmake command
var matchmakeCmd = &cobra.Command{
	Use:   "matchmake",
	Short: "Test multiplayer matchmaking",
	Long: `Test multiplayer matchmaking.  This command will create a matchmake request and wait for a match to be found.
`,
	Run: func(cmd *cobra.Command, args []string) {
		mm := matchmaker.NewMatchmaker(input.GetConfigValue("matchmaking-api-url", "PIXO_PLATFORM_MATCHMAKING_URL"), input.GetConfigValue("token", "SECRET_KEY"))

		moduleID := input.GetIntValue(cmd, "module-id", "PIXO_MODULE_ID")
		orgID := input.GetIntValue(cmd, "org-id", "PIXO_ORG_ID")
		semanticVersion := input.GetStringValue(cmd, "server-version", "PIXO_SERVER_VERSION")

		matchRequest := matchmaker.MatchRequest{
			OrgID:         orgID,
			ModuleID:      moduleID,
			ServerVersion: semanticVersion,
		}
		addr, err := mm.Connect(matchRequest)
		if err != nil {
			log.Error().Err(err).Msg("unable to find a match")
			return
		}

		cmd.Println("Match found! Gameserver connection info:", addr.String())

		viper.Set("gameserver", addr.String())
		if err := viper.WriteConfigAs(cfgFile); err != nil {
			log.Error().Err(err).Msg("Could not write config file")
		}

	},
}

func init() {
	mpCmd.AddCommand(matchmakeCmd)

	matchmakeCmd.PersistentFlags().StringP("server-version", "s", "", "Server version to matchmake against")
}
