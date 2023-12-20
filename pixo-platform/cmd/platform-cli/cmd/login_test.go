package cmd_test

import (
	"bytes"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/cmd"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
	"os"
)

var _ = Describe("Login", func() {

	It("can authenticate with the platform", func() {
		command := cmd.RootCmd()
		b := bytes.NewBufferString("")
		command.SetOut(b)
		command.SetArgs([]string{
			"auth",
			"login",
			"--username",
			os.Getenv("PIXO_USERNAME"),
			"--password",
			os.Getenv("PIXO_PASSWORD"),
		})

		err := command.Execute()

		Expect(err).NotTo(HaveOccurred())

		out, err := io.ReadAll(b)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(out)).To(ContainSubstring("Login successful. Here is your API token:"))
	})

})
