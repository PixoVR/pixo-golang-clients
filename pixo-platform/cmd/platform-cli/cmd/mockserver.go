/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"encoding/json"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/mockserver"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// mockserverCmd represents the mockserver command
var mockserverCmd = &cobra.Command{
	Use:   "mockserver",
	Short: "Starts a mock matchmaking server used for local testing",
	Long:  `Runs a mock matchmaking server that returns a static response determined by the server configuration or user input `,
	Run: func(cmd *cobra.Command, args []string) {

		initLogger(cmd)

		viper.AddConfigPath(".pixo")
		viper.SetConfigName("server")
		viper.SetDefault("module-id", 1)

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				log.Warn().Msg("Config file not found")
			} else {
				log.Error().Err(err).Msg("Config file was found but another error was produced")
			}
		}

		data := matchmaker.MatchResponse{
			Error:   false,
			Message: "Match found",
			MatchDetails: matchmaker.MatchDetails{
				IP:             input.GetStringValue(cmd, "gameserver-ip", "GAMESERVER_IP"),
				Port:           input.GetStringValue(cmd, "gameserver-port", "GAMESERVER_PORT"),
				SessionName:    input.GetStringValue(cmd, "session-name", "SESSION_NAME"),
				SessionID:      input.GetStringValue(cmd, "session-id", "SESSION_ID"),
				OwningUserName: input.GetStringValue(cmd, "owning-user-name", "OWNING_USER_NAME"),
				MapName:        input.GetStringValue(cmd, "map-name", "MAP_NAME"),
				ModuleVersion:  input.GetStringValueOrAskUser(cmd, "server-version", "SERVER_VERSION"),
				ModuleID:       input.GetIntValueOrAskUser(cmd, "module-id", "MODULE_ID"),
				OrgID:          input.GetIntValueOrAskUser(cmd, "org-id", "ORG_ID"),
			},
		}

		response, err := json.Marshal(data)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal response")
			return
		}
		mockserver.Run(matchmaker.MatchmakingEndpoint, response)
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
