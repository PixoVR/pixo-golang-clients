/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"errors"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/loader"
	"github.com/spf13/cobra"
	"os"
)

var (
	ManifestNotInitializedError = errors.New("asset manifest not initialized. Please run 'pixo assets init'")
)

// assetsInitCmd represents the assetsInit rootCmd
var assetsInitCmd = &cobra.Command{
	Use:           "init",
	Short:         "Initialize assets manifest for a module",
	Long:          `Initialize assets manifest for a module. If the module has existing assets, the existing manifest will be pulled.`,
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		questions := []config.Value{
			{Question: moduleQuestion()},
		}

		answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd)
		if err != nil {
			return err
		}

		moduleID := forms.Int(answers["module"])

		manifest, err := NewManifest()
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
			if _, err = os.Create(ManifestFilename); err != nil {
				return err
			}
			manifest = &Manifest{}
		}

		if manifest.ModuleID > 0 && manifest.ModuleID != moduleID {
			return errors.New("assets already initialized for another module")
		}

		manifest.ModuleID = moduleID

		loader.NewLoader(cmd.Context(), "Retrieving assets from the platform...", Ctx.Printer)
		assets, err := Ctx.PlatformClient.GetAssets(cmd.Context(), platform.AssetParams{ModuleID: moduleID})
		if err != nil {
			return err
		}

		for _, asset := range assets {
			manifest.AddAsset(asset)
		}

		return manifest.Save()
	},
}

func init() {
	assetsCmd.AddCommand(assetsInitCmd)
}
