/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"encoding/json"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/mockserver"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// mockserverCmd represents the mockserver rootCmd
var mockserverCmd = &cobra.Command{
	Use:   "mockserver",
	Short: "Starts a mock matchmaking server used for local testing",
	Long:  `Runs a mock matchmaking server that returns a static response determined by the server configuration or user input `,
	Run: func(cmd *cobra.Command, args []string) {

		viper.AddConfigPath(".pixo")
		viper.SetConfigName("server")
		viper.SetDefault("module-id", 1)
		viper.SetDefault("server-port", 8080)

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				log.Debug().Msg("ConfigFile file not found")
			} else {
				log.Debug().Err(err).Msg("ConfigFile file was found but another error was produced")
			}
		}

		address, _ := Ctx.ConfigManager.GetConfigValueOrAskUser("gameserver-ip", cmd)
		port, _ := Ctx.ConfigManager.GetConfigValueOrAskUser("gameserver-port", cmd)

		moduleID, _ := Ctx.ConfigManager.GetIntConfigValueOrAskUser("module-id", cmd)
		orgID, _ := Ctx.ConfigManager.GetIntConfigValueOrAskUser("org-id", cmd)

		sessionName, _ := Ctx.ConfigManager.GetConfigValueOrAskUser("session-name", cmd)
		sessionID, _ := Ctx.ConfigManager.GetConfigValueOrAskUser("session-id", cmd)
		owningUserName, _ := Ctx.ConfigManager.GetConfigValueOrAskUser("owning-user-name", cmd)
		mapName, _ := Ctx.ConfigManager.GetConfigValueOrAskUser("map-name", cmd)
		serverVersion, _ := Ctx.ConfigManager.GetConfigValueOrAskUser("server-version", cmd)

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
			log.Debug().Err(err).Msg("Failed to marshal response")
			return
		}

		mode := gin.ReleaseMode

		if isDebug {
			mode = gin.DebugMode
		}

		mockserver.Run(cmd, Ctx.ConfigManager, mode, matchmaker.MatchmakingEndpoint, response)
	},
}

func init() {
	mpCmd.AddCommand(mockserverCmd)

	mockserverCmd.Flags().StringP("server-port", "p", "8080", "Port to run the server on")
	mockserverCmd.Flags().String("gameserver-ip", "127.0.0.1", "IP address of the game server to be returned in the response")
	mockserverCmd.Flags().String("gameserver-port", "7777", "Port of the game server to be returned in the response")
	mockserverCmd.Flags().String("map-name", "Default", "Name of the map to be returned in the response")
	mockserverCmd.Flags().StringP("session-name", "n", "Test", "Name of the session to be returned in the response")
	mockserverCmd.Flags().StringP("session-id", "i", "FB0HIFBMY8NAME99IS7C3WALKERB4D76", "ID of the session to be returned in the response")
	mockserverCmd.Flags().StringP("owning-user-name", "u", "PixoServer", "Name of the user that owns the session to be returned in the response")
	mockserverCmd.Flags().IntP("org-id", "o", 1, "Org ID to be returned in the response")
}
