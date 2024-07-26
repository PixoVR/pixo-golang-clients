/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/legacy"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/loader"
	"github.com/spf13/cobra"
)

// createUserCmd represents the create user command
var createUserCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user",
	Long:  `Create a new user with the following command:`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		authUsername, _ := Ctx.ConfigManager.GetConfigValue("username")
		authPassword, _ := Ctx.ConfigManager.GetConfigValue("password")
		if err := Ctx.LegacyClient.Login(authUsername, authPassword); err != nil {
			Ctx.Printer.Println(":exclamation: Unable to login")
			return err
		}

		questions := []config.Value{
			{Question: forms.Question{Type: forms.Input, Key: "first-name"}},
			{Question: forms.Question{Type: forms.Input, Key: "last-name"}},
			{Question: forms.Question{Type: forms.Input, Key: "user-email"}},
			{Question: forms.Question{Type: forms.Input, Key: "user-username"}},
			{Question: forms.Question{Type: forms.SensitiveInput, Key: "user-password"}},
		}

		var orgs []legacy.Org
		if _, ok := Ctx.ConfigManager.GetFlagOrConfigValue("org", cmd); !ok {
			orgs, err = Ctx.LegacyClient.GetOrgs()
			if err != nil {
				Ctx.Printer.Println(":exclamation: Unable to get orgs")
				return err
			}
		}

		orgOptions := make([]forms.Option, len(orgs))
		for i, org := range orgs {
			labelPrefix := fmt.Sprintf("Org ID %d: %s", org.ID, org.Name)
			orgOptions[i] = forms.Option{
				Label: labelPrefix,
				Value: fmt.Sprint(org.ID),
			}
		}
		questions = append(questions, config.Value{
			Question: forms.Question{Type: forms.SelectID, Key: "org", Options: orgOptions},
		})

		var roles []platform.Role
		if _, ok := Ctx.ConfigManager.GetFlagOrConfigValue("role", cmd); !ok {
			roles, err = Ctx.PlatformClient.GetRoles(cmd.Context())
			if err != nil {
				Ctx.Printer.Println(":exclamation: Unable to get roles")
				return err
			}
		}

		roleOptions := make([]forms.Option, len(roles))
		for i, role := range roles {
			roleOptions[i] = forms.Option{
				Label: role.Name,
			}
		}
		questions = append(questions, config.Value{
			Question: forms.Question{Type: forms.Select, Key: "role", Options: roleOptions},
		})

		answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd)
		if err != nil {
			Ctx.Printer.Printf(":exclamation: %v\n", err)
		}

		firstName := forms.String(answers["first-name"])
		lastName := forms.String(answers["last-name"])
		email := forms.String(answers["user-email"])
		username := forms.String(answers["user-username"])
		password := forms.String(answers["user-password"])
		orgID := forms.Int(answers["org"])
		role := forms.String(answers["role"])

		user := &platform.User{
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Username:  username,
			OrgID:     orgID,
			Role:      role,
			Password:  password,
		}

		spinner := loader.NewLoader(cmd.Context(), "Creating user...", Ctx.Printer)
		err = Ctx.PlatformClient.CreateUser(cmd.Context(), user)
		spinner.Stop()
		if err != nil {
			Ctx.Printer.Println(":exclamation: Unable to create user: ", err)
			return err
		}

		Ctx.Printer.Printf(":rocket: User created: %s - %s\n", user.Email, user.Role)
		return nil
	},
}

func init() {
	usersCmd.AddCommand(createUserCmd)

	createUserCmd.Flags().String("first-name", "", "First name of the new user")
	createUserCmd.Flags().String("last-name", "", "Last name of the new user")
	createUserCmd.Flags().String("user-email", "", "Email of the new user")
	createUserCmd.Flags().String("user-username", "", "Username of the new user")
	createUserCmd.Flags().String("user-password", "", "Password of the new user")
	createUserCmd.Flags().String("org", "", "Organization ID of the new user")
	createUserCmd.Flags().String("role", "", "Role of the new user")
}
