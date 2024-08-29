/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"context"
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
			{Question: forms.Question{
				Type: forms.SelectID,
				Key:  "module-id",
				LabelFunc: func(item interface{}) string {
					module := item.(platform.Module)
					return fmt.Sprintf("%d: %s - %s", module.ID, module.Abbreviation, module.Name)
				},
				GetItemsFunc: func(ctx context.Context) (interface{}, error) {
					return Ctx.PlatformClient.GetModules(cmd.Context())
				},
			}},
			{Question: forms.Question{Type: forms.Input, Key: "server-version", Optional: true}},
		}

		answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd)
		if err != nil {
			return err
		}

		moduleID := forms.Int(answers["module-id"])

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
				msg := emoji.Sprintf(":exclamation: server version %s already exists", semVer)
				Ctx.Printer.Println(msg)
				return errors.New(msg)
			}

			spinner.Stop()
			Ctx.Printer.Println(":heavy_check_mark: Server version does not exist yet: ", semVer)
			return nil
		}

		var filePath string
		image, ok := Ctx.ConfigManager.GetFlagOrConfigValue("image", cmd)
		if !ok || image == "" {
			filePath, ok = Ctx.ConfigManager.GetFlagOrConfigValue("zip-file", cmd)
			if !ok || filePath == "" {
				question := &forms.Question{Prompt: "DOCKER IMAGE"}
				if err := Ctx.FormHandler.GetResponseFromUser(question); err != nil {
					return err
				}
				if question.Answer == "" {
					return errors.New("no image or zip file provided")
				}
				image = forms.String(question.Answer)
			}
		}

		msg := fmt.Sprint("Deploying server version: ", semVer)
		spinner := loader.NewLoader(cmd.Context(), msg, Ctx.Printer)

		input := platform.MultiplayerServerVersion{
			ModuleID:        moduleID,
			ImageRegistry:   image,
			LocalFilePath:   filePath,
			SemanticVersion: semVer,
			Engine:          "unreal",
		}
		if _, err := Ctx.PlatformClient.CreateMultiplayerServerVersion(cmd.Context(), input); err != nil {
			msg := fmt.Sprintf("Failed to deploy multiplayer server version: %s - %s", semVer, err.Error())
			return errors.New(msg)
		}

		spinner.Stop()
		Ctx.Printer.Println(":cruise_ship: Deployed version: ", semVer)
		return nil
	},
}

func init() {
	serversCmd.AddCommand(mpDeployCmd)

	mpDeployCmd.PersistentFlags().StringP("image", "i", "", "Docker image to deploy as the multiplayer server version")
	mpDeployCmd.Flags().StringP("ini", "f", "", "Path to the ini file to use for the semantic version")
	mpDeployCmd.Flags().StringP("zip-file", "z", "", "Path to the zip file to use for the upload")
	mpDeployCmd.Flags().BoolVarP(&isPrecheck, "pre-check", "p", false, "Check if server version exists already")
}
