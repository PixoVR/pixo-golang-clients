/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/loader"
	graphql_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/rs/zerolog/log"
	"strconv"

	"github.com/spf13/cobra"
)

// createUserCmd represents the create user command
var createUserCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user",
	Long:  `Create a new user with the following command:`,
	RunE: func(cmd *cobra.Command, args []string) error {
		clientConfig := urlfinder.ClientConfig{
			Lifecycle: input.GetConfigValue("lifecycle", "PIXO_LIFECYCLE"),
			Region:    input.GetConfigValue("region", "PIXO_REGION"),
		}
		client, err := graphql_api.NewClientWithBasicAuth(
			input.GetConfigValue("username", "PIXO_USERNAME"),
			input.GetConfigValue("password", "PIXO_PASSWORD"),
			clientConfig,
		)
		if err != nil {
			log.Error().Err(err).Msg("Could not create platform client")
			return err
		}

		firstName := input.GetStringValueOrAskUser(cmd, "first-name", "FIRST_NAME")
		lastName := input.GetStringValueOrAskUser(cmd, "last-name", "LAST_NAME")
		username := input.GetStringValueOrAskUser(cmd, "username", "USERNAME")
		orgIDVal := input.GetStringValueOrAskUser(cmd, "org-id", "ORG_ID")
		role := input.GetStringValueOrAskUser(cmd, "role", "ROLE")
		password := input.GetStringValueOrAskUser(cmd, "password", "PASSWORD")

		orgID, err := strconv.Atoi(orgIDVal)
		if err != nil {
			cmd.Println("Could not parse org id")
			return err
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

		user, err := client.CreateUser(cmd.Context(), input)
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
