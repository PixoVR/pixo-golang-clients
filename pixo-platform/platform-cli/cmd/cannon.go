/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// cannonCmd represents the cannon command
var cannonCmd = &cobra.Command{
	Use:   "cannon",
	Short: "Load testing tool",
	Long:  `Load test various elements of the platform`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(cannonCmd)
	cannonCmd.PersistentFlags().IntP("amount", "a", 50, "Number of requests to simulate")
	cannonCmd.PersistentFlags().IntP("concurrent", "c", 5, "Number of concurrent requests to simulate")
	cannonCmd.PersistentFlags().IntP("timeout", "t", 2, "Max duration of the entire test in minutes")
}

func canRunLoadTests(cmd *cobra.Command) error {
	user, err := Ctx.PlatformClient.GetUser(cmd.Context(), Ctx.PlatformClient.ActiveUserID())
	if err != nil {
		return err
	}

	if user.Org.Type != "platform" {
		return fmt.Errorf("only platform users can run load tests")
	}

	lifecycle := Ctx.ConfigManager.Lifecycle()
	if lifecycle == "prod" || lifecycle == "" {
		return fmt.Errorf("cannot run load tests against production")
	}

	return nil
}
