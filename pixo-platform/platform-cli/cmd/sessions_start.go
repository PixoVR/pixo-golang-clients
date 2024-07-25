/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/loader"
	"github.com/spf13/cobra"
)

// sessionsStartCmd represents the sessions start command
var sessionsStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a session",
	Long:  `Start a session to mimic headset interactions`,
	RunE: func(cmd *cobra.Command, args []string) error {
		moduleID, ok := Ctx.ConfigManager.GetIntConfigValueOrAskUser("module-id", cmd)
		if !ok {
			Ctx.Printer.Println(":exclamation: Module ID not provided")
		}

		spinner := loader.NewLoader(cmd.Context(), "Finding IP Address...", Ctx.Printer)
		ipAddress, err := Ctx.PlatformClient.GetIPAddress()
		spinner.Stop()
		if err != nil {
			ipAddress, ok = Ctx.ConfigManager.GetConfigValueOrAskUser("ip-address", cmd)
			if !ok {
				Ctx.Printer.Println(":exclamation: ip address not provided: ", err)
			}
		}

		spinner = loader.NewLoader(cmd.Context(), "Starting session...", Ctx.Printer)
		session, err := Ctx.PlatformClient.CreateSession(cmd.Context(), moduleID, ipAddress, "")
		spinner.Stop()
		if err != nil {
			Ctx.Printer.Println(":exclamation: Unable to create session: ", err)
			return nil
		}

		Ctx.Printer.Printf(":white_check_mark: Session started for module %d with ID %d", moduleID, session.ID)
		Ctx.ConfigManager.SetIntConfigValue("session-id", session.ID)
		return nil
	},
}

func init() {
	sessionsCmd.AddCommand(sessionsStartCmd)
}
