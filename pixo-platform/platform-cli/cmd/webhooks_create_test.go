package cmd_test

import (
	"bytes"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Webhooks Create", func() {

	BeforeEach(func() {
		executor = NewTestExecutor()
	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("can return an error if the url is missing", func() {
		input := bytes.NewBufferString("")

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
		input := bytes.NewBufferString("https://example.com\n")

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

	It("asks the user if they want the token generated", func() {
		input := bytes.NewBufferString("yes\n")

		output, err := executor.RunCommandWithInput(
			input,
			"webhooks",
			"create",
			"--url",
			"https://example.com",
			"--description",
			"test",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Generate token?"))
		Expect(output).NotTo(ContainSubstring("Enter TOKEN"))
		Expect(executor.MockPlatformClient.NumCalledCreateWebhook).To(Equal(1))
	})

	It("asks the user for the token if they don't want it generated", func() {
		input := bytes.NewBufferString("no\n\n")

		output, err := executor.RunCommandWithInput(
			input,
			"webhooks",
			"create",
			"--url",
			"https://example.com",
			"--description",
			"test",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Generate token?"))
		Expect(output).To(ContainSubstring("Enter WEBHOOK TOKEN"))
		Expect(output).To(ContainSubstring("No token provided. Webhook will be insecure"))
		Expect(executor.MockPlatformClient.NumCalledCreateWebhook).To(Equal(1))
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
