package cmd_test

import (
	"bytes"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/headset"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/cmd"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/ctx"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/editor"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms/basic"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	executor *TestExecutor
)

func TestCLI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pixo Platform CLI Suite")
}

type TestExecutor struct {
	rootCmd               *cobra.Command
	Printer               *printer.EmojiPrinter
	FormHandler           *basic.FormHandler
	ConfigManager         *config.ConfigManager
	InMemoryConfig        *config.InMemoryConfigManager
	MockPlatformClient    *platform.MockClient
	MockHeadsetClient     *headset.MockClient
	MockMatchmakingClient *matchmaker.MockMatchmaker
	MockFileOpener        *editor.MockFileOpener
}

func NewTestExecutor() *TestExecutor {
	formHandler := basic.NewFormHandler(nil, nil)
	inMemoryConfigManager := config.NewInMemoryConfigManager()
	configManager := config.NewConfigManager(inMemoryConfigManager)

	mockPlatformClient := &platform.MockClient{}
	mockHeadsetClient := &headset.MockClient{}
	mockMatchmaker := matchmaker.NewMockMatchmaker()

	emojiPrinter := printer.NewEmojiPrinter(nil)
	mockFileOpener := &editor.MockFileOpener{}

	cmd.Ctx = &ctx.CLIContext{
		Printer:           emojiPrinter,
		FormHandler:       formHandler,
		ConfigManager:     configManager,
		PlatformClient:    mockPlatformClient,
		HeadsetClient:     mockHeadsetClient,
		MatchmakingClient: mockMatchmaker,
		FileOpener:        mockFileOpener,
	}

	rootCmd := cmd.GetRootCmd()
	return &TestExecutor{
		rootCmd:               rootCmd,
		FormHandler:           formHandler,
		Printer:               emojiPrinter,
		InMemoryConfig:        inMemoryConfigManager,
		ConfigManager:         configManager,
		MockPlatformClient:    mockPlatformClient,
		MockHeadsetClient:     mockHeadsetClient,
		MockMatchmakingClient: mockMatchmaker,
		MockFileOpener:        mockFileOpener,
	}
}

func (t *TestExecutor) Cleanup() {
	t.InMemoryConfig.Clear()
	t.MockPlatformClient.Reset()
	t.rootCmd.SetArgs(nil)
	t.clearFlagValues(t.rootCmd)

}

func (t *TestExecutor) clearFlagValues(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		_ = flag.Value.Set(flag.DefValue)
	})
	for _, subCmd := range cmd.Commands() {
		t.clearFlagValues(subCmd)
	}
}

func (t *TestExecutor) RunCommandWithInput(reader io.Reader, args ...string) (string, error) {
	t.rootCmd.SetIn(reader)
	t.FormHandler.SetReader(reader)
	t.ConfigManager.SetReader(reader)

	writer := bytes.NewBufferString("")
	t.rootCmd.SetOut(writer)
	t.FormHandler.SetWriter(writer)
	t.ConfigManager.SetWriter(writer)
	t.Printer.SetWriter(writer)

	t.rootCmd.SetArgs(args)
	err := t.rootCmd.Execute()

	output, _ := io.ReadAll(writer)
	return string(output), err
}

func (t *TestExecutor) RunCommandWithInputAndExpectSuccess(input io.Reader, args ...string) string {
	output, err := t.RunCommandWithInput(input, args...)
	Expect(err).NotTo(HaveOccurred())
	Expect(output).NotTo(BeEmpty())
	return output
}

func (t *TestExecutor) RunCommand(args ...string) (string, error) {
	return t.RunCommandWithInput(os.Stdin, args...)
}

func (t *TestExecutor) RunCommandAndExpectSuccess(args ...string) string {
	output, err := t.RunCommandWithInput(os.Stdin, args...)
	Expect(err).NotTo(HaveOccurred())
	return output
}
