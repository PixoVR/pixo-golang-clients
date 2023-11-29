/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/config"
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
		if cmd.Flag("connect").Value.String() != "true" {
			if err := cmd.Help(); err != nil {
				log.Error().Err(err).Msg("Could not display help")
				return
			}
		} else {

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
			gameserverReadLoop(cmd, mm, udpAddr)
		}
	},
}

func init() {
	mm = matchmaker.NewMatchmaker(input.GetConfigValue("matchmaking-api-url", config.PixoMatchmakingAPIURLEnvVarKey), input.GetConfigValue("token", config.PixoSecretKeyEnvVarKey))

	rootCmd.AddCommand(mpCmd)

	mpCmd.PersistentFlags().StringP("server-version", "v", "", "Semantic Version of the multiplayer server version")
	mpCmd.PersistentFlags().Bool("connect", false, "Whether to connect to the gameserver found from a match request")
}
