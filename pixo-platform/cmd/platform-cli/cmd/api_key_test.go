package cmd_test

import (
	"fmt"
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
		Expect(err).NotTo(HaveOccurred())
		Expect(output).NotTo(BeEmpty())
		Expect(executor.ConfigManager.APIKey()).NotTo(BeEmpty())
	})

	//It("can create an api key for a user", func() {
	//	output, err := executor.RunCommand("keys", "create", "--user-id", "9999999")
	//	Expect(err).NotTo(HaveOccurred())
	//	Expect(output).To(ContainSubstring("API key created for user: 9999999"))
	//	Expect(executor.ConfigManager.APIKey()).NotTo(BeEmpty())
	//})

	It("can list and delete api keys", func() {
		output, err := executor.RunCommand("keys", "list")
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("API Keys:"))
		keyOne := strings.Split(output, "\n")[1]
		keyOneAPIKey := strings.TrimSpace(strings.Split(keyOne, ":")[1])

		output, err = executor.RunCommand("keys", "delete", "--key-id", keyOneAPIKey)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring(fmt.Sprintf("Deleted API key: %s", keyOneAPIKey)))
	})

	It("can list api keys for a user", func() {
		executor.MockPlatformClient.GetAPIKeysEmpty = true
		output, err := executor.RunCommand("keys", "list", "--user-id", "9999999")
		Expect(executor.MockPlatformClient.CalledGetAPIKeys).To(BeTrue())
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("No API Keys found"))
	})

})
