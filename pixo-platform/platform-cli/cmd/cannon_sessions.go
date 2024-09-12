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
	"os"
)

var legacy bool

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

		moduleVersion, _ := Ctx.ConfigManager.GetFlagValue("version", cmd)
		mode, _ := Ctx.ConfigManager.GetFlagValue("mode", cmd)
		scenario, _ := Ctx.ConfigManager.GetFlagValue("scenario", cmd)
		focus, _ := Ctx.ConfigManager.GetFlagValue("focus", cmd)
		specialization, _ := Ctx.ConfigManager.GetFlagValue("specialization", cmd)

		eventPayload, _ := Ctx.ConfigManager.GetFlagValue("payload", cmd)
		if eventPayload == "" {
			payloadFile, _ := Ctx.ConfigManager.GetFlagValue("payload-file", cmd)
			if payloadFile != "" {
				payloadData, err := os.ReadFile(payloadFile)
				if err != nil {
					return err
				}
				eventPayload = string(payloadData)
			}
		}

		moduleID := forms.Int(answers["module"])

		config := sessions.Config{
			Config: fixture.Config{
				PlatformFixture: Ctx,
				Command:         cmd,
				Writer:          cmd.OutOrStdout(),
			},
			Legacy: legacy,

			Session: platform.Session{
				ModuleID:       moduleID,
				ModuleVersion:  moduleVersion,
				Mode:           mode,
				Scenario:       scenario,
				Focus:          focus,
				Specialization: specialization,
			},
			EventPayload: eventPayload,
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
	cannonSessionsCmd.Flags().String("version", "", "Module version to create sessions for")
	cannonSessionsCmd.Flags().String("mode", "", "Mode to create sessions for")
	cannonSessionsCmd.Flags().String("scenario", "", "Scenario to create sessions for")
	cannonSessionsCmd.Flags().String("focus", "", "Focus to create sessions for")
	cannonSessionsCmd.Flags().String("specialization", "", "Specialization to create sessions for")
	cannonSessionsCmd.Flags().String("payload", "", "Event payload to send when creating events")
	cannonSessionsCmd.Flags().String("payload-file", "", "File containing event payload to send when creating events")
	cannonSessionsCmd.Flags().BoolVar(&legacy, "legacy", false, "Use the legacy headset API for load testing")
}
