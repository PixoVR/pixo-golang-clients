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

		Expect(err).To(MatchError("MODULE not provided"))
		Expect(output).To(ContainSubstring("MODULE"))
	})

	It("asks for module if none is provided", func() {
		input := bytes.NewBufferString("TST\n")

		output, err := executor.RunCommandWithInput(
			input,
			"sessions",
			"simulate",
		)

		Expect(executor.MockPlatformClient.NumCalledGetModules).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("MODULE"))
		Expect(output).To(ContainSubstring("Session started for module TST"))
	})

	It("can return an error if the create session call fails", func() {
		executor.MockPlatformClient.CreateSessionError = errors.New("create error")
		input := bytes.NewBufferString("TST\n")

		_, err := executor.RunCommandWithInput(
			input,
			"sessions",
			"simulate",
		)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("create error"))
		Expect(executor.MockPlatformClient.CalledCreateSessionWith).To(HaveLen(1))
	})

	It("can pass in mode, scenario, focus, and specialization as flags", func() {
		input := bytes.NewBufferString("n\n1\n3\n")

		output := executor.RunCommandWithInputAndExpectSuccess(
			input,
			"sessions",
			"simulate",
			"--module",
			"TST",
			"--mode",
			"Challenge",
			"--scenario",
			"test-scenario",
			"--focus",
			"test-focus",
			"--specialization",
			"test-spec",
		)

		Expect(output).To(ContainSubstring("Mode: challenge"))
		Expect(output).To(ContainSubstring("Scenario: test-scenario"))
		Expect(output).To(ContainSubstring("Focus: test-focus"))
		Expect(output).To(ContainSubstring("Specialization: test-spec"))
	})

	It("can return an error if the update session api call fails", func() {
		executor.MockPlatformClient.UpdateSessionError = errors.New("update error")
		input := bytes.NewBufferString("\n\n\n\n\nn\n100\n200\n")

		_, err := executor.RunCommandWithInput(
			input,
			"sessions",
			"simulate",
			"--module",
			"TST",
		)

		Expect(err).To(MatchError("update error"))
		Expect(executor.MockPlatformClient.CalledUpdateSessionWith).To(HaveLen(1))
	})

	It("can ask the user if the session was passed", func() {
		input := bytes.NewBufferString("\n\n\n\nn\n\n\n\n")

		output := executor.RunCommandWithInputAndExpectSuccess(
			input,
			"sessions",
			"simulate",
			"--module",
			"TST",
		)

		Expect(output).To(ContainSubstring("Session completed"))
		Expect(output).To(ContainSubstring("Passed?"))
	})

	It("can simulate a session with no events", func() {
		input := bytes.NewBufferString("\n\n\n\nn\n1\n3\n")

		output, err := executor.RunCommandWithInput(
			input,
			"sessions",
			"simulate",
			"--module",
			"TST",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(executor.MockPlatformClient.CalledCreateSessionWith).To(HaveLen(1))
		Expect(output).To(ContainSubstring("Session started for module TST"))
		Expect(output).To(ContainSubstring("Create event?"))
		Expect(executor.MockPlatformClient.CalledUpdateSessionWith).To(HaveLen(1))
		Expect(output).To(ContainSubstring("Session completed"))
		Expect(output).To(ContainSubstring("Score: 1.00/3.00"))
		Expect(output).To(ContainSubstring("Percent: 33%"))
		Expect(output).To(ContainSubstring("Duration: 1s"))

		output, err = executor.RunCommand("config")
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Session Id: 1"))
	})

	It("can return an error if theres an error creating event", func() {
		executor.MockPlatformClient.CreateEventError = errors.New("create event error")
		input := bytes.NewBufferString("\n\n\n\ny\nsome-type\n\n")

		output, err := executor.RunCommandWithInput(
			input,
			"sessions",
			"simulate",
			"--module",
			"TST",
		)

		Expect(output).To(ContainSubstring("Create event?"))
		Expect(output).To(ContainSubstring("EVENT TYPE: "))
		Expect(err).To(MatchError("create event error"))
		Expect(executor.MockPlatformClient.CalledCreateEventWith).To(HaveLen(1))
	})

	It("can simulate a session with an event and payload", func() {
		input := bytes.NewBufferString("\n\n\n\ny\nsome-event-type\n{\"some\":\"data\"}\n\n1\n3\ny\n")

		output, err := executor.RunCommandWithInput(
			input,
			"sessions",
			"simulate",
			"--module",
			"TST",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(executor.MockPlatformClient.CalledCreateSessionWith).To(HaveLen(1))
		Expect(output).To(ContainSubstring("Session started for module TST"))
		Expect(output).To(ContainSubstring("Create event?"))
		Expect(output).To(ContainSubstring("JSON PAYLOAD: "))
		Expect(executor.MockPlatformClient.CalledCreateEventWith).To(HaveLen(1))
		Expect(output).To(ContainSubstring("Event created for session"))
		Expect(executor.MockPlatformClient.CalledUpdateSessionWith).To(HaveLen(1))
		Expect(output).To(ContainSubstring("Session completed"))
		Expect(output).To(ContainSubstring("Lesson Status: passed"))
	})

	It("can simulate a session with multiple events", func() {
		input := bytes.NewBufferString("\n\n\n\ny\nsome-event-type\n\ny\nsome-event-type\n\nn\n1\n3\n")

		output, err := executor.RunCommandWithInput(
			input,
			"sessions",
			"simulate",
			"--module",
			"TST",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(executor.MockPlatformClient.CalledCreateSessionWith).To(HaveLen(1))
		Expect(output).To(ContainSubstring("Session started for module TST"))
		Expect(output).To(ContainSubstring("Create event?"))
		Expect(executor.MockPlatformClient.CalledCreateEventWith).To(HaveLen(2))
		Expect(output).To(ContainSubstring("Event created for session"))
		Expect(output).To(ContainSubstring("Session completed"))
	})

	Context("using the legacy headset api", func() {

		It("can return an error if the joined event create call fails", func() {
			executor.MockHeadsetClient.StartSessionError = errors.New("start session error")
			input := bytes.NewBufferString("")

			_, err := executor.RunCommandWithInput(
				input,
				"sessions",
				"simulate",
				"--module",
				"TST",
				"--legacy",
			)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("start session error"))
			Expect(executor.MockHeadsetClient.CalledStartSessionWith).To(HaveLen(1))
		})

		It("can return an error if the end session call fails", func() {
			executor.MockHeadsetClient.EndSessionError = errors.New("end session error")
			input := bytes.NewBufferString("\n\n\n\nn\n1\n3\n")

			_, err := executor.RunCommandWithInput(
				input,
				"sessions",
				"simulate",
				"--module",
				"TST",
				"--legacy",
			)

			Expect(err).To(MatchError("end session error"))
			Expect(executor.MockHeadsetClient.CalledEndSessionWith).To(HaveLen(1))
		})

		It("can simulate a session with no events", func() {
			input := bytes.NewBufferString("\n\n\n\nn\n1\n3\n")

			output, err := executor.RunCommandWithInput(
				input,
				"sessions",
				"simulate",
				"--module",
				"TST",
				"--legacy",
			)

			Expect(err).NotTo(HaveOccurred())
			Expect(executor.MockHeadsetClient.CalledStartSessionWith).To(HaveLen(1))
			Expect(output).To(ContainSubstring("Session started using legacy headset API"))
			Expect(output).To(ContainSubstring("Create event?"))
			Expect(executor.MockHeadsetClient.CalledEndSessionWith).To(HaveLen(1))
		})

		It("can simulate a session with an event and payload", func() {
			input := bytes.NewBufferString("\n\n\n\ny\nsome-event-type\n{\"some\":\"data\"}\n1\n3\n")

			output, err := executor.RunCommandWithInput(
				input,
				"sessions",
				"simulate",
				"--module",
				"TST",
				"--legacy",
			)

			Expect(err).NotTo(HaveOccurred())
			Expect(executor.MockHeadsetClient.CalledStartSessionWith).To(HaveLen(1))
			Expect(output).To(ContainSubstring("Session started using legacy headset API"))
			Expect(output).To(ContainSubstring("Create event?"))
			Expect(executor.MockHeadsetClient.CalledSendEventWith).To(HaveLen(1))
			Expect(output).To(ContainSubstring("Event created for session"))
			Expect(executor.MockHeadsetClient.CalledEndSessionWith).To(HaveLen(1))
		})

	})

})
