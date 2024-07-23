package cmd_test

import (
	"bytes"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Webhooks Create", func() {

	var (
		executor *TestExecutor
	)

	BeforeEach(func() {
		executor = NewTestExecutor()
	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("can return an error if the url is missing", func() {
		input := bytes.NewReader([]byte(""))

		output, err := executor.RunCommandWithInput(
			input,
			"webhooks",
			"create",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Enter URL:"))
		Expect(output).To(ContainSubstring("URL not provided"))
		Expect(executor.MockPlatformClient.NumCalledCreateWebhook).To(Equal(0))
	})

	It("can return an error if the description is missing", func() {
		input := bytes.NewReader([]byte("https://example.com\n"))

		output, err := executor.RunCommandWithInput(
			input,
			"webhooks",
			"create",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Enter description:"))
		Expect(output).To(ContainSubstring("DESCRIPTION not provided"))
		Expect(executor.MockPlatformClient.NumCalledCreateWebhook).To(Equal(0))
	})

	It("can return an error if the api call fails", func() {
		executor.MockPlatformClient.CreateWebhookError = errors.New("error")

		output := executor.RunCommandAndExpectSuccess(
			"webhooks",
			"create",
			"--url",
			"https://example.com",
			"--description",
			"test",
		)

		Expect(output).To(ContainSubstring("error"))
		Expect(executor.MockPlatformClient.NumCalledCreateWebhook).To(Equal(1))
	})

	It("can create a webhook", func() {
		output, err := executor.RunCommand(
			"webhooks",
			"create",
			"--url",
			"https://example.com",
			"--description",
			"test",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Webhook created"))
		Expect(executor.MockPlatformClient.NumCalledCreateWebhook).To(Equal(1))
	})

})
