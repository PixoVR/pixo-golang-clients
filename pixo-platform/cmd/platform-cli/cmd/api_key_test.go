package cmd_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
)

var _ = Describe("API Keys", func() {

	It("can create an api key", func() {
		rootCmd, output := GetRootCmd()
		rootCmd.SetArgs([]string{
			"apiKeys",
			"create",
		})

		err := rootCmd.Execute()

		Expect(err).NotTo(HaveOccurred())
		out, err := io.ReadAll(output)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(out)).To(ContainSubstring("API key created"))

		rootCmd.SetArgs([]string{
			"config",
			"list",
		})

		err = rootCmd.Execute()

		Expect(err).NotTo(HaveOccurred())
		out, err = io.ReadAll(output)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(out)).To(ContainSubstring("api-key : "))
	})

	It("can list api keys", func() {
		rootCmd, output := GetRootCmd()
		rootCmd.SetArgs([]string{
			"apiKeys",
			"list",
		})

		err := rootCmd.Execute()

		Expect(err).NotTo(HaveOccurred())
		out, err := io.ReadAll(output)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(out)).To(ContainSubstring("API Keys:"))
	})

})
