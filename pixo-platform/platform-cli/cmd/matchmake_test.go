package cmd_test

import (
	"bytes"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Matchmake", func() {

	BeforeEach(func() {
		executor = NewTestExecutor()
	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("should ask for the module id and server version if it is not provided", func() {
		input := bytes.NewBufferString("TST\n")
		output, err := executor.RunCommandWithInput(
			input,
			"mp",
			"matchmake",
		)
		Expect(err).To(MatchError("SERVER VERSION not provided"))
		Expect(output).To(ContainSubstring("MODULE"))
		Expect(output).To(ContainSubstring("SERVER VERSION"))
	})

	It("can return an error if module id is missing", func() {
		reader := bytes.NewBufferString("\n")

		output, err := executor.RunCommandWithInput(
			reader,
			"mp",
			"matchmake",
		)

		Expect(err).To(MatchError("MODULE not provided"))
		Expect(output).To(ContainSubstring("MODULE"))
		Expect(executor.MockMatchmakingClient.NumCalledDialWebsocket).To(Equal(0))
		Expect(executor.MockMatchmakingClient.NumCalledWriteToWebsocketError).To(Equal(0))
		Expect(executor.MockMatchmakingClient.NumCalledReadFromWebsocket).To(Equal(0))
		Expect(executor.MockMatchmakingClient.NumCalledCloseWebsocket).To(Equal(0))
	})

	It("can return an error if server version is missing", func() {
		reader := bytes.NewBufferString("\n")

		_, err := executor.RunCommandWithInput(
			reader,
			"mp",
			"matchmake",
			"--module",
			"TST",
		)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("SERVER VERSION not provided"))
		Expect(executor.MockMatchmakingClient.NumCalledDialWebsocket).To(Equal(0))
		Expect(executor.MockMatchmakingClient.NumCalledWriteToWebsocketError).To(Equal(0))
		Expect(executor.MockMatchmakingClient.NumCalledReadFromWebsocket).To(Equal(0))
		Expect(executor.MockMatchmakingClient.NumCalledCloseWebsocket).To(Equal(0))
	})

	It("can return an error if the matchmaking request fails", func() {
		executor.MockMatchmakingClient.FindMatchError = errors.New("matchmaking error")

		output, err := executor.RunCommand(
			"mp",
			"matchmake",
			"--module",
			"TST",
			"--server-version",
			"1.00.00",
		)

		Expect(output).To(ContainSubstring("Attempting to find a match"))
		Expect(err).To(MatchError("matchmaking error"))
		Expect(executor.MockMatchmakingClient.NumCalledFindMatch).To(Equal(1))
		Expect(executor.MockMatchmakingClient.NumCalledDialGameserver).To(Equal(0))
		Expect(executor.MockMatchmakingClient.NumCalledCloseGameserver).To(Equal(0))
	})

	It("can connect to the gameserver connection and exit", func() {
		reader := bytes.NewBufferString("ACK\nexit\n")

		output := executor.RunCommandWithInputAndExpectSuccess(
			reader,
			"mp",
			"matchmake",
			"--module",
			"TST",
			"--server-version",
			"1.00.00",
			"--connect",
		)

		Expect(output).To(ContainSubstring("Attempting to find a match"))
		Expect(executor.MockMatchmakingClient.NumCalledFindMatch).To(Equal(1))
		Expect(executor.MockMatchmakingClient.NumCalledDialGameserver).To(Equal(1))
		Expect(executor.MockMatchmakingClient.NumCalledSendToGameserver).To(Equal(1))
		Expect(executor.MockMatchmakingClient.NumCalledReadFromGameserver).To(Equal(1))
		Expect(executor.MockMatchmakingClient.NumCalledCloseGameserver).To(Equal(1))
	})

	It("can perform a single matchmaking request and connect in a subsequent command", func() {
		output := executor.RunCommandAndExpectSuccess(
			"mp",
			"matchmake",
			"--module",
			"TST",
			"--server-version",
			"1.00.00",
		)
		Expect(output).To(ContainSubstring("Attempting to find a match"))
		Expect(executor.MockMatchmakingClient.NumCalledFindMatch).To(Equal(1))

		input := bytes.NewBufferString("exit\n")
		output = executor.RunCommandWithInputAndExpectSuccess(
			input,
			"mp",
			"--connect",
		)
		Expect(output).To(ContainSubstring("Connecting to gameserver at"))
		Expect(executor.MockMatchmakingClient.NumCalledDialGameserver).To(Equal(1))
		Expect(output).To(ContainSubstring("Enter message to gameserver:"))
		Expect(output).To(ContainSubstring("Closing connection to gameserver at"))
		Expect(executor.MockMatchmakingClient.NumCalledCloseGameserver).To(Equal(1))
	})

})
