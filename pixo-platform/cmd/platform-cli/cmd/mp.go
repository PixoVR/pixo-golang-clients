/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// mpCmd represents the mp command
var mpCmd = &cobra.Command{
	Use:   "mp",
	Short: "Manage Pixo Platform multiplayer resources",
	Long:  `Manage resources like server configurations, versions, triggers. Test game servers and matchmaking.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mp called")
	},
}

func init() {
	rootCmd.AddCommand(mpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
