package cmd_test

import (
	"bytes"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sessions Start", func() {

	var (
		executor *TestExecutor
	)

	BeforeEach(func() {
		executor = NewTestExecutor()
	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("can return an error if the module id is missing", func() {
		input := bytes.NewReader([]byte(""))

		output, err := executor.RunCommandWithInput(
			input,
			"sessions",
			"start",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Enter MODULE ID:"))
		Expect(output).To(ContainSubstring("Module ID not provided"))
	})

	It("can ask the user for the ip address if getting the ip fails", func() {
		input := bytes.NewReader([]byte("1\n"))
		executor.MockPlatformClient.GetIPAddressError = errors.New("ip address not provided")

		output, err := executor.RunCommandWithInput(
			input,
			"sessions",
			"start",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Enter IP ADDRESS:"))
		Expect(output).To(ContainSubstring("ip address not provided"))
	})

	It("can return an error if the api call fails", func() {
		input := bytes.NewReader([]byte("1"))
		executor.MockPlatformClient.CreateSessionError = true

		output, err := executor.RunCommandWithInput(
			input,
			"sessions",
			"start",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("error"))
		Expect(executor.MockPlatformClient.CalledCreateSession).To(BeTrue())
	})

	It("can start a session", func() {
		output, err := executor.RunCommand(
			"sessions",
			"start",
			"--module-id",
			"1",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Session started for module 1 with ID 1"))
	})

})
