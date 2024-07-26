package cmd_test

import (
	"bytes"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sessions End", func() {

	BeforeEach(func() {
		executor = NewTestExecutor()
	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("can return an error if the session id is missing", func() {
		input := bytes.NewBufferString("")

		output, err := executor.RunCommandWithInput(
			input,
			"sessions",
			"end",
		)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("SESSION ID not provided"))
		Expect(output).To(ContainSubstring("SESSION ID:"))
	})

	It("can return an error if the update session api call fails", func() {
		executor.MockPlatformClient.UpdateSessionError = errors.New("update error")
		input := bytes.NewBufferString("100\n200\n")

		output, err := executor.RunCommandWithInput(
			input,
			"sessions",
			"end",
			"--session-id",
			"1",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("update error"))
		Expect(executor.MockPlatformClient.NumCalledUpdateSession).To(Equal(1))
	})

	It("can return an error if the create event api call fails", func() {
		executor.MockPlatformClient.PostError = errors.New("error")
		input := bytes.NewBufferString("100\n200\n")

		output := executor.RunCommandWithInputAndExpectSuccess(
			input,
			"sessions",
			"end",
			"--session-id",
			"1",
		)

		Expect(output).To(ContainSubstring("error"))
		Expect(executor.MockPlatformClient.NumCalledPost).To(Equal(1))
	})

	It("can end a session", func() {
		input := bytes.NewBufferString("100\n200\n")

		output := executor.RunCommandWithInputAndExpectSuccess(
			input,
			"sessions",
			"end",
			"--session-id",
			"1",
		)

		Expect(output).To(ContainSubstring("Session completed"))
		Expect(executor.MockPlatformClient.NumCalledUpdateSession).To(Equal(1))
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
		Expect(output).To(ContainSubstring("Session completed"))
		Expect(output).To(ContainSubstring("Score: 1/3"))
		Expect(output).To(ContainSubstring("Percent: 33%"))
		Expect(output).To(ContainSubstring("Duration: 1s"))
		Expect(executor.MockPlatformClient.NumCalledPost).To(Equal(1))
		Expect(executor.MockPlatformClient.NumCalledUpdateSession).To(Equal(1))
	})

})
