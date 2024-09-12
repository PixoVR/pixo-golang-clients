package cmd_test

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"os"
	"strings"
)

var _ = Describe("Sessions Load Testing", func() {

	BeforeEach(func() {
		executor = NewTestExecutor()
		executor.MockPlatformClient.GetUserResponse = &platform.User{
			Org: platform.Org{Type: "platform"},
		}
		Expect(executor.ConfigManager.SetActiveEnv(config.Env{Lifecycle: "stage"})).To(Succeed())
	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("returns an error if unable to get the current user", func() {
		executor.MockPlatformClient.GetUserError = fmt.Errorf("get user error")
		_, err := executor.RunCommand("cannon", "sessions")
		Expect(err).To(MatchError("get user error"))
	})

	It("only allows platform users", func() {
		executor.MockPlatformClient.GetUserResponse = &platform.User{
			Org: platform.Org{Type: "trial"},
		}

		_, err := executor.RunCommand("cannon", "sessions")

		Expect(err).To(MatchError("only platform users can run load tests"))
	})

	It("isn't allowed on prod", func() {
		Expect(executor.ConfigManager.SetActiveEnv(config.Env{Lifecycle: "prod"})).To(Succeed())
		_, err := executor.RunCommand("cannon", "matchmake")
		Expect(err).To(MatchError("cannot run load tests against production"))
	})

	It("can display errors", func() {
		amount := 20
		concurrent := 5
		timeout := 1
		executor.MockPlatformClient.CreateSessionError = fmt.Errorf("create session error")

		output, err := executor.RunCommand(
			"cannon",
			"sessions",
			"--module",
			"TST",
			"--amount",
			fmt.Sprint(amount),
			"--concurrent",
			fmt.Sprint(concurrent),
			"--timeout",
			fmt.Sprint(timeout),
		)

		Expect(err).ToNot(HaveOccurred())
		Expect(output).To(ContainSubstring("create session error"))
		Expect(output).To(ContainLineWithItems("Start Session Errors:", "20"))
		Expect(output).To(ContainLineWithItems("Create Event Errors:", "0"))
		Expect(output).To(ContainLineWithItems("Complete Session Errors:", "0"))
		Expect(output).To(ContainLineWithItems("Unsuccessful Sessions:", "20"))
		Expect(output).To(ContainLineWithItems("Sessions Started:", "0"))
		Expect(output).To(ContainLineWithItems("Events Created:", "0"))
		Expect(output).To(ContainLineWithItems("Sessions Completed:", "0"))
		Expect(executor.MockPlatformClient.NumCalledCreateSession).To(Equal(amount), "incorrect number of create session calls")
		Expect(executor.MockPlatformClient.NumCalledCreateEvent).To(Equal(0), "incorrect number of create event calls")
		Expect(executor.MockPlatformClient.NumCalledUpdateSession).To(Equal(0), "incorrect number of update session calls")
	})

	It("can load test sessions", func() {
		amount := 20
		concurrent := 5
		timeout := 1

		output := executor.RunCommandAndExpectSuccess(
			"cannon",
			"sessions",
			"--module",
			"TST",
			"--amount",
			fmt.Sprint(amount),
			"--concurrent",
			fmt.Sprint(concurrent),
			"--timeout",
			fmt.Sprint(timeout),
		)

		Expect(output).To(ContainSubstring(fmt.Sprintf("Starting load test with %d requests and %d concurrent workers", amount, concurrent)))
		Expect(output).To(ContainLineWithItems("Start Session Errors:", "0"))
		Expect(output).To(ContainLineWithItems("Create Event Errors:", "0"))
		Expect(output).To(ContainLineWithItems("Complete Session Errors:", "0"))
		Expect(output).To(ContainLineWithItems("Unsuccessful Sessions:", "0"))
		Expect(output).To(ContainLineWithItems("Sessions Started:", "20"))
		Expect(output).To(ContainLineWithItems("Events Created:", "20"))
		Expect(output).To(ContainLineWithItems("Sessions Completed:", "20"))
		Expect(executor.MockPlatformClient.NumCalledCreateSession).To(Equal(amount), "incorrect number of create session calls")
		Expect(executor.MockPlatformClient.NumCalledCreateEvent).To(Equal(amount), "incorrect number of create event calls")
		Expect(executor.MockPlatformClient.NumCalledUpdateSession).To(Equal(amount), "incorrect number of update session calls")
	})

	It("can load test sessions with event payloads", func() {
		_ = executor.RunCommandAndExpectSuccess(
			"cannon",
			"sessions",
			"--module",
			"TST",
			"--payload",
			`{"key":"value"}`,
			"-a",
			"1",
		)

		expectedPayload := map[string]interface{}{"key": "value"}
		Expect(executor.MockPlatformClient.NumCalledCreateEvent).To(Equal(1))
		Expect(executor.MockPlatformClient.CalledCreateEventWith[0].Payload).To(Equal(expectedPayload))
	})

	It("returns an error if unable to read the payload file", func() {
		_, err := executor.RunCommand(
			"cannon",
			"sessions",
			"--module",
			"TST",
			"--payload-file",
			"missing.json",
		)

		Expect(err).To(MatchError("open missing.json: no such file or directory"))
	})

	It("can load test sessions with event payloads from a file", func() {
		payload := `{"key":"value"}`
		filename := "payload.json"
		Expect(os.WriteFile(filename, []byte(payload), 0644)).To(Succeed())
		defer func() {
			Expect(os.Remove(filename)).To(Succeed())
		}()

		_ = executor.RunCommandAndExpectSuccess(
			"cannon",
			"sessions",
			"--module",
			"TST",
			"--payload-file",
			filename,
			"-a",
			"1",
		)

		expectedPayload := map[string]interface{}{"key": "value"}
		Expect(executor.MockPlatformClient.NumCalledCreateEvent).To(Equal(1))
		Expect(executor.MockPlatformClient.CalledCreateEventWith[0].Payload).To(Equal(expectedPayload))
	})

	Context("legacy headset api", func() {

		It("can load test sessions", func() {
			amount := 30
			concurrent := 2
			timeout := 2

			output := executor.RunCommandAndExpectSuccess(
				"cannon",
				"sessions",
				"--legacy",
				"--module",
				"TST",
				"--amount",
				fmt.Sprint(amount),
				"--concurrent",
				fmt.Sprint(concurrent),
				"--timeout",
				fmt.Sprint(timeout),
			)

			Expect(output).To(ContainSubstring(fmt.Sprintf("Starting load test with %d requests and %d concurrent workers", amount, concurrent)))
			Expect(output).To(ContainLineWithItems("Start Session Errors:", "0"))
			Expect(output).To(ContainLineWithItems("Create Event Errors:", "0"))
			Expect(output).To(ContainLineWithItems("Complete Session Errors:", "0"))
			Expect(output).To(ContainLineWithItems("Unsuccessful Sessions:", "0"))
			Expect(output).To(ContainLineWithItems("Sessions Started:", fmt.Sprint(amount)))
			Expect(output).To(ContainLineWithItems("Events Created:", fmt.Sprint(amount)))
			Expect(output).To(ContainLineWithItems("Sessions Completed:", fmt.Sprint(amount)))
			Expect(executor.MockHeadsetClient.NumCalledStartSession).To(Equal(amount), "incorrect number of create session calls")
			Expect(executor.MockHeadsetClient.NumCalledSendEvent).To(Equal(amount), "incorrect number of create event calls")
			Expect(executor.MockHeadsetClient.NumCalledEndSession).To(Equal(amount), "incorrect number of update session calls")
		})

		It("can load test sessions with event payloads", func() {
			executor.RunCommandAndExpectSuccess(
				"cannon",
				"sessions",
				"--module",
				"TST",
				"--legacy",
				"--payload",
				`{"key":"value"}`,
				"-a",
				"1",
			)

			expectedPayload := map[string]interface{}{"key": "value"}
			Expect(executor.MockHeadsetClient.NumCalledSendEvent).To(Equal(1))
			Expect(executor.MockHeadsetClient.CalledSendEventsWith[0].Payload).To(Equal(expectedPayload))
		})

	})

})

func ContainLineWithItems(items ...string) types.GomegaMatcher {
	return &containsLineWithItemsMatcher{
		expected: items,
	}
}

type containsLineWithItemsMatcher struct {
	expected []string
}

// checks for expected.Seconds() +- 1 second
func (m *containsLineWithItemsMatcher) Match(actual interface{}) (success bool, err error) {
	actualStr := actual.(string)
	actualStrLines := strings.Split(actualStr, "\n")
	for _, line := range actualStrLines {
		var numFound int
		for _, item := range m.expected {
			if strings.Contains(line, item) {
				numFound++
			}
		}
		if numFound == len(m.expected) {
			return true, nil
		}
	}

	return false, nil
}

func (m *containsLineWithItemsMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected %s to contain %s", actual, m.expected)
}

func (m *containsLineWithItemsMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected %s not to contain %s", actual, m.expected)
}
