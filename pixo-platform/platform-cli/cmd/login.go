/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := Ctx.Authenticate(cmd); err != nil {
			Ctx.Printer.Println(":exclamation: Login failed. Please check your credentials and try again.")
		}

		questions := []config.Value{
			{Question: forms.Question{Type: forms.Input, Key: "username"}},
			{Question: forms.Question{Type: forms.SensitiveInput, Key: "password"}},
		}

		answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd)
		if err != nil {
			Ctx.Printer.Println(":exclamation: Login failed")
			return err
		}

		username := forms.String(answers["username"])
		Ctx.ConfigManager.SetConfigValue("username", username)

		password := forms.String(answers["password"])
		Ctx.ConfigManager.SetConfigValue("password", password)

		if err := Ctx.Authenticate(cmd); err != nil {
			Ctx.Printer.Println(":exclamation: Login failed. Please check your credentials and try again.")
		}

		msg := ":rocket: Login successful. Here is your API "

		token, ok := Ctx.ConfigManager.GetConfigValue("token")
		if ok {
			Ctx.Printer.Println(msg, "token:\n", token)
			return nil
		}

		apiKey, ok := Ctx.ConfigManager.GetConfigValue("api-key")
		if ok {
			Ctx.Printer.Println(msg, "key:\n", apiKey)
		}

		return nil
	},
}

func init() {
	authCmd.AddCommand(loginCmd)
}
