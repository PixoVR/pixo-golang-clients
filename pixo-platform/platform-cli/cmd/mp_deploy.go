/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/loader"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/parser"
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
)

var (
	isPrecheck bool
	isUpdate   bool
)

// mpDeployCmd represents the mp server versions deploy command
var mpDeployCmd = &cobra.Command{
	Use:           "deploy",
	Short:         "Deploy a multiplayer server version",
	Long:          `Deploy a docker image as a multiplayer server version on the Pixo Platform for a module`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		questions := []config.Value{
			{Question: moduleQuestion()},
			{Question: forms.Question{Type: forms.Input, Key: "server-version", Optional: true}},
		}

		answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd)
		if err != nil {
			return err
		}

		moduleID := forms.Int(answers["module"])

		semVer := forms.String(answers["server-version"])
		if semVer == "" {
			if iniPath, ok := Ctx.ConfigManager.GetFlagOrConfigValue("ini", cmd); ok && iniPath != "" {
				iniParser, err := parser.NewIniParser(&iniPath)
				if err != nil {
					msg := emoji.Sprintf(":exclamation_mark: failed to parse ini file %s", iniPath)
					Ctx.Printer.Println(msg)
					return errors.New(msg)
				}

				semVer, err = iniParser.ParseSemanticVersion()
				if err != nil {
					msg := emoji.Sprintf(":exclamation_mark: No semantic version given and failed to parse server version from ini file %s", iniPath)
					Ctx.Printer.Println(msg)
					return errors.New(msg)
				}

			}
		}

		if semVer == "" {
			return errors.New("SERVER VERSION not provided")
		}

		if isPrecheck {
			params := &platform.MultiplayerServerVersionParams{
				ModuleID:        moduleID,
				SemanticVersion: semVer,
			}

			spinner := loader.NewLoader(cmd.Context(), "Getting multiplayer server versions...", Ctx.Printer)

			if versions, err := Ctx.PlatformClient.GetMultiplayerServerVersionsWithConfig(cmd.Context(), params); err != nil {
				Ctx.Printer.Println(":negative_squared_cross_mark: Unable to retrieve server versions from the Pixo Platform")
				spinner.Stop()
				return err

			} else if len(versions) > 0 {
				spinner.Stop()
				msg := emoji.Sprintf(":exclamation: server version %s already exists\n", semVer)
				Ctx.Printer.Println(msg)
				return errors.New(msg)
			}

			spinner.Stop()
			Ctx.Printer.Println(":heavy_check_mark: Server version does not exist yet: ", semVer)
			return nil
		}

		var image string

		filePath, ok := Ctx.ConfigManager.GetFlagOrConfigValue("zip-file", cmd)
		if !ok || filePath == "" {
			question := forms.Question{Type: forms.Input, Key: "image"}
			answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm([]config.Value{{Question: question}}, cmd)
			if err != nil {
				return err
			}
			image = forms.String(answers["image"])
		}

		input := platform.MultiplayerServerVersion{
			ModuleID:        moduleID,
			SemanticVersion: semVer,
			ImageRegistry:   image,
			LocalFilePath:   filePath,
			Engine:          "unreal",
		}

		if isUpdate {
			msg := fmt.Sprint("Updating server version: ", semVer)
			spinner := loader.NewLoader(cmd.Context(), msg, Ctx.Printer)
			serverVersion, err := Ctx.PlatformClient.UpdateMultiplayerServerVersion(cmd.Context(), platform.MultiplayerServerVersion{
				ModuleID:        moduleID,
				SemanticVersion: semVer,
				ImageRegistry:   image,
			})
			spinner.Stop()
			if err != nil {
				msg := fmt.Sprintf("Failed to update multiplayer server version: %s - %s", semVer, err.Error())
				return errors.New(msg)
			}
			Ctx.Printer.Printf(":cruise_ship: Updated server version: %s - %s\n", serverVersion.Module.Abbreviation, semVer)
			return nil
		} else {
			msg := fmt.Sprint("Deploying server version: ", semVer)
			spinner := loader.NewLoader(cmd.Context(), msg, Ctx.Printer)
			serverVersion, err := Ctx.PlatformClient.CreateMultiplayerServerVersion(cmd.Context(), input)
			spinner.Stop()
			if err != nil {
				msg := fmt.Sprintf("Failed to deploy multiplayer server version: %s - %s", semVer, err.Error())
				return errors.New(msg)
			}
			Ctx.Printer.Printf(":cruise_ship: Deployed version: %s - %s\n", serverVersion.Module.Abbreviation, semVer)
			return nil
		}

	},
}

func init() {
	serversCmd.AddCommand(mpDeployCmd)

	mpDeployCmd.PersistentFlags().StringP("image", "i", "", "Docker image to deploy as the multiplayer server version")
	mpDeployCmd.Flags().StringP("ini", "f", "", "Path to the ini file to use for the semantic version")
	mpDeployCmd.Flags().StringP("zip-file", "z", "", "Path to the zip file to use for the upload")
	mpDeployCmd.Flags().BoolVarP(&isPrecheck, "pre-check", "p", false, "Check if server version exists already")
	mpDeployCmd.Flags().BoolVarP(&isUpdate, "update", "u", false, "Update an existing server version")
}
