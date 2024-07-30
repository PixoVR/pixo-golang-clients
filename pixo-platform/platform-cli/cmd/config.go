/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

var (
	edit bool
)

// configCmd represents the config rootCmd
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure the CLI settings",
	Long: `Manage settings like region, org, and module ID.  This commands will prompt you for the settings if they are not already set.
`,
	Run: func(cmd *cobra.Command, args []string) {

		if edit {
			Ctx.Printer.Println(":file_folder: Opening config file in editor")
			if err := Ctx.FileOpener.OpenEditor(activeConfigFile); err != nil {
				Ctx.Printer.Println(":warning: Unable to open editor: ", err)
			}
		}

		Ctx.Printer.Println(":file_folder: Config: ", activeConfigFile)

		if region := Ctx.ConfigManager.Region(); region != "" {
			Ctx.Printer.Println(":earth_americas: Region: ", region)
		}

		if lifecycle := Ctx.ConfigManager.Lifecycle(); lifecycle != "" {
			Ctx.Printer.Println(":gear: Lifecycle: ", lifecycle)
		}

		Ctx.Printer.Println()

		if userID, ok := Ctx.ConfigManager.GetConfigValue("auth-user-id"); ok {
			Ctx.Printer.Println(":id: User ID: ", userID)
		}

		if username, ok := Ctx.ConfigManager.GetConfigValue("auth-username"); ok {
			Ctx.Printer.Println(":bust_in_silhouette: Username: ", username)
		}

		if _, ok := Ctx.ConfigManager.GetConfigValue("auth-password"); ok {
			Ctx.Printer.Println(":lock: Password: ********")
		}

		if _, ok := Ctx.ConfigManager.GetConfigValue("api-key"); ok {
			Ctx.Printer.Println(":key: API Key: ********")
		}

		if _, ok := Ctx.ConfigManager.GetConfigValue("auth-token"); ok {
			Ctx.Printer.Println(":coin: Token: ********")
		}

		Ctx.Printer.Println()

		activeEnv := Ctx.ConfigManager.ActiveEnv()

		for k, v := range activeEnv.EnvMap {
			if isSensitiveOrRepetitive(k) {
				continue
			}
			Ctx.Printer.Println(":arrow_right: ", cleanKey(k), ": ", v)
		}

	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.PersistentFlags().BoolVarP(&edit, "edit", "e", false, "Edit the config file in your default editor")
}

func cleanKey(k string) string {
	k = strings.Replace(k, "id", "ID", -1)
	k = strings.Replace(k, "api", "API", -1)
	k = strings.Replace(k, "-", " ", -1)

	c := cases.Title(language.English)
	return c.String(k)
}

var (
	userInfoList = []string{
		"username",
		"user-id",
		"api-key",
	}
	sensitiveList = []string{
		"password",
		"token",
		"api-key",
	}
)

func isSensitiveOrRepetitive(k string) bool {
	list := append(userInfoList, sensitiveList...)
	for _, s := range list {
		if strings.Contains(k, s) {
			return true
		}
	}

	return false
}
