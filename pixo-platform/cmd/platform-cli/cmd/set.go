/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
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
		lifecycle, _ := cmd.Flags().GetString("lifecycle")
		region, _ := cmd.Flags().GetString("region")

		viper.Set("lifecycle", lifecycle)
		viper.Set("region", region)

		config := urlfinder.ServiceConfig{
			Region:    region,
			Lifecycle: lifecycle,
		}

		config.Service = "api"
		viper.Set("legacy-api-url", config.FormatURL())

		config.Service = "primary"
		viper.Set("platform-api-url", config.FormatURL())

		config.Service = "match"
		viper.Set("matchmaking-api-url", config.FormatURL())

		if err := viper.WriteConfigAs(cfgFile); err != nil {
			log.Error().Err(err).Msg("Could not write config file")
		} else {
			cmd.Println("Config file updated successfully")
		}
	},
}

func init() {
	configCmd.AddCommand(setCmd)

	setCmd.Flags().StringP("lifecycle", "l", "dev", "Lifecycle of Pixo Platform to use")
	setCmd.Flags().StringP("region", "r", "na", "Region of Pixo Platform to use")
}
