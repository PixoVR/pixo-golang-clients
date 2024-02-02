/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/config"
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

		PlatformCtx.ConfigManager.SetActiveEnv(env)

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
				PlatformCtx.ConfigManager.SetConfigValue(key, val)
				cmd.Printf(emoji.Sprintf(":rocket: Config value %s set to %s\n", key, val))
				newCmd := NewRootCmd(nil)
				newCmd.SetArgs([]string{"config", "list"})
				_ = cmd.Execute()
				return
			} else {
				cmd.Println("Value must be provided")
				return
			}
		}

		username := cmd.Flag("username").Value.String()
		if username != "" {
			PlatformCtx.ConfigManager.SetConfigValue("username", username)
		}

		password := cmd.Flag("password").Value.String()
		if password != "" {
			PlatformCtx.ConfigManager.SetConfigValue("password", password)
		}

		cmd.Printf("Config updated successfully: %s\n", cfgFile)
		newCmd := NewRootCmd(nil)
		newCmd.SetArgs([]string{"config", "list"})
		_ = cmd.Execute()
	},
}

func init() {
	configCmd.AddCommand(setCmd)

	setCmd.Flags().StringP("lifecycle", "l", "", "Lifecycle of Pixo Platform to use (dev, stage, prod)")
	setCmd.Flags().StringP("region", "r", "", "Region of Pixo Platform to use (na, saudi)")
	setCmd.Flags().StringP("username", "u", "", "Username to use for Pixo Platform")
	setCmd.Flags().StringP("password", "p", "", "Password to use for Pixo Platform")
	setCmd.Flags().StringP("key", "k", "", "Key of the config value to set")
	setCmd.Flags().StringP("val", "v", "", "Value of the config value to set")
}
