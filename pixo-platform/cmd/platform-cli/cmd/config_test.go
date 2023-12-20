package cmd_test

import (
	"fmt"
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
		command, output := GetRootCmd()
		command.SetArgs([]string{
			"--config",
			testConfigPath,
			"config",
			"set",
			"-l",
			"dev",
		})
		err := command.Execute()
		Expect(err).NotTo(HaveOccurred())
		Expect(output.String()).To(ContainSubstring("lifecycle : dev"))

		command.SetArgs([]string{
			"--config",
			testConfigPath,
			"config",
			"set",
			"-l",
			"stage",
		})
		err = command.Execute()
		Expect(err).NotTo(HaveOccurred())
		Expect(output.String()).To(ContainSubstring("lifecycle : stage"))
	})

	It("can set the region", func() {
		command, output := GetRootCmd()
		command.SetArgs([]string{
			"--config",
			testConfigPath,
			"config",
			"set",
			"-r",
			"saudi",
		})
		err := command.Execute()
		Expect(err).NotTo(HaveOccurred())
		Expect(output.String()).To(ContainSubstring("region : saudi"))

		command.SetArgs([]string{
			"--config",
			testConfigPath,
			"config",
			"set",
			"-r",
			"na",
		})
		err = command.Execute()
		Expect(err).NotTo(HaveOccurred())
		Expect(output.String()).To(ContainSubstring("region : na"))
	})

	It("can set the username and password", func() {
		command, output := GetRootCmd()
		command.SetArgs([]string{
			"--config",
			testConfigPath,
			"config",
			"set",
			"--key",
			"username",
			"--val",
			"test",
		})
		err := command.Execute()
		Expect(err).NotTo(HaveOccurred())
		Expect(output.String()).To(ContainSubstring("username : test"))

		command.SetArgs([]string{
			"--config",
			testConfigPath,
			"config",
			"set",
			"--key",
			"password",
			"--val",
			"test",
		})
		err = command.Execute()
		Expect(err).NotTo(HaveOccurred())
		Expect(output.String()).To(ContainSubstring("password : test"))
	})

})
