/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"encoding/json"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/allocator"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/mockserver"
	"github.com/spf13/cobra"
)

// mockserverCmd represents the mockserver rootCmd
var mockserverCmd = &cobra.Command{
	Use:   "mockserver",
	Short: "Starts a mock matchmaking server used for local testing",
	Long:  `Runs a mock matchmaking server that returns a static response determined by the server configuration or user input `,
	Run: func(cmd *cobra.Command, args []string) {
		address, _ := Ctx.ConfigManager.GetFlagOrConfigValue("gameserver-ip", cmd)
		port, _ := Ctx.ConfigManager.GetFlagOrConfigValue("gameserver-port", cmd)
		serverPort, _ := Ctx.ConfigManager.GetFlagOrConfigValue("matchmaker-port", cmd)
		sessionID, _ := Ctx.ConfigManager.GetFlagOrConfigValue("session-id", cmd)
		owningUserName, _ := Ctx.ConfigManager.GetFlagOrConfigValue("owning-user-name", cmd)
		mapName, _ := Ctx.ConfigManager.GetFlagOrConfigValue("map-name", cmd)

		data := matchmaker.MatchResponse{
			Error:   false,
			Message: "Match found",
			MatchDetails: matchmaker.MatchDetails{
				IP:             address,
				Port:           port,
				SessionID:      sessionID,
				OwningUserName: owningUserName,
				MapName:        mapName,
			},
		}

		response, err := json.Marshal(data)
		if err != nil {
			Ctx.Printer.Println(":exclamation: Could not marshal response: ", err)
			return
		}

		mockserver.Run(serverPort, Ctx.Printer, "matchmaking/"+matchmaker.MatchmakingEndpoint, response)
	},
}

func init() {
	mpCmd.AddCommand(mockserverCmd)

	mockserverCmd.Flags().Int("matchmaker-port", 8080, "Port of the mock matchmaker server")
	mockserverCmd.Flags().String("gameserver-ip", allocator.Localhost, "IP address of the game server to be returned in the response")
	mockserverCmd.Flags().Int("gameserver-port", allocator.DefaultGameserverPort, "Port of the game server to be returned in the response")
	mockserverCmd.Flags().StringP("session-id", "i", "FB0HIFBMY8NAME99IS7C3WALKERB4D76", "ID of the session to be returned in the response")
	mockserverCmd.Flags().StringP("owning-user-name", "u", "PixoServer", "Name of the user that owns the session to be returned in the response")
	mockserverCmd.Flags().String("map-name", "Default", "Name of the map to be returned in the response")
}
