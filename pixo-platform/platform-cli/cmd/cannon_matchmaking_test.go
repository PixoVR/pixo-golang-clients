package cmd_test

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Matchmaking Load Testing", func() {

	BeforeEach(func() {
		executor = NewTestExecutor()
		executor.MockPlatformClient.GetUserResponse = &platform.User{
			Org: platform.Org{Type: "platform"},
		}
		Expect(executor.ConfigManager.SetActiveEnv(config.Env{Lifecycle: "dev"})).To(Succeed())

	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("returns an error if unable to get the current user", func() {
		executor.MockPlatformClient.GetUserError = fmt.Errorf("get user error")
		_, err := executor.RunCommand("cannon", "matchmake")
		Expect(err).To(MatchError("get user error"))
	})

	It("only allows platform users", func() {
		executor.MockPlatformClient.GetUserResponse = &platform.User{
			Org: platform.Org{Type: "trial"},
		}

		_, err := executor.RunCommand("cannon", "matchmake")

		Expect(err).To(MatchError("only platform users can run load tests"))
	})

	It("isn't allowed on prod", func() {
		Expect(executor.ConfigManager.SetActiveEnv(config.Env{Lifecycle: "prod"})).To(Succeed())
		_, err := executor.RunCommand("cannon", "matchmake")
		Expect(err).To(MatchError("cannot run load tests against production"))
	})

	It("can load test matchmaking", func() {
		expectedNumRequests := 50
		output := executor.RunCommandAndExpectSuccess(
			"cannon",
			"matchmake",
			"--module",
			"TST",
			"--server-version",
			"1.00.00",
		)

		Expect(output).To(ContainSubstring(fmt.Sprintf("Starting load test with %d requests and 5 concurrent workers", expectedNumRequests)))
		Expect(output).To(ContainLineWithItems("Connection Errors:", "0"))
		Expect(output).To(ContainLineWithItems("Matching Errors:", "0"))
		Expect(output).To(ContainLineWithItems("Matches Received:", fmt.Sprint(expectedNumRequests)))
		Expect(output).To(ContainLineWithItems("Gameservers Received:", "1"))
		Expect(executor.MockMatchmakingClient.NumCalledDialWebsocket).To(Equal(expectedNumRequests), "incorrect number of dial matchmaker calls")
		Expect(executor.MockMatchmakingClient.NumCalledWriteToWebsocketError).To(Equal(expectedNumRequests), "incorrect number of write to websocket calls")
		Expect(executor.MockMatchmakingClient.NumCalledReadFromWebsocket).To(Equal(expectedNumRequests), "incorrect number of read from websocket calls")
		Expect(executor.MockMatchmakingClient.NumCalledCloseWebsocket).To(Equal(expectedNumRequests), "incorrect number of close websocket calls")
	})

})
