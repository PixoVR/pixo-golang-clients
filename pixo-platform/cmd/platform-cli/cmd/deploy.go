/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/parser"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	platformAPI "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy rootCmd
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a multiplayer server version",
	Long:  `Deploy a docker image as a multiplayer server version on the Pixo Platform for a module`,
	RunE: func(cmd *cobra.Command, args []string) error {

		moduleID := input.GetIntValueOrAskUser(cmd, "module-id", "MODULE_ID")
		if moduleID == 0 {
			return errors.New("no module id provided")
		}

		semanticVersion := input.GetStringValueOrAskUser(cmd, "server-version", "SERVER_VERSION")
		if semanticVersion == "" {
			iniPath := cmd.Flag("ini").Value.String()

			iniParser, err := parser.NewIniParser(&iniPath)
			if err != nil {
				msg := fmt.Sprintf("failed to parse ini file %s", iniPath)
				return errors.New(msg)
			}

			semanticVersion, err = iniParser.ParseSemanticVersion()
			if err != nil {
				msg := fmt.Sprintf("no semantic version given and failed to parse server version from ini file %s", iniPath)
				return errors.New(msg)
			}

		}

		isPrecheck := cmd.Flag("pre-check").Value.String()
		if isPrecheck == "true" {

			params := platformAPI.MultiplayerServerVersionQueryParams{
				ModuleID:        moduleID,
				SemanticVersion: semanticVersion,
			}

			if versions, err := apiClient.GetMultiplayerServerVersions(context.Background(), params); err != nil {
				cmd.Println("unable to retrieve server versions from platform api")
				return err

			} else if len(versions) > 0 {
				errMsg := fmt.Sprintf("server version %s already exists\n", semanticVersion)
				return errors.New(errMsg)
			}

			cmd.Println("server version does not exist yet:", semanticVersion)
			return nil
		}

		image := input.GetStringValueOrAskUser(cmd, "image", "GAMESERVER_IMAGE")
		if image == "" {
			return errors.New("no gameserver image provided")
		}

		if err := apiClient.CreateMultiplayerServerVersion(context.Background(), moduleID, image, semanticVersion); err != nil {
			msg := fmt.Sprintf("Failed to create multiplayer server version: %s - %s", semanticVersion, err.Error())
			return errors.New(msg)
		}

		cmd.Println("Successfully created multiplayer server version:", semanticVersion)
		log.Info().Msgf("Successfully created multiplayer server version: %s", semanticVersion)
		return nil
	},
}

func init() {
	serverVersionsCmd.AddCommand(deployCmd)

	deployCmd.PersistentFlags().StringP("image", "i", "", "Docker image to deploy as the multiplayer server version")
	deployCmd.Flags().BoolP("pre-check", "p", false, "Check if server version exists already")
}
