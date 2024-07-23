package cmd_test

import (
	"bytes"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Webhooks Delete", func() {

	var (
		executor *TestExecutor
	)

	BeforeEach(func() {
		executor = NewTestExecutor()
	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("can return an error if the id is not provided", func() {
		input := bytes.NewReader([]byte(""))

		output := executor.RunCommandWithInputAndExpectSuccess(
			input,
			"webhooks",
			"delete",
		)

		Expect(output).To(ContainSubstring("ID is required"))
		Expect(executor.MockPlatformClient.NumCalledDeleteWebhook).To(Equal(0))
	})

	It("can return an error if the api call fails", func() {
		executor.MockPlatformClient.DeleteWebhookError = errors.New("error")

		output := executor.RunCommandAndExpectSuccess(
			"webhooks",
			"delete",
			"--webhook-id",
			"1",
		)

		Expect(output).To(ContainSubstring("error"))
		Expect(output).To(ContainSubstring("Unable to delete webhook"))
		Expect(executor.MockPlatformClient.NumCalledDeleteWebhook).To(Equal(1))
	})

	It("can delete a webhook", func() {
		output, err := executor.RunCommand(
			"webhooks",
			"delete",
			"--webhook-id",
			"1",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Webhook deleted"))
	})

})
