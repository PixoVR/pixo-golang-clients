/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"net"

	"github.com/spf13/cobra"
)

// matchmakeCmd represents the matchmake command
var matchmakeCmd = &cobra.Command{
	Use:   "matchmake",
	Short: "Test multiplayer matchmaking",
	Long: `Test multiplayer matchmaking.  This command will create a matchmake request and wait for a match to be found.
`,
	Run: func(cmd *cobra.Command, args []string) {
		mm := matchmaker.NewMatchmaker(input.GetConfigValue("matchmaking-api-url", config.PixoMatchmakingAPIURLEnvVarKey), input.GetConfigValue("token", config.PixoSecretKeyEnvVarKey))

		moduleID := input.GetIntValueOrAskUser(cmd, "module-id", "PIXO_MODULE_ID")
		semanticVersion := input.GetStringValueOrAskUser(cmd, "server-version", "PIXO_SERVER_VERSION")

		cmd.Println(fmt.Sprintf("Attempting to find a match for module %d with server version %s...", moduleID, semanticVersion))

		matchRequest := matchmaker.MatchRequest{
			ModuleID:      moduleID,
			ServerVersion: semanticVersion,
		}
		addr, err := mm.FindMatch(matchRequest)
		if err != nil {
			log.Error().Err(err).Msg("unable to find a match")
			return
		}

		cmd.Println("Match found! Gameserver connection info:", addr.String())

		viper.Set("gameserver", addr.String())
		if err := viper.WriteConfigAs(cfgFile); err != nil {
			log.Error().Err(err).Msg("Could not write config file")
		}

		if cmd.Flag("connect").Value.String() == "true" {
			gameserverReadLoop(cmd, mm, addr)
		}

	},
}

func gameserverReadLoop(cmd *cobra.Command, mm matchmaker.Matchmaker, addr *net.UDPAddr) {
	log.Info().Msg("Connecting to gameserver")
	if err := mm.DialGameserver(addr); err != nil {
		log.Error().Err(err).Msg("Could not connect to gameserver")
	}

	for {
		userInput := input.ReadFromUser("Enter message to send to gameserver: ")
		if userInput == "" || userInput == "exit" {
			break
		}

		response, err := mm.SendAndReceiveMessage([]byte(userInput))
		if err != nil {
			log.Error().Err(err).Msg("Could not send and receive message from gameserver")
		}

		cmd.Println(string(response))
	}

	log.Info().Msg("Closing connection to gameserver")
	if err := mm.CloseGameserverConnection(); err != nil {
		log.Error().Err(err).Msg("Could not close connection to gameserver")
	}

}

func init() {
	mpCmd.AddCommand(matchmakeCmd)

	matchmakeCmd.PersistentFlags().StringP("server-version", "s", "", "Server version to matchmake against")
}
