/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"context"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	graphql_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strconv"
)

// userCmd represents the create command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Creating a user",
	Long:  `Creating a user with the following command:`,
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

		user, err := client.CreateUser(context.Background(), input)
		if err != nil {
			cmd.Println("Could not create user: ", err.Error())
			return err
		}

		cmd.Println("Created user " + user.Username)
		return nil
	},
}

func init() {
	createCmd.AddCommand(userCmd)

	userCmd.Flags().String("first-name", "", "First name of the user")
	userCmd.Flags().String("last-name", "", "Last name of the user")
	userCmd.Flags().String("username", "", "Username of the user")
	userCmd.Flags().String("password", "", "Password of the user")
	userCmd.Flags().String("org-id", "", "Organization ID of the user")
	userCmd.Flags().String("role", "user", "Role of the user")
}
