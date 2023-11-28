/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/rs/zerolog/log"
	"net"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	mm matchmaker.Matchmaker
)

// mpCmd represents the mp command
var mpCmd = &cobra.Command{
	Use:   "mp",
	Short: "Manage Pixo Platform multiplayer resources",
	Long:  `Manage resources like server configurations, versions, triggers. Test game servers and matchmaking.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("connect").Value.String() == "true" {

			addr := input.GetConfigValue("gameserver", "PIXO_GAMESERVER")
			if addr == "" {
				log.Error().Msg("No gameserver address provided")
				return
			}

			splitAddr := strings.Split(addr, ":")
			if len(splitAddr) != 2 {
				log.Error().Str("addr", addr).Msg("Invalid gameserver address")
				return
			}

			gameserverHost := splitAddr[0]
			gameserverPort, err := strconv.Atoi(splitAddr[1])
			if err != nil {
				log.Error().Err(err).Msg("Could not parse gameserver port")
				return
			}

			udpAddr := &net.UDPAddr{IP: net.ParseIP(gameserverHost), Port: gameserverPort}
			if err := mm.DialGameserver(udpAddr); err != nil {
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
	},
}

func init() {
	mm = matchmaker.NewMatchmaker(input.GetConfigValue("matchmaking-api-url", "PIXO_PLATFORM_MATCHMAKING_URL"), input.GetConfigValue("token", "SECRET_KEY"))

	rootCmd.AddCommand(mpCmd)

	mpCmd.PersistentFlags().Bool("connect", false, "Whether UDP connection should be made when gameserver is found")
}
