/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set	a config value",
	Long:  `Can set a single value, or change region/lifecycle of the Pixo Platform APIs used`,
	Run: func(cmd *cobra.Command, args []string) {

		if key, err := cmd.Flags().GetString("key"); err != nil {
			log.Error().Err(err).Msg("Could not get key flag")
		} else if key != "" {
			if val, err := cmd.Flags().GetString("val"); err != nil {
				log.Error().Err(err).Msg("Could not get val flag")
			} else if val != "" {
				viper.Set(key, val)
			} else {
				cmd.Println("Value must be provided")
			}
		}

		lifecycle := input.GetStringValueOrAskUser(cmd, "lifecycle", "PIXO_LIFECYCLE", "prod")
		region := input.GetStringValueOrAskUser(cmd, "region", "PIXO_REGION", "na")

		viper.Set("lifecycle", lifecycle)
		viper.Set("region", region)

		if err := viper.WriteConfigAs(cfgFile); err != nil {
			log.Error().Err(err).Msg("Could not write config file")
		} else {
			cmd.Println("Config file updated successfully")
		}
	},
}

func init() {
	configCmd.AddCommand(setCmd)

	setCmd.Flags().StringP("lifecycle", "l", "", "Lifecycle of Pixo Platform to use (dev, stage, prod)")
	setCmd.Flags().StringP("region", "r", "", "Region of Pixo Platform to use (na, saudi)")
	setCmd.Flags().StringP("key", "k", "", "Key of the config value to set")
	setCmd.Flags().StringP("val", "v", "", "Value of the config value to set")
}
