package cmd_test

import (
	"bytes"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Describe("API Keys", func() {

	BeforeEach(func() {
		executor = NewTestExecutor()
	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("can return an error if the create call fails", func() {
		executor.MockPlatformClient.CreateAPIKeyError = errors.New("error")

		_, err := executor.RunCommand(
			"keys",
			"create",
		)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("error"))
		Expect(executor.MockPlatformClient.NumCalledCreateAPIKey).To(Equal(1))
	})

	It("can create an api key", func() {
		output, err := executor.RunCommand("keys", "create")

		Expect(executor.MockPlatformClient.NumCalledCreateAPIKey).To(Equal(1))
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

	It("can return an error if the get call fails", func() {
		executor.MockPlatformClient.GetAPIKeysError = errors.New("error")

		_, err := executor.RunCommand(
			"keys",
			"list",
		)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("error"))
		Expect(executor.MockPlatformClient.NumCalledGetAPIKeys).To(Equal(1))
	})

	It("can get api keys", func() {
		output := executor.RunCommandAndExpectSuccess("keys", "list")
		Expect(executor.MockPlatformClient.NumCalledGetAPIKeys).To(Equal(1))
		Expect(output).To(ContainSubstring("API keys:"))
		keys := strings.Split(output, "\n")
		Expect(len(keys)).To(BeNumerically(">", 1))
	})

	It("can list api keys for a user", func() {
		executor.MockPlatformClient.GetAPIKeysEmpty = true
		output, err := executor.RunCommand("keys", "list", "--user-id", "9999999")
		Expect(executor.MockPlatformClient.NumCalledGetAPIKeys).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("No API keys found"))
	})

	It("can return an error if there is no key provided when deleting", func() {
		input := bytes.NewBufferString("\n")
		_, err := executor.RunCommandWithInput(input, "keys", "delete")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("KEY IDs not provided"))
	})

	It("can return an error if the delete call fails", func() {
		executor.MockPlatformClient.DeleteAPIKeyError = errors.New("error")

		output, err := executor.RunCommand(
			"keys",
			"delete",
			"--key-ids",
			"Key ID 1",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Error deleting API key 1: error"))
		Expect(executor.MockPlatformClient.NumCalledDeleteAPIKey).To(Equal(1))
	})

	It("can delete an api key", func() {
		output := executor.RunCommandAndExpectSuccess("keys", "delete", "--key-ids", "Key ID 1")
		Expect(output).To(ContainSubstring("Deleted API key: 1"))
		Expect(executor.MockPlatformClient.NumCalledDeleteAPIKey).To(Equal(1))
	})

	It("can delete several api keys", func() {
		output := executor.RunCommandAndExpectSuccess("keys", "delete", "--key-ids", "Key ID 1,Key ID 2,Key ID 3")
		Expect(output).To(ContainSubstring("Deleted API key: 1"))
		Expect(output).To(ContainSubstring("Deleted API key: 2"))
		Expect(output).To(ContainSubstring("Deleted API key: 3"))
		Expect(executor.MockPlatformClient.NumCalledDeleteAPIKey).To(Equal(3))
	})

})
