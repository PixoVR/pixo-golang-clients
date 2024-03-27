/*
Copyright © 2023 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"context"
	"fmt"
	graphql_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/forms"
	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

// modulesDeployCmd represents the modules deploy rootCmd
var modulesDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy module versions",
	Long: `Deploy module versions
`,
	Run: func(cmd *cobra.Command, args []string) {

		moduleID, ok := Ctx.ConfigManager.GetIntConfigValueOrAskUser("module-id", cmd)
		if !ok {
			Ctx.ConfigManager.Println(":exclamation: Module ID not provided")
			return
		}

		semVer, ok := Ctx.ConfigManager.GetConfigValueOrAskUser("semantic-version", cmd)
		if !ok {
			Ctx.ConfigManager.Println(":exclamation: Semantic version not provided")
			return
		}

		packageName, ok := Ctx.ConfigManager.GetConfigValueOrAskUser("package", cmd)
		if !ok {
			Ctx.ConfigManager.Println(":exclamation: Package name not provided")
			return
		}

		zipFilepath, ok := Ctx.ConfigManager.GetConfigValueOrAskUser("zip-file", cmd)
		if !ok {
			Ctx.ConfigManager.Println(":exclamation: Zip file not provided")
			return
		}

		var selectedPlatforms []int
		platformsInput, ok := Ctx.ConfigManager.GetFlagValue("platforms", cmd)
		if !ok {
			platforms, err := Ctx.PlatformClient.GetPlatforms(cmd.Context())
			if err != nil {
				Ctx.ConfigManager.Printf(":exclamation: %s\n", err.Error())
				return
			}

			platformOptions := make([]forms.Option, len(platforms))
			for i, platform := range platforms {
				platformOptions[i] = forms.Option{
					Value: fmt.Sprint(platform.ID),
					Label: platform.Name,
				}
			}

			selectedPlatforms, err = Ctx.FormHandler.MultiSelectIDs("Select PLATFORMS:\n", platformOptions)
			if err != nil || len(selectedPlatforms) == 0 {
				Ctx.ConfigManager.Println(":exclamation: Platforms not provided")
				return
			}
		} else {
			selectedPlatformStrings := strings.Split(platformsInput, ",")
			if len(selectedPlatformStrings) == 0 {
				Ctx.ConfigManager.Println(":exclamation: Platforms not provided")
				return
			}

			for _, selectedPlatformString := range selectedPlatformStrings {
				selectedPlatform, err := strconv.Atoi(selectedPlatformString)
				if err != nil {
					Ctx.ConfigManager.Println(":exclamation: Invalid platform ID")
					return
				}
				selectedPlatforms = append(selectedPlatforms, selectedPlatform)
			}
		}

		var selectedControlTypes []int
		controlTypesInput, ok := Ctx.ConfigManager.GetFlagValue("controls", cmd)
		if !ok {

			controlTypes, err := Ctx.PlatformClient.GetControlTypes(cmd.Context())
			if err != nil {
				Ctx.ConfigManager.Printf(":exclamation: %s\n", err.Error())
				return
			}

			controlTypeOptions := make([]forms.Option, len(controlTypes))
			for i, controlType := range controlTypes {
				controlTypeOptions[i] = forms.Option{
					Label: controlType.Name,
					Value: fmt.Sprint(controlType.ID),
				}
			}

			selectedControlTypes, err = Ctx.FormHandler.MultiSelectIDs("Select CONTROL TYPES:\n", controlTypeOptions)
			if err != nil || len(selectedControlTypes) == 0 {
				Ctx.ConfigManager.Println(":exclamation: Control types not provided")
				return
			}
		} else {
			selectedControlTypeStrings := strings.Split(controlTypesInput, ",")
			if len(selectedControlTypeStrings) == 0 {
				Ctx.ConfigManager.Println(":exclamation: Control types not provided")
				return
			}

			for _, selectedControlTypeString := range selectedControlTypeStrings {
				selectedControlType, err := strconv.Atoi(selectedControlTypeString)
				if err != nil {
					Ctx.ConfigManager.Println(":exclamation: Invalid control type ID")
					return
				}
				selectedControlTypes = append(selectedControlTypes, selectedControlType)
			}
		}

		input := graphql_api.ModuleVersion{
			ModuleID:        moduleID,
			LocalFilePath:   zipFilepath,
			SemanticVersion: semVer,
			Package:         packageName,
			PlatformIds:     selectedPlatforms,
			ControlIds:      selectedControlTypes,
		}

		//spinner := loader.NewSpinner(Ctx.ConfigManager)
		ctx, cancel := context.WithCancel(context.Background())
		err := spinner.New().
			Type(spinner.Line).
			Title("Deploying module version...").
			Context(ctx).
			Run()

		moduleVersion, err := Ctx.PlatformClient.CreateModuleVersion(cmd.Context(), input)
		cancel()
		//spinner.Stop()
		if err != nil {
			Ctx.ConfigManager.Printf(":exclamation: %s\n", err.Error())
			return
		}

		Ctx.ConfigManager.Printf("Deployed version %s for module %d\n", moduleVersion.SemanticVersion, moduleVersion.ModuleID)
	},
}

func init() {
	modulesCmd.AddCommand(modulesDeployCmd)

	modulesDeployCmd.Flags().StringP("semantic-version", "v", "", "Semantic version of the module version")
	modulesDeployCmd.Flags().StringP("package", "p", "", "Package name of the module version")
	modulesDeployCmd.Flags().StringP("zip-file", "f", "", "Zip file path of the module version")
	modulesDeployCmd.Flags().String("platforms", "", "Comma separated list of platform IDs")
	modulesDeployCmd.Flags().String("controls", "", "Comma separated list of control type IDs")
}
