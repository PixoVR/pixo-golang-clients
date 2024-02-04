/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"errors"
	"fmt"
	platformAPI "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/parser"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/loader"
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
)

var (
	isPrecheck bool
)

// deployCmd represents the deploy rootCmd
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a multiplayer server version",
	Long:  `Deploy a docker image as a multiplayer server version on the Pixo Platform for a module`,
	RunE: func(cmd *cobra.Command, args []string) error {

		moduleID, ok := Ctx.ConfigManager.GetIntConfigValueOrAskUser("module-id", cmd)
		if !ok {
			cmd.Println(emoji.Sprintf(":warning: Module ID not provided"))
			return errors.New("module ID not provided")
		}

		semanticVersion, ok := Ctx.ConfigManager.GetConfigValueOrAskUser("server-version", cmd)
		if !ok {
			iniPath, _ := Ctx.ConfigManager.GetConfigValueOrAskUser("ini", cmd)
			iniParser, err := parser.NewIniParser(&iniPath)
			if err != nil {
				msg := emoji.Sprintf(":exclamation_mark: failed to parse ini file %s", iniPath)
				return errors.New(msg)
			}

			semanticVersion, err = iniParser.ParseSemanticVersion()
			if err != nil {
				msg := emoji.Sprintf(":exclamation_mark: No semantic version given and failed to parse server version from ini file %s", iniPath)
				return errors.New(msg)
			}

		}

		cmd.Println("Deploying server version: ", semanticVersion)

		if isPrecheck {

			params := &platformAPI.MultiplayerServerVersionQueryParams{
				ModuleID:        moduleID,
				SemanticVersion: semanticVersion,
			}

			spinner := loader.NewSpinner(cmd.OutOrStdout())

			if versions, err := Ctx.PlatformClient.GetMultiplayerServerVersions(cmd.Context(), params); err != nil {
				cmd.Println(emoji.Sprint(":negative_squared_cross_mark: Unable to retrieve server versions from platform api"))
				return err

			} else if len(versions) > 0 {
				spinner.Stop()
				cmd.Println(emoji.Sprintf(":red_square: Server version %s already exists\n", semanticVersion))
				return errors.New("server version already exists")
			}

			spinner.Stop()
			cmd.Println(emoji.Sprintf(":heavy_check_mark: Server version does not exist yet: %s", semanticVersion))
			return nil
		}

		image, ok := Ctx.ConfigManager.GetConfigValueOrAskUser("image", cmd)
		if !ok {
			return errors.New("no gameserver image provided")
		}

		spinner := loader.NewSpinner(cmd.OutOrStdout())

		if err := Ctx.PlatformClient.CreateMultiplayerServerVersion(cmd.Context(), moduleID, image, semanticVersion); err != nil {
			msg := fmt.Sprintf("Failed to create multiplayer server version: %s - %s", semanticVersion, err.Error())
			return errors.New(msg)
		}

		spinner.Stop()
		cmd.Println(emoji.Sprint(":cruise_ship: Successfully created multiplayer server version: ", semanticVersion))
		return nil
	},
}

func init() {
	serverVersionsCmd.AddCommand(deployCmd)

	deployCmd.PersistentFlags().StringP("image", "i", "", "Docker image to deploy as the multiplayer server version")
	deployCmd.Flags().StringP("ini", "f", parser.DefaultConfigFilepath, "Path to the ini file to use for the rootCmd")
	deployCmd.Flags().BoolVarP(&isPrecheck, "pre-check", "p", false, "Check if server version exists already")
}