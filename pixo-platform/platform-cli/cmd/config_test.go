package cmd_test

import (
	"fmt"
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", Ordered, func() {

	BeforeEach(func() {
		executor = NewTestExecutor()
	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("can show the platform urls", func() {
		output := executor.RunCommandAndExpectSuccess("config", "urls")
		Expect(output).To(ContainSubstring("Region: na"))
		Expect(output).To(ContainSubstring("Status: prod"))
		Expect(output).To(ContainSubstring("Web: https://apex.pixovr.com"))
		Expect(output).To(ContainSubstring("Platform API: https://apex.pixovr.com/v2"))
		Expect(output).To(ContainSubstring("Platform API Docs: https://apex.pixovr.com/v2/swagger/index.html"))
		Expect(output).To(ContainSubstring("Heartbeat API: https://apex.pixovr.com/heartbeat"))
		Expect(output).To(ContainSubstring("Matchmaking API Docs: https://apex.pixovr.com/matchmaking/swagger/index.html"))
		Expect(output).To(ContainSubstring("Matchmaking API: https://apex.pixovr.com/matchmaking"))
		Expect(output).To(ContainSubstring("Heartbeat API Docs: https://apex.pixovr.com/heartbeat/swagger/index.html"))

		_ = executor.RunCommandAndExpectSuccess("config", "set", "-l", "dev")
		output = executor.RunCommandAndExpectSuccess("config", "urls")
		Expect(output).To(ContainSubstring("Web: https://apex.dev.pixovr.com"))
		Expect(output).To(ContainSubstring("Platform API: https://apex.dev.pixovr.com/v2"))
		Expect(output).To(ContainSubstring("Platform API Docs: https://apex.dev.pixovr.com/v2/swagger/index.html"))
		Expect(output).To(ContainSubstring("Matchmaking API: https://apex.dev.pixovr.com/matchmaking"))
		Expect(output).To(ContainSubstring("Matchmaking API Docs: https://apex.dev.pixovr.com/matchmaking/swagger/index.html"))
		Expect(output).To(ContainSubstring("Heartbeat API: https://apex.dev.pixovr.com/heartbeat"))
		Expect(output).To(ContainSubstring("Heartbeat API Docs: https://apex.dev.pixovr.com/heartbeat/swagger/index.html"))
	})

	It("can show the current config", func() {
		username := faker.Username()
		password := faker.Password()
		executor.ExpectLoginToSucceed(username, password)

		token, ok := executor.ConfigManager.GetConfigValue("auth-token")
		Expect(ok).To(BeTrue())
		Expect(token).NotTo(BeEmpty())

		userID, ok := executor.ConfigManager.GetIntConfigValue("auth-user-id")
		Expect(ok).To(BeTrue())
		Expect(userID).NotTo(BeZero())

		_ = executor.RunCommandAndExpectSuccess("config", "set", "-k", "test", "-v", "testvalue")
		_ = executor.RunCommandAndExpectSuccess("config", "set", "-k", "api-key", "-v", "testapikey")

		output := executor.RunCommandAndExpectSuccess("config")

		Expect(output).NotTo(ContainSubstring(password))
		Expect(output).NotTo(ContainSubstring(token))
		Expect(output).To(ContainSubstring(fmt.Sprintf("User ID: %d", userID)))
		Expect(output).To(ContainSubstring("Region: na"))
		Expect(output).To(ContainSubstring("Status: prod"))
		Expect(output).To(ContainSubstring("Username: " + username))
		Expect(output).To(ContainSubstring("API Key: ********"))
		Expect(output).NotTo(ContainSubstring("testapikey"))
		Expect(output).To(ContainSubstring("Password: ********"))
		Expect(output).NotTo(ContainSubstring(password))
		Expect(output).To(ContainSubstring("Token: ********"))
		Expect(output).To(ContainSubstring("Test: testvalue"))
	})

	It("can set the lifecycle", func() {
		_ = executor.RunCommandAndExpectSuccess("config", "set", "-l", "dev")
		Expect(executor.ConfigManager.Lifecycle()).To(Equal("dev"))
		Expect(executor.ConfigManager.Region()).To(Equal("na"))

		_ = executor.RunCommandAndExpectSuccess("config", "set", "-l", "stage")
		Expect(executor.ConfigManager.Lifecycle()).To(Equal("stage"))
		Expect(executor.ConfigManager.Region()).To(Equal("na"))
	})

	It("can set the region", func() {
		output := executor.RunCommandAndExpectSuccess("config", "set", "-r", "saudi", "-l", "prod")
		Expect(output).To(ContainSubstring("Region: saudi"))
		Expect(output).To(ContainSubstring("Status: prod"))
		Expect(executor.ConfigManager.Lifecycle()).To(Equal("prod"))
		Expect(executor.ConfigManager.Region()).To(Equal("saudi"))
		output = executor.RunCommandAndExpectSuccess("config")
		Expect(output).To(ContainSubstring("Region: saudi"))
		Expect(output).To(ContainSubstring("Status: prod"))

		output = executor.RunCommandAndExpectSuccess("config", "set", "-r", "na", "-l", "dev")

		Expect(executor.ConfigManager.Lifecycle()).To(Equal("dev"))
		Expect(executor.ConfigManager.Region()).To(Equal("na"))
		Expect(output).To(ContainSubstring("Region: na"))
		Expect(output).To(ContainSubstring("Status: dev"))
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
		output := executor.RunCommandAndExpectSuccess("config", "--edit")
		Expect(output).To(ContainSubstring("Opening config file"))
		Expect(executor.MockFileOpener.CalledOpenEditor).To(BeTrue())
	})

})
