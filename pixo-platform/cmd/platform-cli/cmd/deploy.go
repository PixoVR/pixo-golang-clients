/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/parser"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/loader"
	platformAPI "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
	"os"
)

// deployCmd represents the deploy rootCmd
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a multiplayer server version",
	Long:  `Deploy a docker image as a multiplayer server version on the Pixo Platform for a module`,
	RunE: func(cmd *cobra.Command, args []string) error {

		moduleID := input.GetIntValueOrAskUser(cmd, "module-id", "MODULE_ID")
		if moduleID == 0 {
			cmd.Println(emoji.Sprint(":exclamation_mark:No module id provided"))
			return nil
		}

		semanticVersion := input.GetStringValueOrAskUser(cmd, "server-version", "SERVER_VERSION")
		if semanticVersion == "" {
			iniPath := cmd.Flag("ini").Value.String()

			iniParser, err := parser.NewIniParser(&iniPath)
			if err != nil {
				msg := emoji.Sprintf(":exclamation_mark:failed to parse ini file %s", iniPath)
				return errors.New(msg)
			}

			semanticVersion, err = iniParser.ParseSemanticVersion()
			if err != nil {
				msg := emoji.Sprintf(":exclamation_mark:No semantic version given and failed to parse server version from ini file %s", iniPath)
				return errors.New(msg)
			}

		}

		apiClient := getAuthenticatedClient()

		isPrecheck := cmd.Flag("pre-check").Value.String()
		if isPrecheck == "true" {

			params := platformAPI.MultiplayerServerVersionQueryParams{
				ModuleID:        moduleID,
				SemanticVersion: semanticVersion,
			}

			spinner := loader.NewSpinner(cmd.OutOrStdout())
			defer spinner.Stop()

			if versions, err := apiClient.GetMultiplayerServerVersions(cmd.Context(), params); err != nil {
				cmd.Println(emoji.Sprint(":negative_squared_cross_mark:Unable to retrieve server versions from platform api"))
				return err

			} else if len(versions) > 0 {
				spinner.Stop()
				cmd.Println(emoji.Sprintf(":red_square:Server version %s already exists\n", semanticVersion))
				os.Exit(1)
				return nil
			}

			cmd.Println(emoji.Sprintf(":heavy_check_mark:Server version does not exist yet: %s", semanticVersion))
			return nil
		}

		image := input.GetStringValueOrAskUser(cmd, "image", "GAMESERVER_IMAGE")
		if image == "" {
			return errors.New("no gameserver image provided")
		}

		spinner := loader.NewSpinner(cmd.OutOrStdout())

		if err := apiClient.CreateMultiplayerServerVersion(cmd.Context(), moduleID, image, semanticVersion); err != nil {
			msg := fmt.Sprintf("Failed to create multiplayer server version: %s - %s", semanticVersion, err.Error())
			return errors.New(msg)
		}

		spinner.Stop()
		cmd.Println(emoji.Sprint(":cruise_ship:Successfully created multiplayer server version: ", semanticVersion))
		return nil
	},
}

func init() {
	serverVersionsCmd.AddCommand(deployCmd)

	deployCmd.PersistentFlags().StringP("image", "i", "", "Docker image to deploy as the multiplayer server version")
	deployCmd.Flags().BoolP("pre-check", "p", false, "Check if server version exists already")
}
