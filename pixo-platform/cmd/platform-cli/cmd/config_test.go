package cmd_test

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/cmd"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/editor"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Config", func() {

	var (
		testConfigPath string
	)

	BeforeEach(func() {
		testConfigPath = fmt.Sprintf("%s/.pixo/test-config.yaml", os.Getenv("HOME"))
		if _, err := os.Stat(testConfigPath); err == nil {
			_ = os.Remove(testConfigPath)
		}
	})

	It("can set the lifecycle", func() {
		output, err := RunCommand("config", "set", "-l", "dev")
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("lifecycle : dev"))

		output, err = RunCommand("--config", testConfigPath, "config", "set", "-l", "stage")
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("lifecycle : stage"))
	})

	It("can set the region", func() {
		output, err := RunCommand("--config", testConfigPath, "config", "set", "-r", "saudi")
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("region : saudi"))

		output, err = RunCommand("--config", testConfigPath, "config", "set", "-r", "na")
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("region : na"))
	})

	It("can set the username and password", func() {
		output, err := RunCommand(
			"--config",
			testConfigPath,
			"config",
			"set",
			"--key",
			"username",
			"--val",
			"test",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("username : test"))

		output, err = RunCommand(
			"--config",
			testConfigPath,
			"config",
			"set",
			"--key",
			"password",
			"--val",
			"test",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("password : test"))
	})

	It("can open up the config file", func() {
		fileOpener := &editor.MockFileOpener{}

		cmd.FileOpener = fileOpener

		output, err := RunCommand("--config", testConfigPath, "config", "--edit")
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Opening config file"))
	})

})
