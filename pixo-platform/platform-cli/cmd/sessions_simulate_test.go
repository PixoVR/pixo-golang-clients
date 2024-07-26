package cmd_test

import (
	"bytes"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sessions Simulation", func() {

	BeforeEach(func() {
		executor = NewTestExecutor()
	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("can return an error if the module id is missing", func() {
		input := bytes.NewBufferString("")

		output, err := executor.RunCommandWithInput(
			input,
			"sessions",
			"simulate",
		)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("MODULE ID not provided"))
		Expect(output).To(ContainSubstring("MODULE ID"))
	})

	It("doesnt get modules if id is provided", func() {
		output, err := executor.RunCommand(
			"sessions",
			"simulate",
			"--module-id",
			"1",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Session started for module"))
		Expect(executor.MockPlatformClient.NumCalledGetModules).To(Equal(0))
	})

	It("gets modules if no id is provided", func() {
		input := bytes.NewBufferString("test\n")

		output, err := executor.RunCommandWithInput(
			input,
			"sessions",
			"simulate",
		)

		Expect(executor.MockPlatformClient.NumCalledGetModules).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("MODULE ID"))
		Expect(output).To(ContainSubstring("Session started for module test"))
	})

	It("can ask the user for the ip address if getting the ip fails", func() {
		input := bytes.NewBufferString("test\n")
		executor.MockPlatformClient.GetIPAddressError = errors.New("ip address not provided")

		output, err := executor.RunCommandWithInput(
			input,
			"sessions",
			"simulate",
		)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("ip address not provided"))
		Expect(output).To(ContainSubstring("IP ADDRESS:"))
	})

	It("can return an error if the create call fails", func() {
		input := bytes.NewBufferString("test\n")
		executor.MockPlatformClient.CreateSessionError = errors.New("create error")

		_, err := executor.RunCommandWithInput(
			input,
			"sessions",
			"simulate",
		)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("create error"))
		Expect(executor.MockPlatformClient.NumCalledCreateSession).To(Equal(1))
	})

	It("can simulate a session", func() {
		output, err := executor.RunCommand(
			"sessions",
			"simulate",
			"--module-id",
			"1",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Session started for module"))

		output, err = executor.RunCommand("config")
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Session Id: 1"))
	})

})
