package cmd_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Describe("API Keys", func() {

	var (
		executor *TestExecutor
	)

	BeforeEach(func() {
		executor = NewTestExecutor()
	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("can create an api key", func() {
		output, err := executor.RunCommand("keys", "create")
		Expect(executor.MockPlatformClient.CalledCreateAPIKey).To(BeTrue())
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("API key created"))
		val, ok := executor.ConfigManager.GetConfigValue("api-key")
		Expect(val).NotTo(BeEmpty())
		Expect(ok).To(BeTrue())
	})

	//It("can create an api key for a user", func() {
	//	output, err := executor.RunCommand("keys", "create", "--user-id", "9999999")
	//	Expect(err).NotTo(HaveOccurred())
	//	Expect(output).To(ContainSubstring("API key created for user: 9999999"))
	//	Expect(executor.ConfigManager.APIKey()).NotTo(BeEmpty())
	//})

	It("can list api keys", func() {
		output := executor.RunCommandAndExpectSuccess("keys", "list")
		Expect(executor.MockPlatformClient.CalledGetAPIKeys).To(BeTrue())
		Expect(output).To(ContainSubstring("API keys:"))
		keys := strings.Split(output, "\n")
		Expect(len(keys)).To(BeNumerically(">", 1))
	})

	It("can delete an api key", func() {
		output := executor.RunCommandAndExpectSuccess("keys", "delete", "--key-id", "1")
		Expect(output).To(ContainSubstring("Deleted API key: 1"))
	})

	It("can list api keys for a user", func() {
		executor.MockPlatformClient.GetAPIKeysEmpty = true
		output, err := executor.RunCommand("keys", "list", "--user-id", "9999999")
		Expect(executor.MockPlatformClient.CalledGetAPIKeys).To(BeTrue())
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("No API keys found"))
	})

})
