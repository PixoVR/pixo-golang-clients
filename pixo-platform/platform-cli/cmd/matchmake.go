/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/load"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/loader"
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
			Ctx.ConfigManager.Println(":exclamation: Module ID not provided")
			return
		}

		semanticVersion, ok := Ctx.ConfigManager.GetConfigValueOrAskUser("server-version", cmd)
		if !ok {
			Ctx.ConfigManager.Println(":exclamation: Server version not provided")
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
				Ctx.ConfigManager.Println(":exclamation: Could not create load tester: ", err)
				return
			}

			tester.Run()
			return
		}

		Ctx.ConfigManager.Printf(":magnifying_glass_tilted_left:Attempting to find a match for module %d with server version %s...\n", matchRequest.ModuleID, matchRequest.ServerVersion)

		spinner := loader.NewLoader(cmd.Context(), "", Ctx.ConfigManager)

		addr, err := Ctx.MatchmakingClient.FindMatch(matchRequest)
		spinner.Stop()
		if err != nil {
			Ctx.ConfigManager.Println(":exclamation:Could not find match: ", err)
			return
		}

		Ctx.ConfigManager.Println(":video_game:Match found! Gameserver address: ", addr.String())

		Ctx.ConfigManager.SetConfigValue("gameserver", addr.String())

		if connect {
			gameserverReadLoop(addr)
		}
	},
}

func gameserverReadLoop(addr *net.UDPAddr) {
	Ctx.ConfigManager.Println(":satellite:Connecting to gameserver at ", addr.String())
	if err := Ctx.MatchmakingClient.DialGameserver(addr); err != nil {
		Ctx.ConfigManager.Println(":warning:Could not connect to gameserver at ", addr.String())
		return
	}

	for {
		userInput := Ctx.ConfigManager.ReadFromUser("message to gameserver")
		if userInput == "" || userInput == "exit" {
			break
		}

		if err := Ctx.MatchmakingClient.SendMessageToGameserver([]byte(userInput)); err != nil {
			Ctx.ConfigManager.Println(":warning:Could not send message to gameserver: ", err)
			continue
		}

		res, err := Ctx.MatchmakingClient.ReadMessageFromGameserver()
		if err != nil {
			Ctx.ConfigManager.Println(":warning:Could not read message from gameserver: ", err)
			continue
		}

		Ctx.ConfigManager.Println(":arrow_right:Response: ", string(res))
	}

	Ctx.ConfigManager.Println("\n:stop_sign:Closing connection to gameserver at ", addr.String())
	if err := Ctx.MatchmakingClient.CloseGameserverConnection(); err != nil {
		Ctx.ConfigManager.Println(":warning:Could not close connection to gameserver at ", addr.String())
	}

}

func init() {
	mpCmd.AddCommand(matchmakeCmd)

	matchmakeCmd.Flags().IntP("load", "l", 0, "Number of connections in load test")
	matchmakeCmd.Flags().IntP("timeout", "t", 600, "Timeout in seconds for load test")
}
