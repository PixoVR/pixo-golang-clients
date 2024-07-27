package cmd_test

import (
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Webhooks List", func() {

	BeforeEach(func() {
		executor = NewTestExecutor()
	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("can return an error if the get webhooks call fails", func() {
		executor.MockPlatformClient.GetWebhooksError = errors.New("get webhooks error")

		output, err := executor.RunCommand(
			"webhooks",
			"list",
		)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("get webhooks error"))
		Expect(output).To(ContainSubstring("Failed to get webhooks"))
		Expect(executor.MockPlatformClient.NumCalledGetWebhooks).To(Equal(1))
	})

	It("can list webhooks", func() {
		output, err := executor.RunCommand(
			"webhooks",
			"list",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("1. Description: test"))
		Expect(output).To(ContainSubstring("URL: https://example.com"))
		Expect(output).To(ContainSubstring("2. Description: test-2"))
		Expect(output).To(ContainSubstring("URL: https://example-2.com"))
	})

})
