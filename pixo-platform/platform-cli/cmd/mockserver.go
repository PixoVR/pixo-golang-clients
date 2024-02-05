/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"encoding/json"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/mockserver"
	"github.com/spf13/cobra"
)

// mockserverCmd represents the mockserver rootCmd
var mockserverCmd = &cobra.Command{
	Use:   "mockserver",
	Short: "Starts a mock matchmaking server used for local testing",
	Long:  `Runs a mock matchmaking server that returns a static response determined by the server configuration or user input `,
	Run: func(cmd *cobra.Command, args []string) {

		//viper.AddConfigPath(".pixo")
		//viper.SetConfigName("server")

		address, _ := Ctx.ConfigManager.GetFlagOrConfigValue("gameserver-ip", cmd)
		port, _ := Ctx.ConfigManager.GetFlagOrConfigValue("gameserver-port", cmd)
		serverPort, _ := Ctx.ConfigManager.GetFlagOrConfigValue("matchmaker-port", cmd)
		moduleID, ok := Ctx.ConfigManager.GetIntFlagOrConfigValue("module-id", cmd)
		if !ok {
			moduleID = 1
		}
		serverVersion, _ := Ctx.ConfigManager.GetFlagOrConfigValue("server-version", cmd)
		orgID, _ := Ctx.ConfigManager.GetIntFlagOrConfigValue("org-id", cmd)
		sessionName, _ := Ctx.ConfigManager.GetFlagOrConfigValue("session-name", cmd)
		sessionID, _ := Ctx.ConfigManager.GetFlagOrConfigValue("session-id", cmd)
		owningUserName, _ := Ctx.ConfigManager.GetFlagOrConfigValue("owning-user-name", cmd)
		mapName, _ := Ctx.ConfigManager.GetFlagOrConfigValue("map-name", cmd)

		data := matchmaker.MatchResponse{
			Error:   false,
			Message: "Match found",
			MatchDetails: matchmaker.MatchDetails{
				IP:             address,
				Port:           port,
				SessionName:    sessionName,
				SessionID:      sessionID,
				OwningUserName: owningUserName,
				MapName:        mapName,
				ModuleVersion:  serverVersion,
				ModuleID:       moduleID,
				OrgID:          orgID,
			},
		}

		response, err := json.Marshal(data)
		if err != nil {
			Ctx.ConfigManager.Println(":exclamation: Could not marshal response: ", err)
			return
		}

		mockserver.Run(serverPort, Ctx.ConfigManager, "matchmaking/"+matchmaker.MatchmakingEndpoint, response)
	},
}

func init() {
	mpCmd.AddCommand(mockserverCmd)

	mockserverCmd.Flags().String("gameserver-ip", matchmaker.Localhost, "IP address of the game server to be returned in the response")
	mockserverCmd.Flags().String("server-version", "1.00.00", "Version of the server to be returned in the response")
	mockserverCmd.Flags().Int("gameserver-port", matchmaker.DefaultGameserverPort, "Port of the game server to be returned in the response")
	mockserverCmd.Flags().Int("matchmaker-port", 8080, "Port of the mock matchmaker server")
	mockserverCmd.Flags().String("map-name", "Default", "Name of the map to be returned in the response")
	mockserverCmd.Flags().StringP("session-name", "n", "Test", "Name of the session to be returned in the response")
	mockserverCmd.Flags().StringP("session-id", "i", "FB0HIFBMY8NAME99IS7C3WALKERB4D76", "ID of the session to be returned in the response")
	mockserverCmd.Flags().StringP("owning-user-name", "u", "PixoServer", "Name of the user that owns the session to be returned in the response")
	mockserverCmd.Flags().IntP("org-id", "o", 1, "Org ID to be returned in the response")
}
