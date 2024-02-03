package cmd_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConfigFile", func() {

	var (
		executor *TestExecutor
	)

	BeforeEach(func() {
		executor = NewTestExecutor()
	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("can set the lifecycle", func() {
		output, err := executor.RunCommand("config", "set", "-l", "dev")
		Expect(err).NotTo(HaveOccurred())
		Expect(output).NotTo(BeEmpty())
		Expect(executor.ConfigManager.Lifecycle()).To(Equal("dev"))
		Expect(executor.ConfigManager.Region()).To(Equal("na"))

		output, err = executor.RunCommand("config", "set", "-l", "stage")
		Expect(err).NotTo(HaveOccurred())
		Expect(output).NotTo(BeEmpty())
		Expect(executor.ConfigManager.Lifecycle()).To(Equal("stage"))
		Expect(executor.ConfigManager.Region()).To(Equal("na"))
	})

	It("can set the region", func() {
		output, err := executor.RunCommand("config", "set", "-r", "saudi", "-l", "prod")
		Expect(err).NotTo(HaveOccurred())
		Expect(output).NotTo(BeEmpty())
		Expect(executor.ConfigManager.Lifecycle()).To(Equal("prod"))
		Expect(executor.ConfigManager.Region()).To(Equal("saudi"))
		output, err = executor.RunCommand("config")
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Region: saudi"))
		Expect(output).To(ContainSubstring("Lifecycle: prod"))

		output, err = executor.RunCommand("config", "set", "-r", "na", "-l", "dev")
		Expect(err).NotTo(HaveOccurred())
		Expect(output).NotTo(BeEmpty())
		Expect(executor.ConfigManager.Lifecycle()).To(Equal("dev"))
		Expect(executor.ConfigManager.Region()).To(Equal("na"))
		Expect(output).To(ContainSubstring("Region: na"))
		Expect(output).To(ContainSubstring("Lifecycle: dev"))
	})

	It("can set the username and password", func() {
		output, err := executor.RunCommand(
			"config",
			"set",
			"--username",
			"test",
			"--password",
			"testpassword",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).NotTo(BeEmpty())
		val, ok := executor.ConfigManager.GetConfigValue("username")
		Expect(val).To(Equal("test"))
		Expect(ok).To(BeTrue())
		val, ok = executor.ConfigManager.GetConfigValue("password")
		Expect(val).To(Equal("testpassword"))
		Expect(ok).To(BeTrue())
	})

	It("can open up the config file", func() {
		output, err := executor.RunCommand("config", "--edit")
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Opening config file"))
		Expect(executor.MockFileOpener.CalledOpenEditor).To(BeTrue())
	})

})
