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

var _ = Describe("Assets Init", func() {

	BeforeEach(func() {
		_ = os.Remove(cmd.ManifestFilename)
		executor = NewTestExecutor()
	})

	AfterEach(func() {
		_ = os.Remove(cmd.ManifestFilename)
		executor.Cleanup()
	})

	It("returns an error if the module is not provided", func() {
		input := bytes.NewBufferString("")
		_, err := executor.RunCommandWithInput(
			input,
			"assets",
			"init",
		)
		Expect(err).To(MatchError("MODULE not provided"))
	})

	It("creates the manifest file", func() {
		executor.MockPlatformClient.GetAssetsReturns = []platform.Asset{}

		_, err := executor.RunCommand(
			"assets",
			"init",
			"--module",
			"TST",
		)

		Expect(err).NotTo(HaveOccurred())
		contents, err := os.ReadFile(cmd.ManifestFilename)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(contents)).To(Equal("moduleId: 1\nassets: []\n"))
	})

	It("won't overwrite an existing manifest file", func() {
		Expect(os.WriteFile(cmd.ManifestFilename, []byte("moduleId: 2\n"), 0644)).To(Succeed())

		_, err := executor.RunCommand(
			"assets",
			"init",
			"--module",
			"TST",
		)

		Expect(err).To(MatchError("assets already initialized for another module"))
		contents, err := os.ReadFile(cmd.ManifestFilename)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(contents)).To(Equal("moduleId: 2\n"))
	})

	It("returns an error if unable to get the assets for a module", func() {
		executor.MockPlatformClient.GetAssetsError = errors.New("get assets error")
		_, err := executor.RunCommand(
			"assets",
			"init",
			"--module",
			"TST",
		)
		Expect(err).To(MatchError("get assets error"))
	})

	It("can pull any existing module assets", func() {
		executor.MockPlatformClient.GetAssetsReturns = []platform.Asset{
			{
				ID:       1,
				ModuleID: 1,
				Name:     "test",
				Type:     "image",
				Versions: []platform.AssetVersion{
					{
						ID:           1,
						Status:       "stage",
						LanguageCode: "en",
					},
				},
			},
		}

		_, err := executor.RunCommand(
			"assets",
			"init",
			"--module",
			"TST",
		)

		Expect(err).NotTo(HaveOccurred())
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
							LanguageCode: "en",
						},
					},
				},
			},
		}))
	})

})
