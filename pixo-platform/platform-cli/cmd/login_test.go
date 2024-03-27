package cmd_test

import (
	"bytes"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Login", func() {

	var (
		executor *TestExecutor
	)

	BeforeEach(func() {
		executor = NewTestExecutor()

	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("can login from user input", func() {
		input := bytes.NewReader([]byte("testuser\nfakepassword\n"))
		output, err := executor.RunCommandWithInput(
			input,
			"auth",
			"login",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Enter username:"))
		Expect(output).To(ContainSubstring("Enter password:"))
		Expect(output).To(ContainSubstring("Login successful. Here is your API token:"))

		output, err = executor.RunCommand("config")
		Expect(err).NotTo(HaveOccurred())
		Expect(output).NotTo(ContainSubstring("fakepassword"))
		Expect(output).NotTo(ContainSubstring("token"))
		Expect(output).NotTo(ContainSubstring("api-key"))
	})

})

func (t *TestExecutor) ExpectLoginToSucceed(username, password string) {

	output, err := t.RunCommand(
		"auth",
		"login",
		"--username",
		username,
		"--password",
		password,
	)

	Expect(err).NotTo(HaveOccurred())
	Expect(output).To(ContainSubstring("Login successful. Here is your API token:"))
	userID, ok := t.ConfigManager.GetConfigValue("user-id")
	Expect(ok).To(BeTrue())
	Expect(userID).To(Equal(fmt.Sprint(t.MockPlatformClient.ActiveUserID())))
}
