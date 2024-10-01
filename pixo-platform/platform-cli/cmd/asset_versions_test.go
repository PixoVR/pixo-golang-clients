package cmd_test

import (
	"bytes"
	"errors"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/cmd"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Asset Versions", func() {

	BeforeEach(func() {
		_ = os.Remove(cmd.ManifestFilename)
		executor = NewTestExecutor()
		executor.MockPlatformClient.GetAssetsReturns = []platform.Asset{}
		executor.RunCommandAndExpectSuccess("assets", "init", "--module", "TST")
		executor.RunCommandAndExpectSuccess(
			"assets",
			"add",
			"--name",
			"test",
			"--type",
			"image",
		)
	})

	AfterEach(func() {
		_ = os.Remove(cmd.ManifestFilename)
		executor.Cleanup()
	})

	It("returns an error if a valid status is not provided as the first arg", func() {
		_, err := executor.RunCommand("assets", "versions")
		Expect(err).To(MatchError("status (stage or release) is required as the first argument"))

		_, err = executor.RunCommand("assets", "versions", "nonexistent")
		Expect(err).To(MatchError("status (stage or release) is required as the first argument"))
	})

	It("returns an error if the asset name is not provided", func() {
		input := bytes.NewBufferString("")
		_, err := executor.RunCommandWithInput(
			input,
			"assets",
			"versions",
			"stage",
		)
		Expect(err).To(MatchError("ASSET not provided"))
	})

	Context("staging an asset version", func() {

		It("returns an error if the filepath is not provided", func() {
			executor.MockPlatformClient.GetAssetsReturns = []platform.Asset{
				{
					ID:       1,
					ModuleID: 1,
					Name:     "test",
				},
			}
			input := bytes.NewBufferString("")

			_, err := executor.RunCommandWithInput(
				input,
				"assets",
				"versions",
				"stage",
				"--asset",
				"test",
			)

			Expect(err).To(MatchError("FILEPATH not provided"))
		})

		It("returns an error if the manifest is not initialized", func() {
			_ = os.Remove(cmd.ManifestFilename)
			_, err := executor.RunCommand(
				"assets",
				"versions",
				"stage",
				"--asset",
				"test",
				"--filepath",
				"test.png",
			)
			Expect(err).To(MatchError("asset manifest not initialized. Please run 'pixo assets init'"))
		})

		It("returns an error if the asset is not found in the manifest", func() {
			_, err := executor.RunCommand(
				"assets",
				"versions",
				"stage",
				"--asset",
				"nonexistent",
				"--filepath",
				"test.png",
			)
			Expect(err).To(MatchError("asset not found in manifest. Please run 'pixo assets add'"))
		})

		It("can return an error if the assets cant be retrieved from the platform", func() {
			executor.MockPlatformClient.GetAssetsError = errors.New("get assets error")
			_, err := executor.RunCommand(
				"assets",
				"versions",
				"stage",
				"--asset",
				"test",
				"--filepath",
				"test.png",
			)
			Expect(err).To(MatchError("get assets error"))
		})

		It("can return an error if the asset is not found on the platform", func() {
			_, err := executor.RunCommand(
				"assets",
				"versions",
				"stage",
				"--asset",
				"test",
				"--filepath",
				"test.png",
			)
			Expect(err).To(MatchError("asset not found on pixo platform"))
		})

		Context("asset exists", func() {

			BeforeEach(func() {
				executor.MockPlatformClient.GetAssetsReturns = []platform.Asset{
					{
						ID:       1,
						ModuleID: 1,
						Name:     "test",
					},
				}
			})

			It("returns an error if unable to create the asset version", func() {
				executor.MockPlatformClient.CreateAssetVersionError = errors.New("create asset version error")
				_, err := executor.RunCommand(
					"assets",
					"versions",
					"stage",
					"--asset",
					"test",
					"--filepath",
					"test.png",
				)
				Expect(err).To(MatchError("create asset version error"))
			})

			It("can create an asset version in stage and add to the manifest", func() {
				executor.RunCommandAndExpectSuccess(
					"assets",
					"versions",
					"stage",
					"--asset",
					"test",
					"--lang",
					"es",
					"--filepath",
					"test.png",
				)

				Expect(executor.MockPlatformClient.CalledCreateAssetVersionWith).To(HaveLen(1))
				Expect(executor.MockPlatformClient.CalledCreateAssetVersionWith[0].AssetID).To(Equal(1))
				Expect(executor.MockPlatformClient.CalledCreateAssetVersionWith[0].Status).To(Equal("stage"))
				Expect(executor.MockPlatformClient.CalledCreateAssetVersionWith[0].LanguageCode).To(Equal("es"))
				Expect(executor.MockPlatformClient.CalledCreateAssetVersionWith[0].LocalFilePath).To(Equal("test.png"))

				manifest, err := cmd.NewManifest()
				Expect(err).NotTo(HaveOccurred())
				Expect(manifest).To(Equal(&cmd.Manifest{
					ModuleID:           1,
					ModuleAbbreviation: "",
					Assets: []cmd.Asset{
						{
							Name: "test",
							Type: "image",
							Versions: []cmd.Version{
								{
									ID:           1,
									Status:       "stage",
									LanguageCode: "es",
								},
							},
						},
					},
				}))
			})

		})

	})

	Context("promoting an asset version to release", func() {

		BeforeEach(func() {
			mockAssets := []platform.Asset{
				{
					ID:       1,
					ModuleID: 1,
					Name:     "test",
				},
				{
					ID:       2,
					ModuleID: 1,
					Name:     "test-2",
				},
			}
			executor.MockPlatformClient.GetAssetsReturns = mockAssets
			executor.RunCommandAndExpectSuccess(
				"assets",
				"versions",
				"stage",
				"--asset",
				"test",
				"--lang",
				"fr",
				"--filepath",
				"test.png",
			)
			executor.RunCommandAndExpectSuccess(
				"assets",
				"versions",
				"stage",
				"--asset",
				"test",
				"--lang",
				"es",
				"--filepath",
				"test.png",
			)
			executor.Cleanup()
			executor = NewTestExecutor()
			executor.MockPlatformClient.GetAssetsReturns = mockAssets
		})

		It("asks the user if they are sure with a warning", func() {
			input := bytes.NewBufferString("n\n")
			output, err := executor.RunCommandWithInput(
				input,
				"assets",
				"versions",
				"release",
				"--asset",
				"test",
			)
			Expect(err).To(MatchError("release cancelled"))
			Expect(output).To(ContainSubstring("Are you sure you want to release this asset version to production?"))
		})

		It("can promote asset versions to release without specifying language", func() {
			input := bytes.NewBufferString("\ny\n")

			executor.RunCommandWithInputAndExpectSuccess(
				input,
				"assets",
				"versions",
				"release",
				"--asset",
				"test",
			)

			Expect(executor.MockPlatformClient.CalledUpdateAssetVersionWith).To(HaveLen(2))
			Expect(executor.MockPlatformClient.CalledUpdateAssetVersionWith[0]).To(Equal(platform.AssetVersion{
				ID:            1,
				Status:        "release",
				LanguageCode:  "",
				LocalFilePath: "",
			}))
			Expect(executor.MockPlatformClient.CalledUpdateAssetVersionWith[1]).To(Equal(platform.AssetVersion{
				ID:            2,
				Status:        "release",
				LanguageCode:  "",
				LocalFilePath: "",
			}))

			manifest, err := cmd.NewManifest()
			Expect(err).NotTo(HaveOccurred())
			Expect(manifest).To(Equal(&cmd.Manifest{
				ModuleID:           1,
				ModuleAbbreviation: "",
				Assets: []cmd.Asset{
					{
						Name: "test",
						Type: "image",
						Versions: []cmd.Version{
							{
								ID:           1,
								Status:       "release",
								LanguageCode: "fr",
							},
							{
								ID:           2,
								Status:       "release",
								LanguageCode: "es",
							},
						},
					},
				},
			}))
		})

		It("can promote all asset versions for a given asset and language to release", func() {
			input := bytes.NewBufferString("y\n")

			executor.RunCommandWithInputAndExpectSuccess(
				input,
				"assets",
				"versions",
				"release",
				"--asset",
				"test",
				"--lang",
				"es",
			)

			Expect(executor.MockPlatformClient.CalledUpdateAssetVersionWith).To(HaveLen(1))
			Expect(executor.MockPlatformClient.CalledUpdateAssetVersionWith[0].ID).To(Equal(2))
			Expect(executor.MockPlatformClient.CalledUpdateAssetVersionWith[0].Status).To(Equal("release"))
			Expect(executor.MockPlatformClient.CalledUpdateAssetVersionWith[0].LanguageCode).To(Equal(""))
			Expect(executor.MockPlatformClient.CalledUpdateAssetVersionWith[0].LocalFilePath).To(Equal(""))

			manifest, err := cmd.NewManifest()
			Expect(err).NotTo(HaveOccurred())
			Expect(manifest).To(Equal(&cmd.Manifest{
				ModuleID:           1,
				ModuleAbbreviation: "",
				Assets: []cmd.Asset{
					{
						Name: "test",
						Type: "image",
						Versions: []cmd.Version{
							{
								ID:           1,
								Status:       "stage",
								LanguageCode: "fr",
							},
							{
								ID:           2,
								Status:       "release",
								LanguageCode: "es",
							},
						},
					},
				},
			}))
		})

	})

})
