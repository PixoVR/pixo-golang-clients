/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"errors"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/loader"
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
		apiKey, ok := Ctx.ConfigManager.GetFlagValue("key", cmd)
		if ok {
			Ctx.PlatformClient.SetAPIKey(apiKey)
			if _, err := Ctx.PlatformClient.GetControlTypes(cmd.Context()); err != nil {
				return errors.New("invalid API key")
			}

			Ctx.ConfigManager.SetConfigValue("api-key", apiKey)
			Ctx.Printer.Println(":rocket: Login with API key successful.")
			return nil
		}

		username, _ := Ctx.ConfigManager.GetConfigValue("auth-username")
		password, _ := Ctx.ConfigManager.GetConfigValue("auth-password")
		if username == "" || password == "" {

			var questions []config.Value

			if username == "" {
				questions = append(questions, config.Value{
					Question: forms.Question{
						Type: forms.Input,
						Key:  "username",
					},
				})
			}

			if password == "" {
				questions = append(questions, config.Value{
					Question: forms.Question{
						Type: forms.SensitiveInput,
						Key:  "password",
					},
				})
			}

			answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd)
			if err != nil {
				Ctx.Printer.Println(":exclamation: Login failed")
				return err
			}

			username = forms.String(answers["username"])
			Ctx.ConfigManager.SetConfigValue("auth-username", username)

			password := forms.String(answers["password"])
			Ctx.ConfigManager.SetConfigValue("auth-password", password)
		}

		spinner := loader.NewLoader(cmd.Context(), "Logging into the Pixo Platform...", Ctx.Printer)
		defer spinner.Stop()

		if err := Ctx.PlatformClient.Login(username, password); err != nil {
			return err
		}

		Ctx.ConfigManager.SetConfigValue("auth-token", Ctx.PlatformClient.GetToken())
		Ctx.ConfigManager.SetIntConfigValue("auth-user-id", Ctx.PlatformClient.ActiveUserID())

		msg := ":rocket: Login successful. Here is your API token:\n"
		Ctx.Printer.Println(msg, Ctx.PlatformClient.GetToken())
		return nil
	},
}

func init() {
	authCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringP("key", "k", "", "API key")
}
