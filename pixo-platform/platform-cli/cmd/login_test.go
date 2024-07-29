package cmd_test

import (
	"bytes"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Login", func() {

	BeforeEach(func() {
		executor = NewTestExecutor()
	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("can login with flags", func() {
		output, err := executor.RunCommand(
			"auth",
			"login",
			"--username",
			"testuser1",
			"--password",
			"fakepassword",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Login successful. Here is your API token:"))
		userID, ok := executor.ConfigManager.GetConfigValue("auth-user-id")
		Expect(ok).To(BeTrue())
		Expect(userID).To(Equal(fmt.Sprint(executor.MockPlatformClient.ActiveUserID())))
	})

	It("can login from user input", func() {
		input := bytes.NewBufferString("testuser2\nfakepassword\n")

		output, err := executor.RunCommandWithInput(
			input,
			"auth",
			"login",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Enter USERNAME:"))
		Expect(output).To(ContainSubstring("Enter PASSWORD:"))
		Expect(output).To(ContainSubstring("Login successful. Here is your API token:"))

		output, err = executor.RunCommand("config")
		Expect(err).NotTo(HaveOccurred())
		Expect(output).NotTo(ContainSubstring("fakepassword"))
		Expect(output).NotTo(ContainSubstring("token"))
		Expect(output).NotTo(ContainSubstring("api-key"))
	})

	It("can return an error if unable to use the api key", func() {
		executor.MockPlatformClient.GetControlTypesError = fmt.Errorf("get roles error")

		_, err := executor.RunCommand(
			"auth",
			"login",
			"--key",
			"fake-key",
		)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("invalid API key"))
	})

	It("can login with an api key", func() {
		output, err := executor.RunCommand(
			"auth",
			"login",
			"--key",
			"fake-key",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Login with API key successful."))

		output, err = executor.RunCommand("config")
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("API Key:"))
		Expect(output).NotTo(ContainSubstring("fake-key"))
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
	userID, ok := t.ConfigManager.GetConfigValue("auth-user-id")
	Expect(ok).To(BeTrue())
	Expect(userID).To(Equal(fmt.Sprint(t.MockPlatformClient.ActiveUserID())))
}
