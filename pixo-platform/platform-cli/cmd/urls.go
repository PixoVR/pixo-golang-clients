/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/spf13/cobra"
	"strings"
)

// urlsCmd represents the urls command
var urlsCmd = &cobra.Command{
	Use:   "urls",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		region := Ctx.ConfigManager.Region()
		if region != "" {
			Ctx.ConfigManager.Println(":earth_americas: Region: ", region)
		}

		lifecycle := Ctx.ConfigManager.Lifecycle()
		if lifecycle != "" {
			Ctx.ConfigManager.Println(":gear:  Lifecycle: ", lifecycle)
		}

		Ctx.ConfigManager.Println()

		client := urlfinder.ServiceConfig{
			Region:    region,
			Lifecycle: lifecycle,
		}
		url := strings.Replace(client.FormatURL(), "/v2", "", 1)
		Ctx.ConfigManager.Println(":link: Web: ", url)
		Ctx.ConfigManager.Println("\n:link: Platform API: ", url, "/v2")
		Ctx.ConfigManager.Println(":link: Platform API Docs: ", url, "/v2/swagger/index.html")
		Ctx.ConfigManager.Println("\n:link: Matchmaking API: ", url, "/matchmaking")
		Ctx.ConfigManager.Println(":link: Matchmaking API Docs: ", url, "/matchmaking/swagger/index.html")
		Ctx.ConfigManager.Println("\n:link: Heartbeat API: ", url, "/heartbeat")
		Ctx.ConfigManager.Println(":link: Heartbeat API Docs: ", url, "/heartbeat/swagger/index.html")
	},
}

func init() {
	configCmd.AddCommand(urlsCmd)
}
