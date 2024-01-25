package cmd_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Describe("API Keys", func() {

	It("can create an api key", func() {
		output, err := RunCommand("apiKeys", "create")
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("API key created"))

		output, err = RunCommand("config", "list")
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("api-key : "))
	})

	It("can list and delete api keys", func() {
		output, err := RunCommand("apiKeys", "list")
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("API Keys:"))
		keyOne := strings.Split(output, "\n")[1]
		keyOneAPIKey := strings.TrimSpace(strings.Split(keyOne, ":")[1])

		output, err = RunCommand("apiKeys", "delete", "--key-id", keyOneAPIKey)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring(fmt.Sprintf("Deleted API key: %s", keyOneAPIKey)))
	})

	It("can list api keys for a user", func() {
		output, err := RunCommand("apiKeys", "list", "--user-id", "9999999")
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("No API Keys found"))
	})

})
