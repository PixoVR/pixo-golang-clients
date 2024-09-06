package cmd_test

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
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
