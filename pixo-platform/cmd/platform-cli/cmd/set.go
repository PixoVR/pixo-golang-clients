/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	"github.com/kyokomi/emoji"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// setCmd represents the set rootCmd
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set	a config value",
	Long:  `Can set a single value, or change region/lifecycle of the Pixo Platform APIs used`,
	Run: func(cmd *cobra.Command, args []string) {

		if key, err := cmd.Flags().GetString("key"); err != nil {
			cmd.Println(emoji.Sprintf(":exclamation: Unable to get key flag"))
			return
		} else if key != "" {
			if val, err := cmd.Flags().GetString("val"); err != nil {
				cmd.Println(emoji.Sprintf(":exclamation: Unable to get value flag"))
				return
			} else if val != "" {
				if err := PlatformCtx.ConfigManager.SetConfigValue(key, val); err != nil {
					cmd.Println(emoji.Sprintf(":exclamation: Unable to write config value %s to config file", key))
					return
				} else {
					cmd.Printf(emoji.Sprintf(":rocket: GetConfig value %s set to %s\n", key, val))
					newCmd := NewRootCmd(nil)
					newCmd.SetArgs([]string{"config", "list"})
					_ = cmd.Execute()
					return
				}
			} else {
				cmd.Println("Value must be provided")
				return
			}
		}

		username := input.GetStringValueOrAskUser(cmd, "username", "PIXO_USERNAME", "")
		if err := PlatformCtx.ConfigManager.SetUsername(username); err != nil {
			cmd.Println(emoji.Sprintf(":x: Unable to write username to config file"))
		}

		password := input.GetSensitiveStringValueOrAskUser(cmd, "password", "PIXO_PASSWORD", "")
		if err := PlatformCtx.ConfigManager.SetPassword(password); err != nil {
			cmd.Println(emoji.Sprintf(":x: Unable to write password to config file"))
		}

		if username != "" && password != "" {
			_ = PlatformCtx.ConfigManager.SetUsername(username)
			_ = PlatformCtx.ConfigManager.SetPassword(password)
			if err := PlatformCtx.Authenticate(cmd.InOrStdin(), cmd.OutOrStdout()); err != nil {
				cmd.Println(emoji.Sprintf(":x: Login failed. Please check your credentials and try again."))
			}
		}

		lifecycle := input.GetStringValueOrAskUser(cmd, "lifecycle", "PIXO_LIFECYCLE", "prod")
		region := input.GetStringValueOrAskUser(cmd, "region", "PIXO_REGION", "na")

		env := config.Env{Lifecycle: lifecycle, Region: region}
		if err := PlatformCtx.ConfigManager.SetActiveEnv(env); err != nil {
			log.Error().Err(err).Msg("Could not write env to config file")
		} else {
			cmd.Printf("GetConfig file updated successfully: %s\n", cfgFile)
			newCmd := NewRootCmd(nil)
			newCmd.SetArgs([]string{"config", "list"})
			_ = cmd.Execute()
		}
	},
}

func init() {
	configCmd.AddCommand(setCmd)

	setCmd.Flags().StringP("lifecycle", "l", "prod", "Lifecycle of Pixo Platform to use (dev, stage, prod)")
	setCmd.Flags().StringP("region", "r", "na", "Region of Pixo Platform to use (na, saudi)")
	setCmd.Flags().StringP("username", "u", "", "Username to use for Pixo Platform")
	setCmd.Flags().StringP("password", "p", "", "Password to use for Pixo Platform")
	setCmd.Flags().StringP("key", "k", "", "Key of the config value to set")
	setCmd.Flags().StringP("val", "v", "", "Value of the config value to set")
}
