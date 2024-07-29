/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/load"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/loader"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		questions := []config.Value{
			{Question: forms.Question{
				Type: forms.SelectID,
				Key:  "module-id",
				LabelFunc: func(i interface{}) string {
					item := i.(platform.Module)
					return fmt.Sprintf("%d: %s - %s", item.ID, item.Abbreviation, item.Name)
				},
				GetItemsFunc: func(ctx context.Context) (interface{}, error) {
					return Ctx.PlatformClient.GetModules(cmd.Context())
				},
			}},
			{Question: forms.Question{Type: forms.Input, Key: "server-version"}},
		}

		answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd)
		if err != nil {
			return err
		}

		moduleID := forms.Int(answers["module-id"])
		semVer := forms.String(answers["server-version"])

		matchRequest := matchmaker.MatchRequest{
			ModuleID:      moduleID,
			ServerVersion: semVer,
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
				return err
			}

			tester.Run()
			return nil
		}

		Ctx.Printer.Printf(":magnifying_glass_tilted_left:Attempting to find a match for module %d with server version %s...\n", matchRequest.ModuleID, matchRequest.ServerVersion)

		spinner := loader.NewLoader(cmd.Context(), "", Ctx.Printer)

		addr, err := Ctx.MatchmakingClient.FindMatch(matchRequest)
		spinner.Stop()
		if err != nil {
			return err
		}

		Ctx.Printer.Println(":video_game:Match found! Gameserver address: ", addr.String())
		Ctx.ConfigManager.SetConfigValue("gameserver", addr.String())

		if connect {
			gameserverReadLoop(addr)
		}

		return nil
	},
}

func gameserverReadLoop(addr *net.UDPAddr) {
	Ctx.Printer.Println(":satellite:Connecting to gameserver at ", addr.String())
	if err := Ctx.MatchmakingClient.DialGameserver(addr); err != nil {
		Ctx.Printer.Println(":warning:Could not connect to gameserver at ", addr.String())
		return
	}

	for {
		question := &forms.Question{Prompt: "message to gameserver"}
		err := Ctx.FormHandler.GetResponseFromUser(question)
		answer := forms.String(question.Answer)
		if err != nil || answer == "" || answer == "exit" {
			break
		}

		if err = Ctx.MatchmakingClient.SendMessageToGameserver([]byte(answer)); err != nil {
			Ctx.Printer.Println(":warning:Could not send message to gameserver: ", err)
			continue
		}

		res, err := Ctx.MatchmakingClient.ReadMessageFromGameserver()
		if err != nil {
			Ctx.Printer.Println(":warning:Could not read message from gameserver: ", err)
			continue
		}

		Ctx.Printer.Println(":arrow_right:Response: ", string(res))
	}

	Ctx.Printer.Println("\n:stop_sign:Closing connection to gameserver at ", addr.String())
	if err := Ctx.MatchmakingClient.CloseGameserverConnection(); err != nil {
		Ctx.Printer.Println(":warning:Could not close connection to gameserver at ", addr.String())
	}

}

func init() {
	mpCmd.AddCommand(matchmakeCmd)

	matchmakeCmd.Flags().IntP("load", "l", 0, "Number of connections in load test")
	matchmakeCmd.Flags().IntP("timeout", "t", 600, "Timeout in seconds for load test")
}
