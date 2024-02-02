package cmd_test

import (
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

	It("can create a multiplayer session", func() {
		output, err := executor.RunCommand(
			"mp",
			"matchmake",
			"--module-id",
			"271",
			"--server-version",
			"2.00.01",
			"--connect",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Attempting to find a match for module 271 with server version 2.00.01"))
	})

	It("can load test matchmaking", func() {
		output, err := executor.RunCommand(
			"mp",
			"matchmake",
			"--load",
			"10",
			"--module-id",
			"271",
			"--server-version",
			"2.00.01",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Starting load test with 10 connections to"))
		//Expect(executor.MockMatchmakingClient.NumCalledFindMatch).To(Equal(10))
	})

})
