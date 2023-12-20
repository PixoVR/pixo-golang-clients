/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/parser"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	platformAPI "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/spf13/cobra"
	"os"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a multiplayer server version",
	Long:  `Deploy a new image as a multiplayer server version on the Pixo Platform for a specific module`,
	Run: func(cmd *cobra.Command, args []string) {
		initLogger(cmd)

		moduleID := input.GetIntValueOrAskUser(cmd, "module-id", "MODULE_ID")
		if moduleID == 0 {
			cmd.Println("No module ID provided")
			os.Exit(1)
		}

		semanticVersion := input.GetStringValueOrAskUser(cmd, "server-version", "SERVER_VERSION")
		if semanticVersion == "" {
			iniPath := cmd.Flag("ini").Value.String()

			iniParser, err := parser.NewIniParser(&iniPath)
			if err != nil {
				cmd.Printf("Failed to parse ini file %s", iniPath)
				os.Exit(1)
			}

			semanticVersion, err = iniParser.ParseSemanticVersion()
			if err != nil {
				cmd.Printf("No semantic version given and failed to parse server version from ini file %s\n", iniPath)
				os.Exit(1)
			}

		}

		if cmd.Flag("pre-check").Value.String() == "true" {
			params := platformAPI.MultiplayerServerVersionQueryParams{
				ModuleID:        moduleID,
				SemanticVersion: semanticVersion,
			}
			if versions, err := apiClient.GetMultiplayerServerVersions(params); err != nil {
				cmd.Println("unable to retrieve from platform api")
				os.Exit(1)
			} else if len(versions) > 0 {
				cmd.Printf("server version %s already exists\n", semanticVersion)
				os.Exit(1)
			}

			cmd.Printf("server version %s does not exist yet\n", semanticVersion)
			return
		}

		image := input.GetStringValueOrAskUser(cmd, "image", "GAMESERVER_IMAGE")
		if image == "" {
			cmd.Println("No gameserver image provided")
			os.Exit(1)
		}

		if err := apiClient.CreateMultiplayerServerVersion(moduleID, image, semanticVersion); err != nil {
			cmd.Printf("Failed to create multiplayer server version: %s\n", semanticVersion)
			os.Exit(1)
		}

		cmd.Println("Successfully created multiplayer server version: ", semanticVersion)
	},
}

func init() {
	serverVersionsCmd.AddCommand(deployCmd)

	deployCmd.PersistentFlags().StringP("image", "i", "", "Docker image to deploy as the multiplayer server version")
	deployCmd.Flags().BoolP("pre-check", "p", false, "Check if server version exists already")
}
