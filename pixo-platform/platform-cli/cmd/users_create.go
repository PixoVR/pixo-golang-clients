/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"errors"
	"fmt"
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
		questions := []config.Value{
			{Question: forms.Question{Type: forms.Input, Key: "first-name"}},
			{Question: forms.Question{Type: forms.Input, Key: "last-name"}},
			{Question: forms.Question{Type: forms.Input, Key: "user-email", Optional: true}},
			{Question: forms.Question{Type: forms.Input, Key: "user-username", Optional: true}},
			{Question: forms.Question{Type: forms.SensitiveInput, Key: "user-password"}},
			{Question: forms.Question{
				Type: forms.SelectID,
				Key:  "role",
				GetOptionsFunc: func() ([]forms.Option, error) {
					items, err := Ctx.PlatformClient.GetRoles(cmd.Context())
					if err != nil {
						return nil, errors.New("unable to get roles")
					}

					options := make([]forms.Option, len(items))
					for i, item := range items {
						options[i] = forms.Option{
							Label: item.Name,
							Value: fmt.Sprint(item.ID),
						}
					}

					return options, nil
				},
			}},
			{Question: forms.Question{
				Type: forms.SelectID,
				Key:  "org",
				GetOptionsFunc: func() ([]forms.Option, error) {
					items, err := Ctx.PlatformClient.GetOrgs(cmd.Context())
					if err != nil {
						return nil, errors.New("unable to get orgs")
					}

					options := make([]forms.Option, len(items))
					for i, item := range items {
						labelPrefix := fmt.Sprintf("Org ID %d: %s", item.ID, item.Name)
						options[i] = forms.Option{
							Label: labelPrefix,
							Value: fmt.Sprint(item.ID),
						}
					}

					return options, nil
				},
			}},
		}

		answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd)
		if err != nil {
			return err
		}

		user := &platform.User{
			FirstName: forms.String(answers["first-name"]),
			LastName:  forms.String(answers["last-name"]),
			Email:     forms.String(answers["user-email"]),
			Username:  forms.String(answers["user-username"]),
			OrgID:     forms.Int(answers["org"]),
			Role:      forms.String(answers["role"]),
			Password:  forms.String(answers["user-password"]),
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
