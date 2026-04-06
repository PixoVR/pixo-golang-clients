/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
			Ctx.Println(":file_folder: Opening config file in editor")
			if err := Ctx.FileOpener.OpenEditor(activeConfigFile); err != nil {
				Ctx.Println(":warning: Unable to open editor: ", err)
			}
		}

		Ctx.Println(":file_folder: Config: ", activeConfigFile)

		if region := Ctx.ConfigManager.Region(); region != "" {
			Ctx.Println(":earth_americas: Region: ", region)
		}

		if lifecycle := Ctx.ConfigManager.Lifecycle(); lifecycle != "" {
			Ctx.Println(":gear: Status: ", lifecycle)
		}

		Ctx.Println()

		if userID, ok := Ctx.ConfigManager.GetConfigValue("auth-user-id"); ok {
			Ctx.Println(":id: User ID: ", userID)
		}

		if username, ok := Ctx.ConfigManager.GetConfigValue("auth-username"); ok {
			Ctx.Println(":bust_in_silhouette: Username: ", username)
		}

		if _, ok := Ctx.ConfigManager.GetConfigValue("auth-password"); ok {
			Ctx.Println(":lock: Password: ********")
		}

		if _, ok := Ctx.ConfigManager.GetConfigValue("api-key"); ok {
			Ctx.Println(":key: API Key: ********")
		}

		if _, ok := Ctx.ConfigManager.GetConfigValue("auth-token"); ok {
			Ctx.Println(":coin: Token: ********")
		}

		Ctx.Println()

		activeEnv := Ctx.ConfigManager.ActiveEnv()

		for k, v := range activeEnv.EnvMap {
			if isSensitiveOrRepetitive(k) {
				continue
			}
			Ctx.Println(":arrow_right: ", cleanKey(k), ": ", v)
		}

	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.PersistentFlags().BoolVarP(&edit, "edit", "e", false, "Edit the config file in your default editor")
}

func cleanKey(k string) string {
	k = strings.ReplaceAll(k, "id", "ID")
	k = strings.ReplaceAll(k, "api", "API")
	k = strings.ReplaceAll(k, "-", " ")

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
