/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/loader"
	"github.com/spf13/cobra"
	"net"
)

// matchmakeCmd represents the matchmake rootCmd
var matchmakeCmd = &cobra.Command{
	Use:   "matchmake",
	Short: "Request a multiplayer match",
	Long:  `Connect to the matchmaking service to receive a gameserver address`,
	RunE: func(cmd *cobra.Command, args []string) error {
		answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm(matchmakingQuestions(), cmd)
		if err != nil {
			return err
		}

		moduleID := forms.Int(answers["module"])
		semVer := forms.String(answers["server-version"])

		matchRequest := matchmaker.MatchRequest{
			ModuleID:      moduleID,
			ServerVersion: semVer,
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
}

func matchmakingQuestions() []config.Value {
	return []config.Value{
		{Question: moduleQuestion()},
		{Question: forms.Question{Type: forms.Input, Key: "server-version"}},
	}
}
