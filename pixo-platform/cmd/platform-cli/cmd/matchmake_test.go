package cmd_test

import (
	"bytes"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Matchmake", func() {

	var (
		executor *TestExecutor
	)

	BeforeEach(func() {
		executor = NewTestExecutor()
	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("should ask for the module id and server version if it is not provided", func() {
		input := bytes.NewReader([]byte("1\n"))
		output, err := executor.RunCommandWithInput(
			input,
			"mp",
			"matchmake",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Enter MODULE ID:"))
		Expect(output).To(ContainSubstring("Enter SERVER VERSION:"))
		Expect(output).To(ContainSubstring("Server version not provided"))
	})

	It("can return an error if module id is missing", func() {
		reader := bytes.NewReader([]byte("0\n"))
		output, err := executor.RunCommandWithInput(
			reader,
			"mp",
			"matchmake",
			"--module-id",
			"0",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Enter MODULE ID:"))
		Expect(output).To(ContainSubstring("Module ID not provided"))
		Expect(executor.MockMatchmakingClient.NumCalledDialWebsocket).To(Equal(0))
		Expect(executor.MockMatchmakingClient.NumCalledWriteToWebsocket).To(Equal(0))
		Expect(executor.MockMatchmakingClient.NumCalledReadFromWebsocket).To(Equal(0))
	})

	It("can return an error if server version is missing", func() {
		reader := bytes.NewReader([]byte("\n"))
		output, err := executor.RunCommandWithInput(
			reader,
			"mp",
			"matchmake",
			"--module-id",
			"1",
			"--server-version",
			"",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Server version not provided"))
		Expect(executor.MockMatchmakingClient.NumCalledDialWebsocket).To(Equal(0))
		Expect(executor.MockMatchmakingClient.NumCalledWriteToWebsocket).To(Equal(0))
		Expect(executor.MockMatchmakingClient.NumCalledReadFromWebsocket).To(Equal(0))
	})

	It("can perform a single matchmaking request", func() {
		output, err := executor.RunCommand(
			"mp",
			"matchmake",
			"--module-id",
			"1",
			"--server-version",
			"1.00.00",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Attempting to find a match"))
		Expect(executor.MockMatchmakingClient.NumCalledFindMatch).To(Equal(1))
	})

	It("can load test matchmaking", func() {
		num := 100
		output, err := executor.RunCommand(
			"mp",
			"matchmake",
			"--load",
			fmt.Sprint(num),
			"--module-id",
			"1",
			"--server-version",
			"1.00.00",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring(fmt.Sprintf("Starting load test with %d connections to", num)))
		Expect(executor.MockMatchmakingClient.NumCalledDialWebsocket).To(BeNumerically(">", 0))
		Expect(executor.MockMatchmakingClient.NumCalledWriteToWebsocket).To(BeNumerically(">", 0))
		Expect(executor.MockMatchmakingClient.NumCalledReadFromWebsocket).To(BeNumerically(">", 0))
	})

})
