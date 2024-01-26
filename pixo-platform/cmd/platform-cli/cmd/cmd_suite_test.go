package cmd_test

import (
	"bytes"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/cmd"
	"github.com/joho/godotenv"
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

	_ = godotenv.Load(envPath)

	RunSpecs(t, "Pixo Platform CLI Suite")
}

func GetRootCmd() (*cobra.Command, *bytes.Buffer) {
	rootCmd := cmd.NewRootCmd()
	Expect(rootCmd).NotTo(BeNil())

	output := bytes.NewBufferString("")
	rootCmd.SetOut(output)
	return rootCmd, output
}

func RunCommand(args ...string) (string, error) {
	rootCmd, output := GetRootCmd()
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	out, _ := io.ReadAll(output)
	return string(out), err
}

var _ = BeforeSuite(func() {
	assertLogin()
})

func assertLogin() {
	username := os.Getenv("PIXO_USERNAME")
	password := os.Getenv("PIXO_PASSWORD")

	output, err := RunCommand("auth", "login", "--username", username, "--password", password)
	Expect(err).NotTo(HaveOccurred())
	Expect(output).To(ContainSubstring("Login successful. Here is your API token:"))

	output, err = RunCommand("config", "list")
	Expect(err).NotTo(HaveOccurred())
	Expect(output).To(ContainSubstring("user-id : "))
}
