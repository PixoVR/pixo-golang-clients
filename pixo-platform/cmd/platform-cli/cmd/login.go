/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/kyokomi/emoji"

	"github.com/spf13/cobra"
)

// loginCmd represents the login rootCmd
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the Pixo Platform",
	Long: `Your username and password can be provided in multiple ways:
	- rootCmd line flags --username and --password
	- local config file ./config.yaml
	- environment variables PIXO_USERNAME and PIXO_PASSWORD
	- global config file ~/.pixo/config.yaml
	Will prioritize in order of the above list, and will prompt the user if none is found.	
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := PlatformCtx.Authenticate(cmd.InOrStdin(), cmd.OutOrStdout()); err != nil {
			cmd.Println(emoji.Sprintf(":exclamation: Login failed. Please check your credentials and try again."))
			return nil
		}

		cmd.Println(emoji.Sprintf(":rocket: Login successful. Here is your API token: \n%s", PlatformCtx.PlatformClient.GetToken()))
		return nil
	},
}

func init() {
	authCmd.AddCommand(loginCmd)
}
