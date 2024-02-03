/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/loader"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
)

// createUserCmd represents the create user command
var createUserCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user",
	Long:  `Create a new user with the following command:`,
	RunE: func(cmd *cobra.Command, args []string) error {
		firstName, _ := Ctx.ConfigManager.GetConfigValueOrAskUser("first-name", cmd)
		lastName, _ := Ctx.ConfigManager.GetConfigValueOrAskUser("last-name", cmd)
		username, _ := Ctx.ConfigManager.GetConfigValueOrAskUser("username", cmd)
		orgID, _ := Ctx.ConfigManager.GetIntConfigValueOrAskUser("org-id", cmd)
		role, _ := Ctx.ConfigManager.GetConfigValueOrAskUser("role", cmd)
		password, ok := Ctx.ConfigManager.GetSensitiveConfigValueOrAskUser("password", cmd)
		if !ok {
			cmd.Println(emoji.Sprintf(":exclamation: Password not provided"))
			return nil
		}

		input := platform.User{
			FirstName: firstName,
			LastName:  lastName,
			Username:  username,
			OrgID:     orgID,
			Role:      role,
			Password:  password,
		}

		spinner := loader.NewSpinner(cmd.OutOrStdout())

		user, err := Ctx.PlatformClient.CreateUser(cmd.Context(), input)
		if err != nil {
			cmd.Println("Could not create user: ", err.Error())
			return err
		}

		spinner.Stop()
		cmd.Println("Created user " + user.Username)
		return nil
	},
}

func init() {
	usersCmd.AddCommand(createUserCmd)

	createUserCmd.Flags().String("first-name", "", "First name of the user")
	createUserCmd.Flags().String("last-name", "", "Last name of the user")
	createUserCmd.Flags().String("username", "", "Username of the user")
	createUserCmd.Flags().String("password", "", "Password of the user")
	createUserCmd.Flags().String("org-id", "", "Organization ID of the user")
	createUserCmd.Flags().String("role", "user", "Role of the user")
}
