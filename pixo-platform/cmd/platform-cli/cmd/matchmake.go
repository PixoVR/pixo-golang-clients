/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/load"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/loader"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/kyokomi/emoji"
	"net"
	"time"

	"github.com/spf13/cobra"
)

// matchmakeCmd represents the matchmake rootCmd
var matchmakeCmd = &cobra.Command{
	Use:   "matchmake",
	Short: "Connect to the matchmaking service to receive a gameserver",
	Long: `Connect to the matchmaking service to receive a gameserver.
`,
	Run: func(cmd *cobra.Command, args []string) {

		moduleID, ok := Ctx.ConfigManager.GetIntConfigValueOrAskUser("module-id", cmd)
		if !ok {
			cmd.Println(emoji.Sprintf(":exclamation: Module ID not provided"))
			return
		}

		semanticVersion, ok := Ctx.ConfigManager.GetConfigValueOrAskUser("server-version", cmd)
		if !ok {
			cmd.Println(emoji.Sprintf(":exclamation: Server version not provided"))
			return
		}

		matchRequest := matchmaker.MatchRequest{
			ModuleID:      moduleID,
			ServerVersion: semanticVersion,
		}

		if numRequests, ok := Ctx.ConfigManager.GetIntFlagOrConfigValue("load", cmd); ok {
			timeout, _ := Ctx.ConfigManager.GetIntFlagOrConfigValue("timeout", cmd)
			config := load.Config{
				MatchmakingClient: Ctx.MatchmakingClient,
				Request:           matchRequest,
				Connections:       numRequests,
				Duration:          time.Duration(timeout) * time.Second,
				Reader:            cmd.InOrStdin(),
				Writer:            cmd.OutOrStdout(),
			}

			tester, err := load.NewLoadTester(config)
			if err != nil {
				cmd.Println(emoji.Sprintf(":exclamation: Could not create load tester: %s", err))
				return
			}

			tester.Run()
			return
		}

		cmd.Println(emoji.Sprintf(":magnifying_glass_tilted_left:Attempting to find a match for module %d with server version %s...", matchRequest.ModuleID, matchRequest.ServerVersion))

		spinner := loader.NewSpinner(cmd.OutOrStdout())

		addr, err := Ctx.MatchmakingClient.FindMatch(matchRequest)
		if err != nil {
			spinner.Stop()
			cmd.Println(emoji.Sprintf(":exclamation:Could not find match: %s", err))
			return
		}

		spinner.Stop()
		cmd.Println(emoji.Sprintf(":video_game:Match found! Gameserver address: %s", addr.String()))

		Ctx.ConfigManager.SetConfigValue("gameserver", addr.String())

		if connect {
			gameserverReadLoop(cmd, Ctx.MatchmakingClient, addr)
		}

	},
}

func gameserverReadLoop(cmd *cobra.Command, mm matchmaker.Matchmaker, addr *net.UDPAddr) {
	cmd.Println(emoji.Sprintf(":satellite:Connecting to gameserver at %s", addr.String()))
	if err := mm.DialGameserver(addr); err != nil {
		cmd.Println(emoji.Sprintf(":warning:Could not connect to gameserver at %s", addr.String()))
		return
	}

	Ctx.ConfigManager.SetWriter(cmd.OutOrStdout())
	for {
		Ctx.ConfigManager.ReadFromUser("Press enter to send a message to gameserver: ")
		userInput := Ctx.ConfigManager.ReadFromUser("Message to gameserver: ")
		if userInput == "" || userInput == "exit" {
			break
		}

		response, err := mm.SendAndReceiveMessage([]byte(userInput))
		if err != nil {
			cmd.Println(emoji.Sprintf(":warning:Could not send message to gameserver: %s", err))
			continue
		}

		cmd.Print(emoji.Sprintf(":arrow_right:Response: %s", response))
	}

	cmd.Println(emoji.Sprintf("\n:stop_sign:Closing connection to gameserver at %s", addr.String()))
	if err := mm.CloseGameserverConnection(); err != nil {
		cmd.Println(emoji.Sprintf(":warning:Could not close connection to gameserver at %s", addr.String()))
	}

}

func init() {
	mpCmd.AddCommand(matchmakeCmd)

	matchmakeCmd.Flags().IntP("load", "l", 0, "Number of connections in load test")
	matchmakeCmd.Flags().IntP("timeout", "t", 600, "Timeout in seconds for load test")
}
