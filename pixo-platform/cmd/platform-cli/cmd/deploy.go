/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/parser"
	platformAPI "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strconv"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a multiplayer server version",
	Long:  `Deploy a new image as a multiplayer server version on the Pixo Platform for a specific module`,
	Run: func(cmd *cobra.Command, args []string) {
		moduleIDVal := cmd.Flag("module-id").Value.String()
		if moduleIDVal == "" {
			log.Panic().Msg("No module ID given")
		}
		moduleID, err := strconv.Atoi(moduleIDVal)
		if err != nil {
			log.Panic().Err(err).Msgf("Failed to parse module ID: %s", moduleIDVal)
		}

		semanticVersion := cmd.Flag("version").Value.String()
		if semanticVersion == "" {
			iniPath := cmd.Flag("ini").Value.String()

			iniParser, err := parser.NewIniParser(&iniPath)
			if err != nil {
				log.Panic().Err(err).Msg("Failed to create ini parser")
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

			log.Info().Msgf("server version %s does not exist yet", semanticVersion)
			return
		}

		image := cmd.Flag("image").Value.String()
		if image == "" {
			log.Panic().Msg("No image given")
		}

		if err := apiClient.CreateMultiplayerServerVersion(moduleID, image, semanticVersion); err != nil {
			log.Fatal().Err(err).Msgf("Failed to create multiplayer server version: %s", semanticVersion)
		}
	},
}

func init() {
	serverVersionsCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	deployCmd.PersistentFlags().StringP("image", "i", "", "Docker image to deploy as the multiplayer server version")
	deployCmd.PersistentFlags().StringP("version", "v", "1.00.00", "Semantic Version of the multiplayer server version")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	deployCmd.Flags().BoolP("pre-check", "p", false, "Check if server version exists already")
}
