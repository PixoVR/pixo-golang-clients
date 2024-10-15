package cmd_test

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"math/rand"
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

	Context("authorization", func() {

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
		Expect(executor.MockPlatformClient.CalledCreateSessionWith).To(HaveLen(amount), "incorrect number of create session calls")
		Expect(executor.MockPlatformClient.CalledCreateEventWith).To(HaveLen(0), "incorrect number of create event calls")
		Expect(executor.MockPlatformClient.CalledUpdateSessionWith).To(HaveLen(0), "incorrect number of update session calls")
	})

	It("can load test sessions", func() {
		amount := rand.Intn(1000) + 1
		concurrent := 5
		actualConcurrent := min(amount, concurrent)
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

		Expect(output).To(ContainSubstring(fmt.Sprintf("Starting load test with %d requests and %d concurrent workers", amount, actualConcurrent)))
		Expect(output).To(ContainLineWithItems("Start Session Errors:", "0"))
		Expect(output).To(ContainLineWithItems("Create Event Errors:", "0"))
		Expect(output).To(ContainLineWithItems("Complete Session Errors:", "0"))
		Expect(output).To(ContainLineWithItems("Unsuccessful Sessions:", "0"))
		Expect(output).To(ContainLineWithItems("Sessions Started:", fmt.Sprint(amount)))
		Expect(output).To(ContainLineWithItems("Events Created:", fmt.Sprint(amount)))
		Expect(output).To(ContainLineWithItems("Sessions Completed:", fmt.Sprint(amount)))
		Expect(executor.MockPlatformClient.CalledCreateSessionWith).To(HaveLen(amount), "incorrect number of create session calls")
		Expect(executor.MockPlatformClient.CalledCreateEventWith).To(HaveLen(amount), "incorrect number of create event calls")
		Expect(executor.MockPlatformClient.CalledUpdateSessionWith).To(HaveLen(amount), "incorrect number of update session calls")
	})

	It("can load test sessions with session details", func() {
		moduleVersion := "1.00.00"
		scenario := "kitchen"
		mode := "practice"
		focus := "milkshake"
		specialization := "chocolate"

		executor.RunCommandAndExpectSuccess(
			"cannon",
			"sessions",
			"--module",
			"TST",
			"-a",
			"1",
			"--version",
			moduleVersion,
			"--mode",
			mode,
			"--scenario",
			scenario,
			"--focus",
			focus,
			"--specialization",
			specialization,
			"--score",
			"100",
			"--max-score",
			"200",
			"--passed",
		)

		Expect(executor.MockPlatformClient.CalledCreateSessionWith).To(HaveLen(1))
		Expect(executor.MockPlatformClient.CalledCreateSessionWith[0].ModuleVersion).To(Equal(moduleVersion))
		Expect(executor.MockPlatformClient.CalledCreateSessionWith[0].Mode).To(Equal(mode))
		Expect(executor.MockPlatformClient.CalledCreateSessionWith[0].Scenario).To(Equal(scenario))
		Expect(executor.MockPlatformClient.CalledCreateSessionWith[0].Focus).To(Equal(focus))
		Expect(executor.MockPlatformClient.CalledCreateSessionWith[0].Specialization).To(Equal(specialization))
		Expect(executor.MockPlatformClient.CalledCreateSessionWith[0].LessonStatus).To(Equal(""))
		Expect(executor.MockPlatformClient.CalledCreateSessionWith[0].RawScore).To(Equal(0.0))
		Expect(executor.MockPlatformClient.CalledCreateSessionWith[0].MaxScore).To(Equal(0.0))
		Expect(executor.MockPlatformClient.CalledUpdateSessionWith[0].LessonStatus).To(Equal("passed"))
		Expect(executor.MockPlatformClient.CalledUpdateSessionWith[0].RawScore).To(Equal(100.0))
		Expect(executor.MockPlatformClient.CalledUpdateSessionWith[0].MaxScore).To(Equal(200.0))
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

		Expect(executor.MockPlatformClient.CalledCreateEventWith).To(HaveLen(1))
		Expect(executor.MockPlatformClient.CalledCreateEventWith[0].Payload).To(Equal(`{"key":"value"}`))
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
		file, cleanup := NewTestFile(payload)
		defer cleanup()

		executor.RunCommandAndExpectSuccess(
			"cannon",
			"sessions",
			"--module",
			"TST",
			"--payload-file",
			file.Name(),
			"-a",
			"1",
		)

		Expect(executor.MockPlatformClient.CalledCreateEventWith).To(HaveLen(1))
		Expect(executor.MockPlatformClient.CalledCreateEventWith[0].Payload).To(Equal(`{"key":"value"}`))
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
			Expect(executor.MockHeadsetClient.CalledStartSessionWith).To(HaveLen(amount), "incorrect number of create session calls")
			Expect(executor.MockHeadsetClient.CalledSendEventWith).To(HaveLen(amount), "incorrect number of create event calls")
			Expect(executor.MockHeadsetClient.CalledEndSessionWith).To(HaveLen(amount), "incorrect number of update session calls")
		})

		It("can load test sessions with session details", func() {
			moduleVersion := "1.00.00"
			scenario := "kitchen"
			mode := "practice"
			focus := "milkshake"
			specialization := "chocolate"
			score := 100.0
			maxScore := 200.0

			executor.RunCommandAndExpectSuccess(
				"cannon",
				"sessions",
				"--legacy",
				"--module",
				"TST",
				"-a",
				"1",
				"--version",
				moduleVersion,
				"--mode",
				mode,
				"--scenario",
				scenario,
				"--focus",
				focus,
				"--specialization",
				specialization,
				"--score",
				fmt.Sprint(score),
				"--max-score",
				fmt.Sprint(maxScore),
				"--passed",
			)

			expectedStartPayload := map[string]interface{}{
				"object": map[string]interface{}{
					"id": fmt.Sprintf("https://pixovr.com/xapi/objects/%d/%s", 1, scenario),
				},
				"context": map[string]interface{}{
					"revision": moduleVersion,
					"extensions": map[string]interface{}{
						"https://pixovr.com/xapi/extension/sessionMode":           mode,
						"https://pixovr.com/xapi/extension/sessionFocus":          focus,
						"https://pixovr.com/xapi/extension/sessionSpecialization": specialization,
					},
				},
			}
			expectedEndPayload := map[string]interface{}{
				"lessonStatus": "passed",
				"result": map[string]interface{}{
					"score": map[string]interface{}{
						"raw": score,
						"max": maxScore,
					},
				},
			}
			Expect(executor.MockHeadsetClient.CalledStartSessionWith).To(HaveLen(1))
			Expect(executor.MockHeadsetClient.CalledStartSessionWith[0].Payload).To(Equal(expectedStartPayload))

			Expect(executor.MockHeadsetClient.CalledSendEventWith).To(HaveLen(1))
			Expect(executor.MockHeadsetClient.CalledSendEventWith[0].Payload).To(BeNil())

			payload := executor.MockHeadsetClient.CalledEndSessionWith[0].Payload
			Expect(payload["sessionDuration"]).NotTo(BeNil())
			delete(payload, "sessionDuration")
			Expect(executor.MockHeadsetClient.CalledEndSessionWith).To(HaveLen(1))
			Expect(executor.MockHeadsetClient.CalledEndSessionWith[0].Payload).To(Equal(expectedEndPayload))
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
			Expect(executor.MockHeadsetClient.CalledSendEventWith).To(HaveLen(1))
			Expect(executor.MockHeadsetClient.CalledSendEventWith[0].Payload).To(Equal(expectedPayload))
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

func NewTestFile(contents string) (*os.File, func()) {
	randomFilename := fmt.Sprintf("testfile-%d", rand.Intn(10000000))
	Expect(os.WriteFile(randomFilename, []byte(contents), 0644)).To(Succeed())
	cleanup := func() {
		Expect(os.Remove(randomFilename)).To(Succeed())
	}

	file, err := os.Open(randomFilename)
	Expect(err).ToNot(HaveOccurred())

	return file, cleanup
}
