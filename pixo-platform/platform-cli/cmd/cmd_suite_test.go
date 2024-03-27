package cmd_test

import (
	"bytes"
	"fmt"
	graphql_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/cmd"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/clients"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/editor"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/forms/basic"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"io"
	"math/rand"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCLI(t *testing.T) {
	RegisterFailHandler(Fail)

	//root := GetProjectRoot()
	//envPath := filepath.Join(root, "../../../../.env")
	//
	//_ = godotenv.Load(envPath)

	RunSpecs(t, "Pixo Platform CLI Suite")
}

type TestExecutor struct {
	ConfigManager         config.Manager
	MockPlatformClient    *graphql_api.MockGraphQLClient
	MockMatchmakingClient *matchmaker.MockMatchmaker
	MockOldAPIClient      *graphql_api.MockGraphQLClient
	MockFileOpener        *editor.MockFileOpener
	configFile            string
}

func NewTestExecutor() *TestExecutor {
	randomID := rand.Intn(1000000000)
	testConfigPath := fmt.Sprintf("%s/.pixo/.test-config-%d.yaml", os.Getenv("HOME"), randomID)
	if _, err := os.Stat(testConfigPath); err == nil {
		if err = os.Remove(testConfigPath); err != nil {
			log.Warn().Msgf("unable to remove test config file: %s", err)
		} else {
			log.Info().Msgf("removed existing test config file: %s", testConfigPath)
		}
	} else {
		log.Info().Msgf("test config file does not exist: %s", testConfigPath)
	}

	formHandler := basic.NewFormHandler(nil, nil)
	configManager := config.NewFileManager("", formHandler)
	err := configManager.SetConfigFile(testConfigPath)
	Expect(err).NotTo(HaveOccurred())

	mockOldAPIClient := &graphql_api.MockGraphQLClient{}
	mockPlatformClient := &graphql_api.MockGraphQLClient{}
	mockMatchmaker := matchmaker.NewMockMatchmaker()
	mockFileOpener := &editor.MockFileOpener{}

	mockPlatformCtx := &clients.CLIContext{
		FormHandler:       formHandler,
		ConfigManager:     configManager,
		PlatformClient:    mockPlatformClient,
		MatchmakingClient: mockMatchmaker,
		OldAPIClient:      mockOldAPIClient,
		FileOpener:        mockFileOpener,
	}
	cmd.Ctx = mockPlatformCtx

	executor := &TestExecutor{
		ConfigManager:         configManager,
		MockPlatformClient:    mockPlatformClient,
		MockMatchmakingClient: mockMatchmaker,
		MockOldAPIClient:      mockOldAPIClient,
		MockFileOpener:        mockFileOpener,
		configFile:            testConfigPath,
	}

	return executor
}

func (t *TestExecutor) Cleanup() {
	log.Debug().Msgf("Cleaning up test config file: %s", t.configFile)
	viper.Reset()
	_ = os.Remove(t.configFile)
}

func (t *TestExecutor) RunCommandWithInput(reader io.Reader, args ...string) (string, error) {
	rootCmd := cmd.GetRootCmd()
	Expect(rootCmd).NotTo(BeNil())
	rootCmd.SetIn(reader)
	t.ConfigManager.SetReader(reader)

	writer := bytes.NewBufferString("")
	rootCmd.SetOut(writer)
	t.ConfigManager.SetWriter(writer)

	args = append([]string{"--config", t.configFile}, args...)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

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
	Expect(output).NotTo(BeEmpty())
	return output
}
