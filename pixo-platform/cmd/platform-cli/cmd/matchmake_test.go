package cmd_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
)

var _ = Describe("Matchmake", func() {

	It("can create a multiplayer session", func() {
		rootCmd, output := GetRootCmd()
		rootCmd.SetArgs([]string{
			"mp",
			"matchmake",
			"--module-id",
			"271",
			"--server-version",
			"2.00.01",
		})
		err := rootCmd.Execute()
		Expect(err).NotTo(HaveOccurred())

		out, err := io.ReadAll(output)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(out)).To(ContainSubstring("Attempting to find a match for module 271 with server version 2.00.01"))
	})

})
