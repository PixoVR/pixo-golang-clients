/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	multiplayerAllocator "github.com/PixoVR/pixo-golang-clients/pixo-platform/multiplayer-allocator"
	"github.com/spf13/cobra"
)

// buildCmd represents the build rootCmd
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Retrieve logs from the platform for a specific gameserver build",
	Long:  `Retrieve logs for a specific build`,
	Run: func(cmd *cobra.Command, args []string) {

		allocatorClient := multiplayerAllocator.NewClient(
			input.GetConfigValue("token", "PIXO_TOKEN"),
			input.GetConfigValue("lifecycle", "PIXO_LIFECYCLE"),
			input.GetConfigValue("region", "PIXO_REGION"),
		)

		workflows, err := allocatorClient.GetBuildWorkflows()
		if err != nil {
			cmd.Println("Error getting build workflows: ", err)
			return
		}

		logsCh, err := allocatorClient.GetBuildWorkflowLogs(workflows[0].Name)
		if err != nil {
			cmd.Println("Error getting build workflow logs: ", err)
			return
		}

		for log := range logsCh {
			cmd.Println(log)
		}

		cmd.Println("Stop streaming logs")
	},
}

func init() {
	logsCmd.AddCommand(buildCmd)
}
