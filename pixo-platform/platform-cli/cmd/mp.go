/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"errors"
	"net"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	connect bool
)

// mpCmd represents the mp rootCmd
var mpCmd = &cobra.Command{
	Use:   "mp",
	Short: "Manage Pixo Platform multiplayer resources",
	Long:  `Manage resources like server configurations, versions, triggers. Test game servers and matchmaking.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if connect {
			addr, ok := Ctx.ConfigManager.GetConfigValueOrAskUser("gameserver", cmd)
			if !ok {
				return errors.New("no gameserver address provided")
			}

			splitAddr := strings.Split(addr, ":")
			if len(splitAddr) != 2 {
				return errors.New("invalid gameserver address")
			}

			gameserverHost := splitAddr[0]
			gameserverPort, err := strconv.Atoi(splitAddr[1])
			if err != nil {
				return err
			}

			udpAddr := &net.UDPAddr{IP: net.ParseIP(gameserverHost), Port: gameserverPort}
			gameserverReadLoop(udpAddr)
		} else {
			_ = cmd.Help()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(mpCmd)

	mpCmd.PersistentFlags().StringP("server-version", "v", "", "Semantic SemanticVersion of the multiplayer server version")
	mpCmd.PersistentFlags().BoolVarP(&connect, "connect", "c", false, "Whether to connect to the gameserver found from a match request")
	mpCmd.PersistentFlags().StringP("module-id", "m", "", "Module ID")
}
