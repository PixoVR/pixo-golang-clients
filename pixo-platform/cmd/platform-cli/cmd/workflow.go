/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// workflowCmd represents the workflow command
var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Retrieve logs from the platform for a specific build workflow",
	Long:  `Retrieve logs for a specific build workflow`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("workflow called")
	},
}

func init() {
	logsCmd.AddCommand(workflowCmd)
}
