/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/spf13/cobra"
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
			Ctx.ConfigManager.Println(":file_folder: Opening config file in editor")
			if err := Ctx.FileOpener.OpenEditor(Ctx.ConfigManager.ConfigFile()); err != nil {
				Ctx.ConfigManager.Println(":warning: Unable to open editor: ", err)
			}
		}

		Ctx.ConfigManager.Println(":file_folder: Config: ", Ctx.ConfigManager.ConfigFile())

		if region := Ctx.ConfigManager.Region(); region != "" {
			Ctx.ConfigManager.Println(":earth_americas: Region: ", region)
		}

		if lifecycle := Ctx.ConfigManager.Lifecycle(); lifecycle != "" {
			Ctx.ConfigManager.Println(":gear:  Lifecycle: ", lifecycle)
		}

		Ctx.ConfigManager.Println()

		if userID, ok := Ctx.ConfigManager.GetConfigValue("user-id"); ok {
			Ctx.ConfigManager.Println(":id: User ID: ", userID)
		}

		if username, ok := Ctx.ConfigManager.GetConfigValue("username"); ok {
			Ctx.ConfigManager.Println(":bust_in_silhouette: Username: ", username)
		}

		if _, ok := Ctx.ConfigManager.GetConfigValue("password"); ok {
			Ctx.ConfigManager.Println(":lock: Password: ********")
		}

		if _, ok := Ctx.ConfigManager.GetConfigValue("api-key"); ok {
			Ctx.ConfigManager.Println(":key: API Key: ********")
		}

		if _, ok := Ctx.ConfigManager.GetConfigValue("token"); ok {
			Ctx.ConfigManager.Println(":key: Token: ********")
		}

		Ctx.ConfigManager.Println()

		activeEnv := Ctx.ConfigManager.GetActiveEnv()

		userInfoList := []string{
			"username",
			"user-id",
			"api-key",
		}
		sensitiveList := []string{
			"password",
			"token",
			"api-key",
		}
		isSensitiveOrRepetitive := func(k string) bool {
			list := append(userInfoList, sensitiveList...)
			for _, s := range list {
				if strings.Contains(k, s) {
					return true
				}
			}

			return false
		}

		cleanKey := func(k string) string {
			k = strings.Replace(k, "id", "ID", -1)
			k = strings.Replace(k, "api", "API", -1)
			k = strings.Replace(k, "-", " ", -1)
			return strings.Title(k)
		}

		for k, v := range activeEnv.EnvMap {
			if isSensitiveOrRepetitive(k) {
				continue
			}
			Ctx.ConfigManager.Println(":arrow_right: ", cleanKey(k), ": ", v)
		}

	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.PersistentFlags().BoolVarP(&edit, "edit", "e", false, "Edit the config file in your default editor")
}
