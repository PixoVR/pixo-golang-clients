/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists configuration settings",
	Long:  `Lists configuration settings like region, org, and module ID.`,
	Run: func(cmd *cobra.Command, args []string) {
		initLogger(cmd)

		config := viper.AllSettings()
		for key, value := range config {
			cmd.Println(key, ":", value)
		}
	},
}

func init() {
	configCmd.AddCommand(listCmd)
}
