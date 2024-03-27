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

// mpDeployCmd represents the deploy rootCmd
var mpDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a multiplayer server version",
	Long:  `Deploy a docker image as a multiplayer server version on the Pixo Platform for a module`,
	RunE: func(cmd *cobra.Command, args []string) error {

		moduleID, ok := Ctx.ConfigManager.GetIntConfigValueOrAskUser("module-id", cmd)
		if !ok {
			Ctx.ConfigManager.Println(":warning: Module ID not provided")
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

		Ctx.ConfigManager.Println("Deploying server version: ", semanticVersion)

		if isPrecheck {

			params := &platformAPI.MultiplayerServerVersionQueryParams{
				ModuleID:        moduleID,
				SemanticVersion: semanticVersion,
			}

			spinner := loader.NewSpinner(Ctx.ConfigManager)

			if versions, err := Ctx.PlatformClient.GetMultiplayerServerVersions(cmd.Context(), params); err != nil {
				Ctx.ConfigManager.Println(":negative_squared_cross_mark: Unable to retrieve server versions from the Pixo Platform")
				return err

			} else if len(versions) > 0 {
				spinner.Stop()
				Ctx.ConfigManager.Printf(":exclamation: Server version %s already exists\n", semanticVersion)
				return errors.New("server version already exists")
			}

			spinner.Stop()
			Ctx.ConfigManager.Println(":heavy_check_mark: Server version does not exist yet: ", semanticVersion)
			return nil
		}

		var filePath string
		image, ok := Ctx.ConfigManager.GetFlagOrConfigValue("image", cmd)
		if !ok || image == "" {
			filePath, ok = Ctx.ConfigManager.GetFlagOrConfigValue("zip-file", cmd)
			if !ok || filePath == "" {
				return errors.New("no gameserver image or zip file provided")
			}
		}

		spinner := loader.NewSpinner(Ctx.ConfigManager)

		input := platformAPI.MultiplayerServerVersion{
			ModuleID:        moduleID,
			ImageRegistry:   image,
			LocalFilePath:   filePath,
			SemanticVersion: semanticVersion,
			Engine:          "unreal",
		}
		if _, err := Ctx.PlatformClient.CreateMultiplayerServerVersion(cmd.Context(), input); err != nil {
			msg := fmt.Sprintf("Failed to deploy multiplayer server version: %s - %s", semanticVersion, err.Error())
			return errors.New(msg)
		}

		spinner.Stop()
		Ctx.ConfigManager.Println(":cruise_ship: Deployed version: ", semanticVersion)
		return nil
	},
}

func init() {
	serverVersionsCmd.AddCommand(mpDeployCmd)

	mpDeployCmd.PersistentFlags().StringP("image", "i", "", "Docker image to deploy as the multiplayer server version")
	mpDeployCmd.Flags().StringP("ini", "f", parser.DefaultConfigFilepath, "Path to the ini file to use for the semantic version")
	mpDeployCmd.Flags().StringP("zip-file", "z", "", "Path to the zip file to use for the upload")
	mpDeployCmd.Flags().BoolVarP(&isPrecheck, "pre-check", "p", false, "Check if server version exists already")
}
