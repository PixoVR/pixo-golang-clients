/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/parser"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	platformAPI "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a multiplayer server version",
	Long:  `Deploy a new image as a multiplayer server version on the Pixo Platform for a specific module`,
	Run: func(cmd *cobra.Command, args []string) {
		moduleID := input.GetIntValue(cmd, "module-id", "MODULE_ID")
		if moduleID == 0 {
			cmd.Println("No module ID provided")
			return
		}

		semanticVersion := input.GetStringValue(cmd, "server-version", "SERVER_VERSION")
		if semanticVersion == "" {
			iniPath := cmd.Flag("ini").Value.String()

			iniParser, err := parser.NewIniParser(&iniPath)
			if err != nil {
				log.Debug().Err(err).Msgf("Failed to parse ini file %s", iniPath)
				return
			}

			semanticVersion, err = iniParser.ParseSemanticVersion()
			if err != nil {
				log.Fatal().Err(err).Msgf("No semantic version given and failed to parse server version from ini file %s", iniPath)
			}

		}

		if cmd.Flag("pre-check").Value.String() == "true" {
			params := platformAPI.MultiplayerServerVersionQueryParams{
				ModuleID:        moduleID,
				SemanticVersion: semanticVersion,
			}
			if versions, err := apiClient.GetMultiplayerServerVersions(params); err != nil {
				log.Fatal().Err(err).Msgf("unable to retrieve from platform api")
			} else if len(versions) > 0 {
				log.Fatal().Msgf("server version %s already exists", semanticVersion)
			}

			cmd.Println(fmt.Sprintf("server version %s does not exist yet", semanticVersion))
			return
		}

		image := input.GetStringValue(cmd, "image", "GAMESERVER_IMAGE")
		if image == "" {
			log.Fatal().Msg("No gameserver image provided")
		}

		if err := apiClient.CreateMultiplayerServerVersion(moduleID, image, semanticVersion); err != nil {
			log.Fatal().Err(err).Msgf("Failed to create multiplayer server version: %s", semanticVersion)
		}

		cmd.Println("Successfully created multiplayer server version: ", semanticVersion)
	},
}

func init() {
	serverVersionsCmd.AddCommand(deployCmd)

	deployCmd.PersistentFlags().StringP("image", "i", "", "Docker image to deploy as the multiplayer server version")
	deployCmd.PersistentFlags().StringP("server-version", "v", "", "Semantic Version of the multiplayer server version")

	deployCmd.Flags().BoolP("pre-check", "p", false, "Check if server version exists already")
}
