/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/loader"
	"github.com/spf13/cobra"
)

// sessionsSimulateCmd represents the sessions start command
var sessionsSimulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "Simulate a session in headset",
	Long:  `Start a session, create events, and end the session to mimic headset interactions`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		questions := []config.Value{
			{Question: forms.Question{
				Type: forms.SelectID,
				Key:  "module-id",
				LabelFunc: func(item interface{}) string {
					module := item.(platform.Module)
					return fmt.Sprintf("%d: %s - %s", module.ID, module.Abbreviation, module.Name)
				},
				GetItemsFunc: func(ctx context.Context) (interface{}, error) {
					items, err := Ctx.PlatformClient.GetModules(cmd.Context())
					if err != nil {
						return nil, errors.New("unable to get modules")
					}
					return items, nil
				},
			}},
		}

		answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd)
		if err != nil {
			return err
		}

		moduleID := forms.Int(answers["module-id"])

		spinner := loader.NewLoader(cmd.Context(), "Finding IP Address...", Ctx.Printer)
		ipAddress, err := Ctx.PlatformClient.GetIPAddress()
		spinner.Stop()
		if err != nil {
			var ok bool
			ipAddress, ok = Ctx.ConfigManager.GetConfigValueOrAskUser("ip-address", cmd)
			if !ok {
				return errors.New("ip address not provided")
			}
		}

		spinner = loader.NewLoader(cmd.Context(), "Starting session...", Ctx.Printer)
		session := &platform.Session{
			ModuleID:  moduleID,
			IPAddress: ipAddress,
		}
		err = Ctx.PlatformClient.CreateSession(cmd.Context(), session)
		spinner.Stop()
		if err != nil {
			return err
		}

		Ctx.Printer.Printf(":white_check_mark: Session started for module %s - %s", session.Module.Abbreviation, session.Module.Name)
		Ctx.ConfigManager.SetIntConfigValue("session-id", session.ID)
		return nil
	},
}

func init() {
	sessionsCmd.AddCommand(sessionsSimulateCmd)
	sessionsSimulateCmd.Flags().StringP("module-id", "m", "", "Module ID")
}
