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
	"github.com/spf13/cobra"
)

// modulesDeployCmd represents the modules deploy rootCmd
var modulesDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy module versions",
	Long: `Deploy module versions
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		questions := []config.Value{
			{Question: forms.Question{
				Type: forms.SelectID,
				Key:  "module-id",
				LabelFunc: func(i interface{}) string {
					item := i.(platform.Module)
					return fmt.Sprintf("%d: %s - %s", item.ID, item.Abbreviation, item.Name)
				},
				GetItemsFunc: func(ctx context.Context) (interface{}, error) {
					items, err := Ctx.PlatformClient.GetModules(cmd.Context())
					if err != nil {
						Ctx.Printer.Println(":exclamation: Unable to get modules")
						return nil, errors.New("unable to get modules")
					}

					return items, nil
				},
			}},
			{Question: forms.Question{Type: forms.Input, Key: "semantic-version"}},
			{Question: forms.Question{Type: forms.Input, Key: "package"}},
			{Question: forms.Question{Type: forms.Input, Key: "zip-file"}},
			{Question: forms.Question{
				Type: forms.MultiSelectIDs,
				Key:  "platforms",
				GetItemsFunc: func(ctx context.Context) (interface{}, error) {
					items, err := Ctx.PlatformClient.GetPlatforms(cmd.Context())
					if err != nil {
						return nil, errors.New("unable to get platforms")
					}
					return items, nil
				},
			}},
			{Question: forms.Question{
				Type: forms.MultiSelectIDs,
				Key:  "controls",
				GetItemsFunc: func(ctx context.Context) (interface{}, error) {
					items, err := Ctx.PlatformClient.GetControlTypes(cmd.Context())
					if err != nil {
						return nil, errors.New("unable to get control types")
					}

					return items, nil
				},
			}},
		}

		answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd)
		if err != nil {
			return err
		}

		moduleID := forms.Int(answers["module-id"])
		semVer := forms.String(answers["semantic-version"])
		packageName := forms.String(answers["package"])
		zipfilePath := forms.String(answers["zip-file"])
		platforms := forms.IntSlice(answers["platforms"])
		controlTypes := forms.IntSlice(answers["control-types"])

		input := platform.ModuleVersion{
			ModuleID:        moduleID,
			LocalFilePath:   zipfilePath,
			SemanticVersion: semVer,
			Package:         packageName,
			PlatformIds:     platforms,
			ControlIds:      controlTypes,
		}

		spinner := loader.NewLoader(cmd.Context(), "Deploying module version...", Ctx.Printer)
		defer spinner.Stop()

		moduleVersion, err := Ctx.PlatformClient.CreateModuleVersion(cmd.Context(), input)
		if err != nil {
			return err
		}

		Ctx.Printer.Printf("Deployed version %s for module %d\n", moduleVersion.SemanticVersion, moduleVersion.ModuleID)
		return nil
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
