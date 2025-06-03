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

var _ = Describe("Assets Add", func() {

	BeforeEach(func() {
		_ = os.Remove(cmd.ManifestFilename)
		executor = NewTestExecutor()
		executor.MockPlatformClient.GetAssetsReturns = []platform.Asset{}
		executor.RunCommandAndExpectSuccess("assets", "init", "--module", "TST")
	})

	AfterEach(func() {
		_ = os.Remove(cmd.ManifestFilename)
		executor.Cleanup()
	})

	It("returns an error if the name is not provided", func() {
		input := bytes.NewBufferString("")
		_, err := executor.RunCommandWithInput(
			input,
			"assets",
			"add",
		)
		Expect(err).To(MatchError("NAME not provided"))
	})

	It("returns an error if the type is not provided", func() {
		input := bytes.NewBufferString("test\n")
		_, err := executor.RunCommandWithInput(
			input,
			"assets",
			"add",
		)
		Expect(err).To(MatchError("TYPE not provided"))
	})

	It("returns an error if the manifest is not initialized", func() {
		_ = os.Remove(cmd.ManifestFilename)
		_, err := executor.RunCommand(
			"assets",
			"add",
			"--name",
			"test",
			"--type",
			"image",
		)
		Expect(err).To(MatchError("asset manifest not initialized. Please run 'pixo assets init'"))
	})

	It("returns an error if unable to create the asset", func() {
		executor.MockPlatformClient.CreateAssetError = errors.New("create asset error")
		_, err := executor.RunCommand(
			"assets",
			"add",
			"--name",
			"test",
			"--type",
			"image",
		)
		Expect(err).To(MatchError("create asset error"))
	})

	It("can create a module asset and add to the manifest", func() {
		_, err := executor.RunCommand(
			"assets",
			"add",
			"--name",
			"test",
			"--type",
			"image",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(executor.MockPlatformClient.CalledCreateAssetWith).To(HaveLen(1))
		Expect(executor.MockPlatformClient.CalledCreateAssetWith[0].ModuleID).To(Equal(1))
		Expect(executor.MockPlatformClient.CalledCreateAssetWith[0].Name).To(Equal("test"))
		Expect(executor.MockPlatformClient.CalledCreateAssetWith[0].Type).To(Equal("image"))
		manifest, err := cmd.NewManifest()
		Expect(err).NotTo(HaveOccurred())
		Expect(manifest).To(Equal(&cmd.Manifest{
			ModuleID:           1,
			ModuleAbbreviation: "",
			Assets: []cmd.Asset{
				{
					Name:     "test",
					Type:     "image",
					Versions: nil,
				},
			},
		}))
	})

})
