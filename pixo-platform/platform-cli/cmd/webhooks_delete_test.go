package cmd_test

import (
	"bytes"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Webhooks Delete", func() {

	BeforeEach(func() {
		executor = NewTestExecutor()
	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("can return an error if the id is not provided", func() {
		input := bytes.NewBufferString("\n")

		_, err := executor.RunCommandWithInput(
			input,
			"webhooks",
			"delete",
		)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("WEBHOOK IDs not provided"))
		Expect(executor.MockPlatformClient.NumCalledDeleteWebhook).To(Equal(0))
	})

	It("can return an error if the get call fails", func() {
		executor.MockPlatformClient.GetWebhooksError = errors.New("get webhooks error")

		_, err := executor.RunCommand(
			"webhooks",
			"delete",
		)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("get webhooks error"))
		Expect(executor.MockPlatformClient.NumCalledGetWebhooks).To(Equal(1))
		Expect(executor.MockPlatformClient.NumCalledDeleteWebhook).To(Equal(0))
	})

	It("can return an error if the create call fails", func() {
		executor.MockPlatformClient.DeleteWebhookError = errors.New("error")

		output := executor.RunCommandAndExpectSuccess(
			"webhooks",
			"delete",
			"--webhook-ids",
			"Org ID 1: https://example.com",
		)

		Expect(output).To(ContainSubstring("error"))
		Expect(output).To(ContainSubstring("Unable to delete webhook"))
		Expect(executor.MockPlatformClient.NumCalledDeleteWebhook).To(Equal(1))
	})

	It("can delete a webhook", func() {
		output, err := executor.RunCommand(
			"webhooks",
			"delete",
			"--webhook-ids",
			"Org ID 1: https://example.com",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Webhook 1 deleted"))
	})

	It("can delete several webhooks", func() {
		output, err := executor.RunCommand(
			"webhooks",
			"delete",
			"--webhook-ids",
			"Org ID 1: https://example.com,Org ID 2: https://example-2.com",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Webhook 1 deleted"))
		Expect(output).To(ContainSubstring("Webhook 2 deleted"))
		Expect(executor.MockPlatformClient.NumCalledGetWebhooks).To(Equal(1))
		Expect(executor.MockPlatformClient.NumCalledDeleteWebhook).To(Equal(2))
	})

	It("can delete several reading from user input", func() {
		input := bytes.NewBufferString("Org ID 1: https://example.com,Org ID 2: https://example-2.com\n")

		output, err := executor.RunCommandWithInput(
			input,
			"webhooks",
			"delete",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Webhook 1 deleted"))
		Expect(output).To(ContainSubstring("Webhook 2 deleted"))
		Expect(executor.MockPlatformClient.NumCalledGetWebhooks).To(Equal(1))
		Expect(executor.MockPlatformClient.NumCalledDeleteWebhook).To(Equal(2))
	})

})
