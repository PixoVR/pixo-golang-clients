/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"context"
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
			{Question: forms.Question{Type: forms.Input, Key: "email", Optional: true}},
			{Question: forms.Question{Type: forms.Input, Key: "username", Optional: true}},
			{Question: forms.Question{Type: forms.SensitiveInput, Key: "password"}},
			{Question: forms.Question{Type: forms.Select, Key: "role",
				GetItemsFunc: func(ctx context.Context) (interface{}, error) {
					return Ctx.PlatformClient.GetRoles(cmd.Context())
				},
			}},
			{Question: forms.Question{
				Type: forms.SelectID, Key: "org",
				LabelFunc: func(item interface{}) string {
					org := item.(platform.Org)
					return fmt.Sprintf("Org ID %d: %s", org.ID, org.Name)
				},
				GetItemsFunc: func(ctx context.Context) (interface{}, error) {
					return Ctx.PlatformClient.GetOrgs(cmd.Context())
				},
			}},
		}

		answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd)
		if err != nil {
			return err
		}

		user := &platform.User{
			OrgID:     forms.Int(answers["org"]),
			FirstName: forms.String(answers["first-name"]),
			LastName:  forms.String(answers["last-name"]),
			Email:     forms.String(answers["email"]),
			Username:  forms.String(answers["username"]),
			Password:  forms.String(answers["password"]),
			Role:      forms.String(answers["role"]),
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
	createUserCmd.Flags().String("email", "", "Email of the new user")
	createUserCmd.Flags().String("username", "", "Username of the new user")
	createUserCmd.Flags().String("password", "", "Password of the new user")
	createUserCmd.Flags().String("org", "", "Organization ID of the new user")
	createUserCmd.Flags().String("role", "", "Role of the new user")
}
