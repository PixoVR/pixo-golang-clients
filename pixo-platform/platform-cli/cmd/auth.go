/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// authCmd represents the auth rootCmd
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with the Pixo Platform",
	Long:  `Manage authentication and authorization with the Pixo Platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	authCmd.PersistentFlags().StringP("username", "u", "", "PixoVR username")
	authCmd.PersistentFlags().StringP("password", "p", "", "PixoVR password")
}
