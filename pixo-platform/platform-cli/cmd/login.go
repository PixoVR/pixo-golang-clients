/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// loginCmd represents the login rootCmd
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the Pixo Platform",
	Long: `Your username and password can be provided in multiple ways:
	- rootCmd line flags --username and --password
	- environment variables PIXO_USERNAME and PIXO_PASSWORD
	- local config file ./config.yaml
	- global config file ~/.pixo/config.yaml
	Will prioritize in order of the above list, and will prompt the user if none is found.	
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := Ctx.Authenticate(cmd); err != nil {
			Ctx.ConfigManager.Println(":exclamation: Login failed. Please check your credentials and try again.")
		}

		msg := ":rocket: Login successful. Here is your API "
		token, ok := Ctx.ConfigManager.GetConfigValue("token")
		if ok {
			Ctx.ConfigManager.Println(msg, "token:\n", token)
			return
		}

		apiKey, ok := Ctx.ConfigManager.GetConfigValue("api-key")
		if ok {
			Ctx.ConfigManager.Println(msg, "key:\n", apiKey)
		}
	},
}

func init() {
	authCmd.AddCommand(loginCmd)
}
