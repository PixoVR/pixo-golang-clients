package cmd_test

import (
	"bytes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sessions End", func() {

	var (
		executor *TestExecutor
	)

	BeforeEach(func() {
		executor = NewTestExecutor()
	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("can return an error if the session id is missing", func() {
		input := bytes.NewReader([]byte(""))

		output, err := executor.RunCommandWithInput(
			input,
			"sessions",
			"end",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Enter SESSION ID:"))
		Expect(output).To(ContainSubstring("Session ID not provided"))
	})

	It("can return an error if the api call fails", func() {
		executor.MockPlatformClient.UpdateSessionError = true

		output, err := executor.RunCommand(
			"sessions",
			"end",
			"--session-id",
			"1",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("error"))
		Expect(executor.MockPlatformClient.CalledUpdateSession).To(BeTrue())
	})

	It("can end a session", func() {
		output, err := executor.RunCommand(
			"sessions",
			"end",
			"--session-id",
			"1",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Session completed"))
	})

	It("can end a session with scores", func() {
		output, err := executor.RunCommand(
			"sessions",
			"end",
			"--session-id",
			"1",
			"--score",
			"1",
			"--max-score",
			"3",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Session completed with score 1/3 - 33%"))
	})

})
