/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/load/fixture"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/load/sessions"
	"github.com/spf13/cobra"
)

// cannonSessionsCmd represents the sessions start command
var cannonSessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "Load test sessions",
	Long:  `Run a load test simulating many sessions at once`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := canRunLoadTests(cmd); err != nil {
			return err
		}

		questions := []config.Value{{Question: moduleQuestion()}}

		answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd)
		if err != nil {
			return err
		}

		moduleID := forms.Int(answers["module"])

		config := sessions.Config{
			Config: fixture.Config{PlatformFixture: Ctx, Command: cmd},
			Module: platform.Module{ID: moduleID},
		}

		tester, err := sessions.NewLoadTester(config)
		if err != nil {
			return err
		}

		tester.Run()
		return nil
	},
}

func init() {
	cannonCmd.AddCommand(cannonSessionsCmd)
}
