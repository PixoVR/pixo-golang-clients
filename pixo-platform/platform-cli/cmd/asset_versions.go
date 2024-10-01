/*
Copyright © 2024 Walker O'Brien walker.obrien@pixovr.com
*/
package cmd

import (
	"errors"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/spf13/cobra"
)

// assetsVersionCmd represents the assetsVersion rootCmd
var assetsVersionCmd = &cobra.Command{
	Use:           "versions",
	Short:         "Manage physical asset versions",
	Long:          `Manage physical asset version for an abstract asset. The asset must be created before versions can be managed.`,
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 || (args[0] != "stage" && args[0] != "release") {
			return errors.New("status (stage or release) is required as the first argument")
		}
		status := args[0]

		questions := []config.Value{
			{Question: forms.Question{Type: forms.Input, Key: "asset"}},
			{Question: forms.Question{Type: forms.Input, Key: "lang", Optional: true}},
		}

		answers, err := Ctx.ConfigManager.GetValuesOrSubmitForm(questions, cmd)
		if err != nil {
			return err
		}

		assetName := forms.String(answers["asset"])
		languageCode := forms.String(answers["lang"])

		manifest, err := NewManifest()
		if err != nil {
			return ManifestNotInitializedError
		}

		if manifest.GetAsset(assetName) == nil {
			return errors.New("asset not found in manifest. Please run 'pixo assets add'")
		}

		assets, err := Ctx.PlatformClient.GetAssets(cmd.Context(), platform.AssetParams{
			ModuleID:     manifest.ModuleID,
			Name:         assetName,
			LanguageCode: languageCode,
		})
		if err != nil {
			return err
		} else if len(assets) == 0 {
			return errors.New("asset not found on pixo platform")
		}

		var assetVersion platform.AssetVersion

		if status == "release" {
			confirmQuestion := forms.Question{
				Type:   forms.Confirm,
				Key:    "confirmed",
				Prompt: ":warning: Are you sure you want to release this asset version to production? This will cause the version to be used by modules that use the release status to pull versions. (y/n)",
			}
			answers, err = Ctx.FormHandler.AskQuestions(confirmQuestion)
			if err != nil {
				return err
			}

			if !forms.Bool(answers["confirmed"]) {
				return errors.New("release cancelled")
			}

			manifestAsset := manifest.GetAsset(assetName)
			for _, version := range manifestAsset.Versions {
				if version.Status != "stage" || (languageCode != "" && version.LanguageCode != languageCode) {
					continue
				}

				assetVersion = platform.AssetVersion{
					ID:     version.ID,
					Status: "release",
				}
				if err := Ctx.PlatformClient.UpdateAssetVersion(cmd.Context(), &assetVersion); err != nil {
					return err
				}
				manifest.SetVersionToStatus(assetName, version.ID, languageCode, "stage", "release")
			}

		} else if status == "stage" {
			filepathForm := []config.Value{
				{Question: forms.Question{Type: forms.Input, Key: "filepath"}},
			}
			answers, err = Ctx.ConfigManager.GetValuesOrSubmitForm(filepathForm, cmd)
			if err != nil {
				return err
			}

			// remove this once language is not enforced
			if languageCode == "en" {
				status = "release"
			}

			assetVersion = platform.AssetVersion{
				AssetID:       assets[0].ID,
				LocalFilePath: forms.String(answers["filepath"]),
				Status:        status,
				LanguageCode:  languageCode,
			}
			if err := Ctx.PlatformClient.CreateAssetVersion(cmd.Context(), &assetVersion); err != nil {
				return err
			}
			manifest.AddVersion(assetName, assetVersion)
		}

		return manifest.Save()
	},
}

func init() {
	assetsCmd.AddCommand(assetsVersionCmd)
	assetsVersionCmd.Flags().StringP("asset", "a", "", "Name of the abstract asset to manage versions for")
	assetsVersionCmd.Flags().StringP("filepath", "f", "", "Path to the asset version file")
	assetsVersionCmd.Flags().StringP("lang", "l", "", "Language code for the asset version")
}
