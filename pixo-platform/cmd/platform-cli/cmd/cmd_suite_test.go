package cmd_test

import (
	"bytes"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/cmd"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/clients"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/editor"
	graphql_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
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

	configManager := config.NewFileManager("")
	err := configManager.ReadConfigFile(testConfigPath)
	Expect(err).NotTo(HaveOccurred())

	mockOldAPIClient := &graphql_api.MockGraphQLClient{}
	mockPlatformClient := &graphql_api.MockGraphQLClient{}
	mockMatchmaker := matchmaker.NewMockMatchmaker()
	mockFileOpener := &editor.MockFileOpener{}

	mockPlatformCtx := &clients.PlatformContext{
		ConfigManager:     configManager,
		PlatformClient:    mockPlatformClient,
		MatchmakingClient: mockMatchmaker,
		OldAPIClient:      mockOldAPIClient,
		FileOpener:        mockFileOpener,
	}
	cmd.PlatformCtx = mockPlatformCtx

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

func (t *TestExecutor) RunCommandWithReadWriter(reader io.Reader, writer io.Writer, args ...string) (string, error) {
	mockPlatformCtx := &clients.PlatformContext{
		ConfigManager:     t.ConfigManager,
		PlatformClient:    t.MockPlatformClient,
		MatchmakingClient: t.MockMatchmakingClient,
		OldAPIClient:      t.MockOldAPIClient,
		FileOpener:        t.MockFileOpener,
	}
	rootCmd := cmd.NewRootCmd(mockPlatformCtx)
	Expect(rootCmd).NotTo(BeNil())
	rootCmd.SetOut(writer)
	rootCmd.SetIn(reader)

	output := bytes.NewBufferString("")
	rootCmd.SetOut(output)

	args = append([]string{"--config", t.configFile}, args...)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	out, _ := io.ReadAll(output)
	return string(out), err
}
func (t *TestExecutor) RunCommand(args ...string) (string, error) {
	return t.RunCommandWithReadWriter(os.Stdin, os.Stdout, args...)
}

func (t *TestExecutor) ExpectLoginToSucceed(username, password string) {

	output, err := t.RunCommand(
		"auth",
		"login",
		"--username",
		username,
		"--password",
		password,
	)

	Expect(err).NotTo(HaveOccurred())
	Expect(output).To(ContainSubstring("Login successful. Here is your API token:"))
	userID, ok := t.ConfigManager.GetConfigValue("user-id")
	Expect(ok).To(BeTrue())
	Expect(userID).To(Equal(fmt.Sprint(t.MockPlatformClient.ActiveUserID())))
}
