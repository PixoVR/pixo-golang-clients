/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/load/fixture"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/load/matchmaking"
	"github.com/spf13/cobra"
)

// cannonMatchmakingCmd represents the sessions start command
var cannonMatchmakingCmd = &cobra.Command{
	Use:   "matchmake",
	Short: "Load test matchmaking",
	Long:  `Run a load test simulating many matchmaking requests at once`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := canRunLoadTests(cmd); err != nil {
			return err
		}

		answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm(matchmakingQuestions(), cmd)
		if err != nil {
			return err
		}

		moduleID := forms.Int(answers["module"])
		semVer := forms.String(answers["server-version"])

		config := matchmaking.Config{
			Config: fixture.Config{PlatformFixture: Ctx, Command: cmd},
			Request: matchmaker.MatchRequest{
				ModuleID:      moduleID,
				ServerVersion: semVer,
			},
		}

		tester, err := matchmaking.NewLoadTester(config)
		if err != nil {
			return err
		}

		tester.Run()
		return nil
	},
}

func init() {
	cannonCmd.AddCommand(cannonMatchmakingCmd)
	cannonMatchmakingCmd.Flags().StringP("server-version", "v", "", "Semantic version of the multiplayer server version")
}
