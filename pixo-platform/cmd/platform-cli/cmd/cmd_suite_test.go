package cmd_test

import (
	"bytes"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/cmd"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func GetProjectRoot() (root string) {
	_, b, _, _ := runtime.Caller(0)
	root = filepath.Dir(b)
	return root
}

func TestCLI(t *testing.T) {
	RegisterFailHandler(Fail)

	root := GetProjectRoot()
	envPath := filepath.Join(root, "../../../../.env")

	if err := godotenv.Load(envPath); err != nil {
		log.Warn().Msgf("Failed to load .env file at %s", envPath)
	}

	RunSpecs(t, "Pixo Platform CLI Suite")
}

func GetRootCmd() (*cobra.Command, *bytes.Buffer) {
	rootCmd := cmd.NewRootCmd()
	Expect(rootCmd).NotTo(BeNil())

	output := bytes.NewBufferString("")
	rootCmd.SetOut(output)
	return rootCmd, output
}

var _ = BeforeSuite(func() {
	assertLogin()
})

func assertLogin() {
	username := os.Getenv("PIXO_USERNAME")
	password := os.Getenv("PIXO_PASSWORD")

	rootCmd, output := GetRootCmd()
	rootCmd.SetArgs([]string{
		"auth",
		"login",
		"--username",
		username,
		"--password",
		password,
	})
	err := rootCmd.Execute()
	Expect(err).NotTo(HaveOccurred())

	out, err := io.ReadAll(output)
	Expect(err).NotTo(HaveOccurred())
	Expect(string(out)).To(ContainSubstring("Login successful. Here is your API token:"))
}
