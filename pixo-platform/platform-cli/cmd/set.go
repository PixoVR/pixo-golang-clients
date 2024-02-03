/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/config"
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
)

// setCmd represents the set rootCmd
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set	a config value",
	Long:  `Can set a single value, or change region/lifecycle of the Pixo Platform APIs used`,
	Run: func(cmd *cobra.Command, args []string) {

		var env config.Env

		lifecycle := cmd.Flag("lifecycle").Value.String()
		if lifecycle != "" {
			env.Lifecycle = lifecycle
		}

		region := cmd.Flag("region").Value.String()
		if region != "" {
			env.Region = region
		}

		Ctx.ConfigManager.SetActiveEnv(env)

		key, err := cmd.Flags().GetString("key")
		if err != nil {
			cmd.Println(emoji.Sprintf(":exclamation: Unable to get key flag"))
			return
		}

		if key != "" {
			if val, err := cmd.Flags().GetString("val"); err != nil {
				cmd.Println(emoji.Sprintf(":exclamation: Unable to get value flag"))
				return
			} else if val != "" {
				Ctx.ConfigManager.SetConfigValue(key, val)
				cmd.Printf(emoji.Sprintf(":rocket: Config value %s set to %s\n", key, val))
				rootCmd.SetArgs([]string{"config"})
				_ = rootCmd.Execute()
				return
			} else {
				cmd.Println("Value must be provided")
				return
			}
		}

		username := cmd.Flag("username").Value.String()
		if username != "" {
			Ctx.ConfigManager.SetConfigValue("username", username)
		}

		password := cmd.Flag("password").Value.String()
		if password != "" {
			Ctx.ConfigManager.SetConfigValue("password", password)
		}

		cmd.Printf("Config updated successfully: %s\n", cfgFileFlagInput)
		rootCmd.SetArgs([]string{"config"})
		if cmd != nil {
			_ = cmd.Execute()
		}
	},
}

func init() {
	configCmd.AddCommand(setCmd)

	setCmd.Flags().StringP("username", "u", "", "Username to use for Pixo Platform")
	setCmd.Flags().StringP("password", "p", "", "Password to use for Pixo Platform")
	setCmd.Flags().StringP("key", "k", "", "Key of the config value to set")
	setCmd.Flags().StringP("val", "v", "", "Value of the config value to set")
}
