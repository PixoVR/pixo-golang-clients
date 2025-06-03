/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/spf13/cobra"
)

// assetsAddCmd represents the assetsAdd rootCmd
var assetsAddCmd = &cobra.Command{
	Use:           "add",
	Short:         "Add an abstract asset to a module",
	Long:          `Add an abstract asset to a module. Asset versions can be added to the asset after creation.`,
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		questions := []config.Value{
			{Question: forms.Question{Type: forms.Input, Key: "name"}},
			{Question: forms.Question{Type: forms.Input, Key: "type"}},
		}

		answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd)
		if err != nil {
			return err
		}

		name := forms.String(answers["name"])
		assetType := forms.String(answers["type"])

		manifest, err := NewManifest()
		if err != nil {
			return ManifestNotInitializedError
		}

		asset := platform.Asset{
			ModuleID: manifest.ModuleID,
			Name:     name,
			Type:     assetType,
		}
		if err := Ctx.PlatformClient.CreateAsset(cmd.Context(), &asset); err != nil {
			return err
		}

		manifest.AddAsset(asset)
		return manifest.Save()
	},
}

func init() {
	assetsCmd.AddCommand(assetsAddCmd)
	assetsAddCmd.Flags().StringP("name", "n", "", "Name of the abstract asset")
	assetsAddCmd.Flags().StringP("type", "t", "", "Type of the abstract asset")
}
