/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/loader"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/kyokomi/emoji"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"net"

	"github.com/spf13/cobra"
)

// matchmakeCmd represents the matchmake rootCmd
var matchmakeCmd = &cobra.Command{
	Use:   "matchmake",
	Short: "Test multiplayer matchmaking",
	Long: `Test multiplayer matchmaking.  This rootCmd will create a matchmake request and wait for a match to be found.
`,
	Run: func(cmd *cobra.Command, args []string) {

		moduleID := input.GetIntValueOrAskUser(cmd, "module-id", "PIXO_MODULE_ID")
		semanticVersion := input.GetStringValueOrAskUser(cmd, "server-version", "PIXO_SERVER_VERSION")

		cmd.Println(emoji.Sprintf(":magnifying_glass_tilted_left: Attempting to find a match for module %d with server version %s...", moduleID, semanticVersion))

		matchRequest := matchmaker.MatchRequest{
			ModuleID:      moduleID,
			ServerVersion: semanticVersion,
		}

		spinner := loader.NewSpinner(cmd.OutOrStdout())

		addr, err := PlatformCtx.MatchmakingClient.FindMatch(matchRequest)
		if err != nil {
			return
		}

		spinner.Stop()
		cmd.Println(emoji.Sprintf(":video_game: Match found! Gameserver address: %s", addr.String()))

		viper.Set("gameserver", addr.String())
		_ = viper.WriteConfigAs(cfgFile)

		if cmd.Flag("connect").Value.String() == "true" {
			gameserverReadLoop(cmd, PlatformCtx.MatchmakingClient, addr)
		}

	},
}

func gameserverReadLoop(cmd *cobra.Command, mm matchmaker.Matchmaker, addr *net.UDPAddr) {
	cmd.Println(emoji.Sprintf(":satellite: Connecting to gameserver at %s", addr.String()))
	if err := mm.DialGameserver(addr); err != nil {
		log.Error().Err(err).Msg("Could not connect to gameserver")
	}

	for {
		userInput := input.ReadFromUser(cmd.InOrStdin(), cmd.OutOrStdout(), "Enter message to send to gameserver: ")
		if userInput == "" || userInput == "exit" {
			break
		}

		response, err := mm.SendAndReceiveMessage([]byte(userInput))
		if err != nil {
			log.Error().Err(err).Msg("Could not send and receive message from gameserver")
		}

		cmd.Print(emoji.Sprintf(":arrow_right: Response: %s", response))
	}

	cmd.Println(emoji.Sprintf("\n:stop_sign: Closing connection to gameserver at %s", addr.String()))
	if err := mm.CloseGameserverConnection(); err != nil {
		cmd.Println(emoji.Sprintf(":warning: Could not close connection to gameserver at %s", addr.String()))
	}

}

func init() {
	mpCmd.AddCommand(matchmakeCmd)
}
